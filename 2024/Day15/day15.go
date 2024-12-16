package main

import (
	util "2024/until"
	"fmt"
	"slices"
	"strings"
)

const (
	up       = "^"
	down     = "v"
	left     = "<"
	right    = ">"
	robot    = "@"
	edge     = "#"
	box      = "O"
	empty    = "."
	leftBox  = "["
	rightBox = "]"
)

var directions = map[string][2]int{
	up:    {-1, 0},
	down:  {1, 0},
	left:  {0, -1},
	right: {0, 1},
}

type boxx struct {
	leftSide  [2]int
	rightSide [2]int
}

/*
	Advent of Code Day 15:
		Part 1: Straight forward grid traversal with some recursive logic to move stuff around
		Part 2: I haven't successfully gotten this to work, still need to do some debugging. I got my answer but I will need to fix it more
				I collected the effected boxes with a BFS, but I this I need to fix my logic around when I can move them
*/

func main() {
	input := util.ReadInput(util.Parameter())
	splitIndex := slices.Index(input, "")
	grid, robotDirections := util.TransformStringSliceInto2DMatrix(input[:splitIndex]), strings.Split(strings.Join(input[splitIndex:], ""), "")
	doubledGrid := doubleGrid(grid)
	//part1(grid, robotDirections)
	printGrid(doubledGrid, "Initial State:")
	part2(doubledGrid, robotDirections)
}

func part1(grid [][]string, robotDirections []string) {
	robotPosition := findRobot(grid)

	for _, dir := range robotDirections {
		move := directions[dir]
		item := moveItem(grid, [2]int{robotPosition[0] + move[0], robotPosition[1] + move[1]}, move)
		if item == empty {
			grid[robotPosition[0]][robotPosition[1]] = empty
			robotPosition = [2]int{robotPosition[0] + move[0], robotPosition[1] + move[1]}
			grid[robotPosition[0]][robotPosition[1]] = robot
		}
		//printGrid(grid, dir)
	}
	gps := calculateGPS(grid)
	fmt.Println("(part 1) Ans: ", gps)
}

func part2(grid [][]string, robotDirections []string) {
	robotPosition := findRobot(grid)

	for _, dir := range robotDirections {
		move := directions[dir]
		targetPos := [2]int{robotPosition[0] + move[0], robotPosition[1] + move[1]}

		// Handle vertical movement into wide boxes
		if dir == up || dir == down {
			if grid[targetPos[0]][targetPos[1]] == leftBox || grid[targetPos[0]][targetPos[1]] == rightBox {
				if moveVerticalBoxes(grid, targetPos, move) {
					grid[robotPosition[0]][robotPosition[1]] = empty
					robotPosition = targetPos
					grid[robotPosition[0]][robotPosition[1]] = robot
				}
			} else if grid[targetPos[0]][targetPos[1]] == empty {
				grid[robotPosition[0]][robotPosition[1]] = empty
				robotPosition = targetPos
				grid[robotPosition[0]][robotPosition[1]] = robot
			}
		} else { // Handle horizontal movement
			item := moveItem(grid, [2]int{robotPosition[0] + move[0], robotPosition[1] + move[1]}, move)
			if item == empty {
				grid[robotPosition[0]][robotPosition[1]] = empty
				robotPosition = [2]int{robotPosition[0] + move[0], robotPosition[1] + move[1]}
				grid[robotPosition[0]][robotPosition[1]] = robot
			}
		}
	}

	gps := calculateGPS(grid)
	fmt.Println("(part 2) Ans: ", gps)
}

func moveVerticalBoxes(grid [][]string, robotPos [2]int, move [2]int) bool {
	boxes := findConnectedBoxes(grid, robotPos, move)
	// Validate if all boxes can move

	rows := make(map[int][]boxx)
	counter := make(map[int]int)
	for _, box := range boxes {
		currRow := box.left[0]
		rows[currRow] = append(rows[currRow], boxx{box.left, box.right})
		counter[currRow]++
	}

	for i := len(boxes) - 1; i >= 0; i-- {
		newLeft := [2]int{boxes[i].left[0] + move[0], boxes[i].left[1] + move[1]}
		newRight := [2]int{boxes[i].right[0] + move[0], boxes[i].right[1] + move[1]}
		if grid[newLeft[0]][newLeft[1]] == empty && grid[newRight[0]][newRight[1]] == empty {
			counter[boxes[i].left[0]]--
			if counter[boxes[i].left[0]] == 0 {
				// can move all boxes in the row
				for _, box := range rows[boxes[i].left[0]] {
					newLeft := [2]int{box.leftSide[0] + move[0], box.leftSide[1] + move[1]}
					newRight := [2]int{box.rightSide[0] + move[0], box.rightSide[1] + move[1]}
					grid[box.leftSide[0]][box.leftSide[1]] = empty
					grid[box.rightSide[0]][box.rightSide[1]] = empty
					grid[newLeft[0]][newLeft[1]] = leftBox
					grid[newRight[0]][newRight[1]] = rightBox
				}
			}
		} else {
			return false
		}
	}

	return true
}

func findConnectedBoxes(grid [][]string, start [2]int, move [2]int) []struct {
	left  [2]int // Position of `[`
	right [2]int // Position of `]`
} {
	boxes := []struct {
		left  [2]int
		right [2]int
	}{}

	queue := [][2]int{start} // Start from the robot's initial position
	visited := make(map[[2]int]bool)

	// BFS to find all connected boxes in the path
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// Skip if already visited
		if visited[curr] {
			continue
		}
		visited[curr] = true

		// Check if we are at the left side of a box
		if grid[curr[0]][curr[1]] == leftBox {
			boxes = append(boxes, struct {
				left  [2]int
				right [2]int
			}{
				left:  curr,
				right: [2]int{curr[0], curr[1] + 1},
			})
			// Add the next position in the movement direction to the queue
			next := [2]int{curr[0] + move[0], curr[1] + move[1]}
			queue = append(queue, next)
		}

		// Check if we are at the right side of a box
		if grid[curr[0]][curr[1]] == rightBox {
			// Redirect to the corresponding left part of the box
			left := [2]int{curr[0], curr[1] - 1}
			next := [2]int{curr[0] + move[0], curr[1] + move[1]}
			if !visited[left] {
				queue = append(queue, left)
			}
			if !visited[next] {
				queue = append(queue, next)
			}
		}
	}

	return boxes
}

func calculateGPS(grid [][]string) int {
	toReturn := 0
	for r := 1; r < len(grid)-1; r++ {
		for c := 1; c < len(grid[0])-1; c++ {
			if grid[r][c] == box || grid[r][c] == leftBox {
				toReturn += (100 * r) + c
			}
		}
	}
	return toReturn
}

func doubleGrid(grid [][]string) [][]string {
	toReturn := make([][]string, len(grid))
	for r := 0; r < len(grid); r++ {
		row := make([]string, 0)
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] == box {
				row = append(row, "[")
				row = append(row, "]")
			} else if grid[r][c] == robot {
				row = append(row, robot)
				row = append(row, empty)
			} else {
				row = append(row, strings.Split(strings.Repeat(grid[r][c], 2), "")...)
			}
		}
		toReturn[r] = row
	}
	return toReturn
}

func printGrid(grid [][]string, move string) {
	fmt.Println(move)
	for _, row := range grid {
		fmt.Println(strings.Join(row, ""))
	}
	fmt.Println()
}

func moveItem(grid [][]string, itemPos [2]int, move [2]int) string {
	// base case is hit wall exit
	if grid[itemPos[0]][itemPos[1]] == empty || grid[itemPos[0]][itemPos[1]] == edge {
		return grid[itemPos[0]][itemPos[1]]
	}

	dx, dy := itemPos[0]+move[0], itemPos[1]+move[1]

	item := moveItem(grid, [2]int{dx, dy}, move)
	if item == empty {
		prev := grid[itemPos[0]][itemPos[1]]
		grid[itemPos[0]][itemPos[1]] = empty
		grid[dx][dy] = prev
	}

	return item
}

func findRobot(grid [][]string) [2]int {
	for r := 1; r < len(grid)-1; r++ {
		for c := 1; c < len(grid[r])-1; c++ {
			if grid[r][c] == robot {
				return [2]int{r, c}
			}
		}
	}
	return [2]int{-1, -1}
}

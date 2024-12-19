package main

import (
	util "2024/until"
	"fmt"
	"slices"
	"sort"
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
	part1(grid, robotDirections)
	//printGrid(doubledGrid, "Initial State:")
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
				if moveVerticalBoxes(grid, targetPos, move, dir) {
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
		//	printGrid(grid, dir)
	}

	gps := calculateGPS(grid)
	fmt.Println("(part 2) Ans: ", gps)
}

func moveVerticalBoxes(grid [][]string, robotPos [2]int, move [2]int, dir string) bool {
	boxes := findConnectedBoxes(grid, robotPos, move)
	// Validate if all boxes can move

	rows := make(map[int][]boxx)
	keys := make([]int, 0)
	// sort boxes by row, ascending if moving up, descending if moving down
	for _, box := range boxes {
		row := box.leftSide[0]
		if !slices.Contains(keys, row) {
			keys = append(keys, row)
		}
		rows[row] = append(rows[row], box)
	}

	sort.Slice(keys, func(i, j int) bool {
		if dir == up {
			return keys[i] < keys[j]
		}
		return keys[i] > keys[j]
	})

	simulatedGrid := make([][]string, len(grid))
	for r := 0; r < len(grid); r++ {
		simulatedGrid[r] = make([]string, len(grid[0]))
		copy(simulatedGrid[r], grid[r])
	}

	for _, row := range keys {
		boxes := rows[row]
		// Check if the boxes can move
		for _, box := range boxes {
			left, right := box.leftSide, box.rightSide
			if simulatedGrid[left[0]+move[0]][left[1]+move[1]] == empty && simulatedGrid[right[0]+move[0]][right[1]+move[1]] == empty {
				moveItem(simulatedGrid, left, move)
				moveItem(simulatedGrid, right, move)
			} else {
				return false
			}

		}
	}

	// If all boxes can move, update the grid
	for r := 0; r < len(grid); r++ {
		copy(grid[r], simulatedGrid[r])
	}

	return true
}

func findConnectedBoxes(grid [][]string, start [2]int, move [2]int) []boxx {

	boxes := make([]boxx, 0)
	visited := make(map[boxx]bool)
	var initialBox boxx
	if grid[start[0]][start[1]] == leftBox {
		initialBox = boxx{start, [2]int{start[0], start[1] + 1}}
	} else {
		initialBox = boxx{[2]int{start[0], start[1] - 1}, start}
	}

	queue := []boxx{initialBox}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if _, ok := visited[current]; ok {
			continue
		}
		visited[current] = true
		boxes = append(boxes, current)

		// Check if there is a box to the left
		if grid[current.leftSide[0]+move[0]][current.leftSide[1]+move[1]] == leftBox {
			leftBox := boxx{[2]int{current.leftSide[0] + move[0], current.leftSide[1] + move[1]}, [2]int{current.leftSide[0] + move[0], current.leftSide[1] + move[1] + 1}}
			queue = append(queue, leftBox)
		}

		if grid[current.leftSide[0]+move[0]][current.leftSide[1]+move[1]] == rightBox {
			leftBox := boxx{[2]int{current.leftSide[0] + move[0], current.leftSide[1] + move[1] - 1}, [2]int{current.leftSide[0] + move[0], current.leftSide[1] + move[1]}}
			queue = append(queue, leftBox)
		}
		// Check if there is a box to the right
		if grid[current.rightSide[0]+move[0]][current.rightSide[1]+move[1]] == rightBox {
			rightBox := boxx{[2]int{current.rightSide[0] + move[0], current.rightSide[1] + move[1] - 1}, [2]int{current.rightSide[0] + move[0], current.rightSide[1] + move[1]}}
			queue = append(queue, rightBox)
		}

		if grid[current.rightSide[0]+move[0]][current.rightSide[1]+move[1]] == leftBox {
			rightBox := boxx{[2]int{current.rightSide[0] + move[0], current.rightSide[1] + move[1]}, [2]int{current.rightSide[0] + move[0], current.rightSide[1] + move[1] + 1}}
			queue = append(queue, rightBox)
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

package main

import (
	"2024/util"
	"fmt"
	"os"
)

const (
	up       = "^"
	down     = "V"
	left     = "<"
	right    = ">"
	obstacle = "#"
	free     = "."
)

type state struct {
	X, Y             int
	currentDirection string
}

/*
	Day 6 Advent of Code
	Part 1: Just needed to follow the path of the guard and capture all his unique positions
	Part 2: I got stuck on this for a long time, I was trying to find some novel solution where on the guards path I could calculate the best spot to put an obstacle then detect the loop.
			I gave up and took the brute force path and simulated the guards path for every possible obstacle placement, and this ended up working
*/

func main() {
	data := util.TransformStringSliceInto2DMatrix(util.ReadInput(util.Parameter()))

	changeDirections := map[string]string{
		up:    right,
		right: down,
		down:  left,
		left:  up,
	}

	moveDirections := map[string][2]int{
		up:    {-1, 0},
		down:  {1, 0},
		left:  {0, -1},
		right: {0, 1},
	}

	part1(data, changeDirections, moveDirections)
	part2(data, changeDirections, moveDirections)
}

// part1 finds all the unique positions of the guard's path
func part1(grid [][]string, changeDirections map[string]string, moveDirections map[string][2]int) {

	seenSpaces := util.NewHashSet()

	currentX, currentY, currentDirection := findInitialPosition(grid)
	if currentX == -1 {
		fmt.Println("Something is wrong with the grid, no guard was found")
		os.Exit(1)
	}

	seenSpaces.Add([2]int{currentX, currentY})
	stillInGrid := true
	for stillInGrid {

		// check space in front, if there is a obstacle then turn
		if obstacleInFront(currentX, currentY, grid, currentDirection, moveDirections) {
			currentDirection = changeDirections[currentDirection]
		} else {
			// move once space
			dx, dy := moveDirections[currentDirection][0], moveDirections[currentDirection][1]
			currentX, currentY = currentX+dx, currentY+dy

			stillInGrid = inBounds(currentX, currentY, grid)

			if !seenSpaces.Contains([2]int{currentX, currentY}) && stillInGrid {
				seenSpaces.Add([2]int{currentX, currentY})
			}
		}

	}

	fmt.Println("(part 1) Ans: ", seenSpaces.Size())

}

// part2 calculates every possible obstacle position to force loops in the guard's path
func part2(grid [][]string, changeDirections map[string]string, moveDirections map[string][2]int) {
	validObstacles := 0

	R, C := len(grid), len(grid[0]) // Grid dimensions

	// Iterate over all possible positions in the grid
	for obsR := 0; obsR < R; obsR++ {
		for obsC := 0; obsC < C; obsC++ {
			// Skip if the cell already has an obstacle
			if grid[obsR][obsC] == "#" {
				continue
			}

			// Simulate guard's movement with the obstacle at (obsR, obsC)
			if causedLoop(grid, obsR, obsC, moveDirections, changeDirections) {
				validObstacles++
			}
		}
	}

	fmt.Printf("(part 2) Ans: %d\n", validObstacles)
}

// causedLoop is the helper function for part2, where it simulates the guards path until we possibly encounter obsR & obsC
// if we do then we force the guard to turn and simulate hitting a obstacle
func causedLoop(grid [][]string, obsR, obsC int, moveDirections map[string][2]int, changeDirections map[string]string) bool {
	R, C := len(grid), len(grid[0])
	startingRow, startingCol, currentDirection := findInitialPosition(grid) // Guard's starting position
	SEEN := make(map[state]bool)                                            // Track visited (r, c, direction)

	r, c := startingRow, startingCol
	for {
		currentState := state{r, c, currentDirection}
		if SEEN[currentState] {
			return true // Loop detected
		}
		SEEN[currentState] = true

		// Calculate next position
		dr, dc := moveDirections[currentDirection][0], moveDirections[currentDirection][1]
		nr, nc := r+dr, c+dc

		// Check for grid exit
		if nr < 0 || nr >= R || nc < 0 || nc >= C {
			return false // Guard exited the grid, no loop
		}

		// Simulate obstacle at (obsR, obsC)
		if grid[nr][nc] == "#" || (nr == obsR && nc == obsC) {
			// Turn right
			currentDirection = changeDirections[currentDirection]
		} else {
			// Move forward
			r, c = nr, nc
		}
	}
}

// obstacleInFront helper function for part 1, checks if the square in front of guard has be obstacle
func obstacleInFront(x int, y int, grid [][]string, direction string, moveDirections map[string][2]int) bool {
	dx, dy := moveDirections[direction][0]+x, moveDirections[direction][1]+y
	if inBounds(dx, dy, grid) {
		if grid[dx][dy] == obstacle {
			return true
		}
	}
	return false
}

// findInitialPosition locates the guard's initial position
func findInitialPosition(grid [][]string) (int, int, string) {
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] == up || grid[r][c] == down || grid[r][c] == left || grid[r][c] == right {
				return r, c, grid[r][c]
			}
		}
	}
	return -1, -1, ""
}

// inBounds is a helper bounds checking function
func inBounds(x, y int, grid [][]string) bool {
	if x < 0 || x >= len(grid) || y < 0 || y >= len(grid[0]) {
		return false
	}
	return true
}

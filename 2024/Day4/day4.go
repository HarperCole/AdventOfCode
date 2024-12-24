package main

import (
	"2024/util"
	"fmt"
	"strings"
)

/*
	Advent of Code Day 4:
	Part 1: Search the grid of words for target XMAS. It can be in any direction
	Part 2: Search the grid of workds for target:
	M S
	 A
	M S
	It can be in any configuration where A is the center and spells MAS along the axis

	Both of these parts were fairly straight forward. Part 1 was simple 2D matrix search looking for the target. Boundary checking is very important here
	Part 2 added the interesting wrinkle, however it was a simple matter of searching for instances of "A", then checking the diagonals
*/

const targetPart1 = "XMAS"

func main() {

	input := util.TransformStringSliceInto2DMatrix(util.ReadInput(util.Parameter()))

	part1(input)
	part2(input)

}

// part1 search for "XMAS" inside of grid. Need to check every direction that's possible
func part1(grid [][]string) {
	rows := len(grid)
	cols := len(grid[0])
	target := strings.Split(targetPart1, "")
	directions := [][2]int{
		{0, -1},  // Left
		{0, 1},   // Right
		{-1, 0},  // Up
		{1, 0},   // Down
		{-1, -1}, // Up-Left
		{-1, 1},  // Up-Right
		{1, -1},  // Down-Left
		{1, 1},   // Down-Right
	}
	matches := 0

	// checkDirections is a helper function that once we have found an "X" we proceed to check every direction for the rest of the word,
	// We conduct a DFS for the rest of the target word by using k which will scale the dx, dy to the next correct grid position and check that it matches with the correct position in the target word
	// returns false if incorrect next character or becomes out of bounds
	// returns true if match is found
	checkDirections := func(x, y, dx, dy int) bool {
		for k := 0; k < len(targetPart1); k++ {
			nx, ny := x+k*dx, y+k*dy
			if nx < 0 || nx >= rows || ny < 0 || ny >= cols || grid[nx][ny] != target[k] {
				return false
			}
		}
		return true
	}

	// basic grid traversal coupled with O(len(directions)) looping through every direction from the grid position
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == target[0] {
				for _, dir := range directions {
					dx, dy := dir[0], dir[1]
					if checkDirections(r, c, dx, dy) {
						matches++
					}
				}
			}
		}
	}

	fmt.Println("(part 1) Ans: ", matches)

}

// part2 searches grid for x MAS pattern, searches for "A" in grid, then checks the diagonals for "M","S"
func part2(grid [][]string) {
	rows := len(grid)
	cols := len(grid[0])
	target := "A"
	directions := [][2]int{
		{-1, -1}, // Up-Left
		{-1, 1},  // Up-Right
		{1, -1},  // Down-Left
		{1, 1},   // Down-Right
	}
	matches := 0
	combinationMap := map[string]string{
		"M": "S",
		"S": "M",
	}

	// checkCorner ensures that the position is within the grid, and makes sure that letter matches the correctLetter
	checkCorner := func(x, y int, correctLetter string) bool {
		if x < 0 || x >= rows || y < 0 || y >= cols || grid[x][y] != correctLetter {
			return false
		}
		return true
	}

	inBounds := func(x, y int) bool {
		return x >= 0 && x < rows && y >= 0 && y < cols
	}

	for r := 1; r < rows-1; r++ {
		for c := 1; c < cols-1; c++ {
			if grid[r][c] == target {
				upLeftX, upLeftY := directions[0][0]+r, directions[0][1]+c
				upRightX, upRightY := directions[1][0]+r, directions[1][1]+c
				if inBounds(upLeftX, upLeftY) && inBounds(upRightX, upRightY) {
					upLeftLetter := grid[upLeftX][upLeftY]
					upRightLetter := grid[upRightX][upRightY]
					bottomLEftX, bottomLeftY := directions[2][0]+r, directions[2][1]+c
					bottomRightX, bottomRightY := directions[3][0]+r, directions[3][1]+c
					targetBottomRightLetter := combinationMap[upLeftLetter]
					targetBottomLeftLetter := combinationMap[upRightLetter]
					if checkCorner(bottomRightX, bottomRightY, targetBottomRightLetter) && checkCorner(bottomLEftX, bottomLeftY, targetBottomLeftLetter) {
						matches++
					}
				}
			}
		}
	}

	fmt.Println("(part 2) Ans: ", matches)

}

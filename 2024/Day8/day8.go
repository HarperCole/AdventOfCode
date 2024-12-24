package main

import (
	util /until"
	"fmt"
)

const empty = "."

// position holds the x and y coordinates of an antinode
type position struct {
	x, y int
}

/*
	Advent of Code Day 8
	part 1: Need to calculate the antidotes between two pairs of antennas, need to calculate the directional vector dx = x2 - x1, dy = y2 - y1, calculate the two antinotes (above & below)
			above = (x2 + dx, y2 + dy) below = (x1 - dx, y1 - dy) then check if this is within hte bounds of the grid. Add to set to track unique positions
	part 2: So need to calculate all the valid antinodes along the line between the two antenna, this includes the antenna positions
*/

func main() {
	grid := util.TransformStringSliceInto2DMatrix(util.ReadInput(util.Parameter()))
	antennaPositions := findAntennaPositions(grid)
	part1(grid, antennaPositions)
	part2(grid, antennaPositions)
}

// findAntennaPositions finds all antenna positions (non empty) spots in the grid
func findAntennaPositions(grid [][]string) map[string][]position {
	positions := make(map[string][]position)

	for y, row := range grid {
		for x, cell := range row {
			if cell != empty {
				positions[cell] = append(positions[cell], position{x, y})
			}
		}
	}
	return positions
}

// part1 solves part1 as described above
func part1(grid [][]string, locations map[string][]position) {
	antinodeLocations := util.NewHashSet()
	for _, positions := range locations {
		for i := 0; i < len(positions); i++ {
			for j := i + 1; j < len(positions); j++ {
				x1, y1 := positions[i].x, positions[i].y
				x2, y2 := positions[j].x, positions[j].y

				dx, dy := x2-x1, y2-y1
				antinode1 := position{x1 - dx, y1 - dy}
				antinode2 := position{x2 + dx, y2 + dy}
				possiblePositions := []position{antinode1, antinode2}
				for _, pos := range possiblePositions {
					if withInBounds(pos.x, pos.y, len(grid), len(grid[0])) {
						antinodeLocations.Add(pos)
					}
				}
			}
		}
	}

	fmt.Println("(part 1) Ans: ", antinodeLocations.Size())
}

// part2 solves part2 as described above
func part2(grid [][]string, locations map[string][]position) {
	antinodeLocations := util.NewHashSet()

	// Add all antennas as valid antinodes
	for _, positions := range locations {
		for _, antenna := range positions {
			antinodeLocations.Add(antenna)
		}
	}

	// Process pairs of antennas to compute antinodes
	for _, positions := range locations {
		for i := 0; i < len(positions); i++ {
			for j := i + 1; j < len(positions); j++ {
				x1, y1 := positions[i].x, positions[i].y
				x2, y2 := positions[j].x, positions[j].y

				dx, dy := x2-x1, y2-y1
				expandAntidote(antinodeLocations, x1, y1, -dx, -dy, grid)
				expandAntidote(antinodeLocations, x2, y2, dx, dy, grid)
			}
		}
	}

	fmt.Println("(part 2) Ans: ", antinodeLocations.Size())
}

// expandAntidote is a helper function for calculating the harmonic resonance of the antenna pairs
func expandAntidote(locations *util.HashSet, startX int, startY int, dx int, dy int, grid [][]string) {
	x, y := startX, startY

	for {
		x, y = x+dx, y+dy
		if !withInBounds(x, y, len(grid), len(grid[0])) {
			break
		}

		if locations.Contains(position{x, y}) {
		} else {
			locations.Add(position{x, y})
		}
	}
}

// withInBounds is a helper function for calculating the bounds of the grid
func withInBounds(x int, y int, rowLen int, colLen int) bool {
	return x >= 0 && x < rowLen && y >= 0 && y < colLen
}

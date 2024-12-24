package main

import (
	"2024/util"
	"container/list"
	"fmt"
)

const (
	up    = "up"
	down  = "down"
	left  = "left"
	right = "right"
)

// region represents a distinct area in the grid with cells of the same crop type.
//
// Fields:
//   - area: A slice of [2]int representing the row and column indices of all cells in the region.
//   - crop: The crop type represented as a string.
//   - perimeter: The calculated perimeter of the region, based on its boundary cells.
//   - cost: The fencing cost for the region, typically calculated as perimeter * area size.
//   - discount: The calculated discount for the region, typically based on its corners and area size.
type region struct {
	area      [][2]int // List of coordinates of all cells in the region.
	crop      string   // The crop type of the region.
	perimeter int      // The total perimeter of the region.
	cost      int      // The calculated fencing cost for the region.
	discount  int      // The calculated discount for the region.
}

/*
	Advent of Code Day 12:
		Part 1: Very simple flood fill concept, decided to create the region struct to better track the data
		Part 2: This part really stumped me. At first, I thought this was some interval tracking to where you needed to track the direction and continuous interval,
				then I saw a hint on reddit where someone said that you the number of sides is equal to the number of corners of the region. This makes sense because in the context of this problem
				all shapes that can be created in this grid will be simple polygons
*/

func main() {
	input := util.TransformStringSliceInto2DMatrix(util.ReadInput(util.Parameter()))
	regions := computeGrid(input)
	part1(regions)
	part2(regions)
}

// part1 computes and prints the total cost of fencing all regions.
// The cost is calculated as the product of the region's perimeter and its area size.
func part1(regions []region) {

	costs := 0
	for _, reg := range regions {
		costs += reg.cost
	}

	fmt.Println("(part 1) Ans: ", costs)
}

// part2 computes and prints the total discount for fencing all regions.
// The discount is calculated as the product of the region's area and the number of corners in the region.
func part2(regions []region) {

	discount := 0
	for _, reg := range regions {
		discount += reg.discount
	}
	fmt.Println("(part 2) Ans: ", discount)
}

// computeGrid identifies distinct regions in the grid and computes their properties.
//
// It uses a flood-fill algorithm to find connected regions, calculates their perimeter,
// and computes the cost and discount for each region based on its area and corners.
//
// Returns a slice of `region` structs, where each struct represents a region's properties.
func computeGrid(input [][]string) []region {
	seen := make(map[[2]int]bool)
	rows := len(input)
	cols := len(input[0])

	regions := make([]region, 0)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if !seen[[2]int{r, c}] {
				regions = append(regions, captureRegion(input, r, c, seen, input[r][c]))
			}
		}

	}

	for index, _ := range regions {
		regions[index].perimeter = calculatePerimeter(input, regions[index])
		regions[index].cost = regions[index].perimeter * len(regions[index].area)
		regions[index].discount = len(regions[index].area) * calculateCorners(input, regions[index])
	}

	return regions
}

// calculateCorners calculates the total number of corners (both exterior and interior)
// for a given region in the grid.
//
// A corner is defined as:
// - An exterior corner if neither orthogonal neighbor matches the region's crop.
// - An interior corner if both orthogonal neighbors match the region's crop, but the diagonal does not.
//
// Returns the total count of corners for the given region.
func calculateCorners(input [][]string, r region) int {
	// Directions for orthogonal neighbors
	directions := map[string][2]int{
		up:    {-1, 0},
		down:  {1, 0},
		left:  {0, -1},
		right: {0, 1},
	}

	corners := 0

	// Check all diagonal corners
	for _, reg := range r.area {
		// Define diagonal and corresponding orthogonal neighbors
		cornerChecks := []struct {
			diag   [2]int // Diagonal position
			ortho1 [2]int // First orthogonal position
			ortho2 [2]int // Second orthogonal position
		}{
			{[2]int{reg[0] + directions[up][0], reg[1] + directions[left][1]}, // NW
				[2]int{reg[0] + directions[up][0], reg[1]},
				[2]int{reg[0], reg[1] + directions[left][1]}},
			{[2]int{reg[0] + directions[up][0], reg[1] + directions[right][1]}, // NE
				[2]int{reg[0] + directions[up][0], reg[1]},
				[2]int{reg[0], reg[1] + directions[right][1]}},
			{[2]int{reg[0] + directions[down][0], reg[1] + directions[right][1]}, // SE
				[2]int{reg[0] + directions[down][0], reg[1]},
				[2]int{reg[0], reg[1] + directions[right][1]}},
			{[2]int{reg[0] + directions[down][0], reg[1] + directions[left][1]}, // SW
				[2]int{reg[0] + directions[down][0], reg[1]},
				[2]int{reg[0], reg[1] + directions[left][1]}},
		}

		// Check each corner
		for _, check := range cornerChecks {
			if exteriorCorner(check.ortho1, check.ortho2, input, r.crop) || interiorCorner(check.ortho1, check.ortho2, check.diag, input, r.crop) {
				corners++
			}
		}
	}

	return corners
}

// interiorCorner checks if a given corner is an interior corner.
//
// A corner is considered interior if:
// - Both orthogonal neighbors match the region's crop.
// - The diagonal position does not match the region's crop.
// - All positions are within the bounds of the grid.
//
// Returns true if the corner is an interior corner, false otherwise.
func interiorCorner(ortho1, ortho2, diag [2]int, input [][]string, crop string) bool {
	// Ensure all positions are within bounds
	if !isWithinBounds(ortho1[0], ortho1[1], len(input), len(input[0])) ||
		!isWithinBounds(ortho2[0], ortho2[1], len(input), len(input[0])) ||
		!isWithinBounds(diag[0], diag[1], len(input), len(input[0])) {
		return false
	}

	// Both orthogonal neighbors must match, diagonal must not
	return input[ortho1[0]][ortho1[1]] == crop &&
		input[ortho2[0]][ortho2[1]] == crop &&
		input[diag[0]][diag[1]] != crop
}

// exteriorCorner checks if a given corner is an exterior corner.
//
// A corner is considered exterior if:
// - Neither orthogonal neighbor matches the region's crop.
// - Either one or both orthogonal neighbors are out of bounds.
//
// Returns true if the corner is an exterior corner, false otherwise.
func exteriorCorner(ortho1, ortho2 [2]int, input [][]string, crop string) bool {
	// Check if orthogonal neighbors are within bounds and do not match the crop
	out1 := !isWithinBounds(ortho1[0], ortho1[1], len(input), len(input[0])) || input[ortho1[0]][ortho1[1]] != crop
	out2 := !isWithinBounds(ortho2[0], ortho2[1], len(input), len(input[0])) || input[ortho2[0]][ortho2[1]] != crop

	return out1 && out2
}

// isWithinBounds checks if a given cell position is within the bounds of the grid.
//
// Returns true if the position is within bounds, false otherwise.
func isWithinBounds(r, c, rows, cols int) bool {
	return r >= 0 && c >= 0 && r < rows && c < cols
}

// captureRegion performs a flood-fill to identify all cells in a connected region.
//
// Starting from a specific cell, it traverses all adjacent cells of the same crop type,
// marking them as visited and adding them to the region.
//
// Returns a `region` struct representing the connected region's properties.
func captureRegion(grid [][]string, row int, col int, seen map[[2]int]bool, target string) region {
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	reg := region{}
	reg.area = make([][2]int, 0)
	reg.crop = target
	queue := list.New()
	queue.PushBack([2]int{row, col})

	for queue.Len() > 0 {
		element := queue.Remove(queue.Front()).([2]int)

		if seen[element] {
			continue
		}
		seen[element] = true
		reg.area = append(reg.area, element)

		for _, direction := range directions {
			dr, dc := direction[0]+element[0], direction[1]+element[1]
			if dr >= 0 && dr < len(grid) && dc >= 0 && dc < len(grid[0]) && grid[dr][dc] == target && !seen[[2]int{dr, dc}] {
				queue.PushBack([2]int{dr, dc})
			}
		}
	}

	return reg
}

// calculatePerimeter calculates the total perimeter of a given region in the grid.
//
// The perimeter is calculated by checking each cell in the region and subtracting
// the number of adjacent cells of the same crop from the total potential perimeter (4 per cell).
//
// Returns the total perimeter for the region.
func calculatePerimeter(grid [][]string, selectedRegion region) int {
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	totalPerimeter := 0
	for _, area := range selectedRegion.area {
		areaPerimeter := 4
		for _, direction := range directions {
			dr, dc := area[0]+direction[0], area[1]+direction[1]
			if dr >= 0 && dr < len(grid) && dc >= 0 && dc < len(grid[0]) && grid[dr][dc] == selectedRegion.crop {
				areaPerimeter -= 1
			}
		}
		totalPerimeter += areaPerimeter
	}
	return totalPerimeter
}

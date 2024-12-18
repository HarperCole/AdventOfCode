package main

import (
	util "2024/until"
	"container/list"
	"fmt"
	"regexp"
)

const (
	pattern    = `-?\d+`
	part1Bytes = 0x400
)

var direction = [][2]int{
	{0, 1},
	{0, -1},
	{1, 0},
	{-1, 0},
}

type grid struct {
	pos   [2]int
	steps int
}

func main() {
	input := util.ReadInput(util.Parameter())
	corruptedCoordinates := parseCoordinates(input)
	part1(corruptedCoordinates, 71, 71, part1Bytes)
	part2(corruptedCoordinates, 71, 71)
}

// part1 calculates the minimum number of steps to reach the end position in a grid,
// avoiding corrupted coordinates, and prints the result.
//
// Parameters:
// - coordinates: a slice of 2-element integer arrays representing corrupted coordinates.
// - col: the number of columns in the grid.
// - row: the number of rows in the grid.
// - bytes: the number of corrupted coordinates to consider.
func part1(coordinates [][2]int, col, row, bytes int) {
	start := [2]int{0, 0}
	end := [2]int{row - 1, col - 1}

	corruptedCoordinates := util.NewHashSet()
	for i := 0; i < bytes; i++ {
		corruptedCoordinates.Add(coordinates[i])
	}

	steps := bfs(*corruptedCoordinates, start, end, row, col)
	fmt.Println("Part 1: Ans", steps)
}

// part2 finds the first corrupted coordinate that makes the end position unreachable
// in a grid and prints the result.
//
// Parameters:
// - coordinates: a slice of 2-element integer arrays representing corrupted coordinates.
// - col: the number of columns in the grid.
// - row: the number of rows in the grid.
func part2(coordinates [][2]int, col, row int) {
	start := [2]int{0, 0}
	end := [2]int{row - 1, col - 1}

	corruptedCoordinates := util.NewHashSet()
	for i := 0; i < len(coordinates); i++ {
		corruptedCoordinates.Add(coordinates[i])
		steps := bfs(*corruptedCoordinates, start, end, row, col)
		if steps == -1 {
			fmt.Println("Part 2: Ans", coordinates[i])
			break
		}
	}
}

// bfs performs a breadth-first search to find the minimum number of steps
// from the start position to the end position in a grid, avoiding corrupted coordinates.
//
// Parameters:
// - coordinates: a HashSet containing the corrupted coordinates to avoid.
// - current: the starting position as a 2-element integer array.
// - end: the target position as a 2-element integer array.
// - row: the number of rows in the grid.
// - col: the number of columns in the grid.
//
// Returns:
//   - The minimum number of steps to reach the end position from the start position.
//     Returns -1 if the end position is not reachable.
func bfs(coordinates util.HashSet, current [2]int, end [2]int, row, col int) int {

	// Initialize a queue for BFS and add the starting position with 0 steps.
	queue := list.New()
	queue.PushBack(grid{current, 0})

	// Initialize a set to keep track of visited positions and add the starting position.
	visited := util.NewHashSet()
	visited.Add(current)

	// Perform BFS until the queue is empty.
	for queue.Len() > 0 {
		// Dequeue the front element.
		position := queue.Remove(queue.Front()).(grid)

		// Check if the current position is the end position.
		if position.pos == end {
			return position.steps
		}

		// Explore all possible directions from the current position.
		for _, dir := range direction {
			newPos := [2]int{position.pos[0] + dir[0], position.pos[1] + dir[1]}

			// Skip positions that are out of bounds or corrupted.
			if newPos[0] < 0 || newPos[0] >= row || newPos[1] < 0 || newPos[1] >= col {
				continue
			}
			if coordinates.Contains(newPos) {
				continue
			}

			// Skip positions that have already been visited.
			if visited.Contains(newPos) {
				continue
			}

			// Mark the new position as visited and add it to the queue with incremented steps.
			visited.Add(newPos)
			queue.PushBack(grid{newPos, position.steps + 1})
		}
	}

	// Return -1 if the end position is not reachable.
	return -1
}

// parseCoordinates takes a slice of strings as input and returns a slice of 2-element integer arrays.
// Each string in the input is expected to contain two integer coordinates.
// The function uses a regular expression to extract the coordinates from each string and converts them to integers.
func parseCoordinates(input []string) [][2]int {
	// Initialize an empty slice to store the parsed coordinates.
	toReturn := make([][2]int, 0)

	// Compile the regular expression pattern to match integers.
	reg := regexp.MustCompile(pattern)

	// Iterate over each line in the input slice.
	for _, line := range input {
		// Find all substrings in the line that match the regular expression.
		coordinates := reg.FindAllString(line, -1)
		// Convert the matched substrings to integers and append them as a 2-element array to the result slice.
		toReturn = append(toReturn, [2]int{util.ParseInt(coordinates[0]), util.ParseInt(coordinates[1])})
	}

	// Return the slice of parsed coordinates.
	return toReturn
}

package main

import (
	util "2024/until"
	"container/list"
	"fmt"
	"strings"
	"sync"
)

// position in the grid
//   - row denotes row index
//   - col denotes col index
//   - height denotes the height at grid[row][col]
type position struct {
	row, col, height int
}

/*
	Advent of Code Day 10:
		Part 1: Given the grid that represents the topography of an area, find the number of unique peaks (9) that can be found from a trailhead (0). I did this by passing through the grid
				locating a trailhead, then performing BFS to find all the reachable peaks
		Part 2: Find all the unique paths from a trailhead to each reachable peak. Funny enough, with how I solved part 1, the BFS I implemented was already finding all the unique paths.
*/

func main() {
	input := convertData(util.ReadInput(util.Parameter()))
	part1(input)
	part2(input)
}

// searchForTrailHead searches for a trailhead, and then begins the BFS to find reachable peaks
//   - uniquePaths indicates whether to find all uniquePaths or not
func searchForTrailHead(input [][]int, uniquePaths bool) int {
	rows := len(input)
	cols := len(input[0])

	results := make(chan int)
	var wg sync.WaitGroup

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if input[row][col] == 0 {
				wg.Add(1)
				go bfs(input, position{row, col, 0}, results, &wg, uniquePaths)
			}
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for result := range results {
		sum += result
	}

	return sum
}

func part1(input [][]int) {
	scoredTrailheads := searchForTrailHead(input, false)
	fmt.Println("(part 1) Ans: ", scoredTrailheads)
}

func part2(input [][]int) {
	uniquePaths := searchForTrailHead(input, true)
	fmt.Println("(part 2) Ans: ", uniquePaths)
}

// bfs conducts a breath first search from the given starting position
func bfs(grid [][]int, start position, results chan int, wg *sync.WaitGroup, uniquePaths bool) {
	defer wg.Done()

	// set up
	rows := len(grid)
	cols := len(grid[0])

	queue := list.New()
	seenPeaks := util.NewHashSet()
	queue.PushBack([]position{start}) // queue stores path found
	score := 0
	directions := [][2]int{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0},
	}

	// loop over queue while is still populated
	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)

		currentPath := element.Value.([]position)
		current := currentPath[len(currentPath)-1]

		// if a valid path has been found from a trailhead to a peak need to score
		// if we just want to find the number of reachable peaks from a trailhead,
		// then add peak to set, otherwise we're looking for unique paths
		if current.height == 9 && (!seenPeaks.Contains(current) || uniquePaths) {
			if !uniquePaths {
				seenPeaks.Add(current)
			}
			score++

			continue
		}

		// looping over directions map and add valid possible next positions to queue
		// to be valid it must be in bounds and strictly current height + 1
		for _, direction := range directions {
			nr, nc := direction[0]+current.row, direction[1]+current.col
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
				if grid[nr][nc] == current.height+1 {
					//
					newPath := append([]position{}, currentPath...)
					newPath = append(newPath, position{nr, nc, grid[nr][nc]})
					queue.PushBack(newPath)
				}

			}
		}

	}

	results <- score
}

func convertData(input []string) [][]int {
	toReturn := make([][]int, 0)

	for _, line := range input {
		split := strings.Split(line, "")
		converted := make([]int, 0)
		for _, num := range split {
			converted = append(converted, util.ParseInt(num))
		}
		toReturn = append(toReturn, converted)
	}
	return toReturn
}

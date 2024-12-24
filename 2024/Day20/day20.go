package main

import (
	"2024/util"
	"container/list"
	"fmt"
	"math"
)

type point struct {
	x int
	y int
}

type cheat struct {
	original point
	end      point
}

const (
	startChar = "S"
	endChar   = "E"
)

func main() {
	grid := util.TransformStringSliceInto2DMatrix(util.ReadInput(util.Parameter()))
	start, end := find(startChar, grid), find(endChar, grid)
	part1(grid, start, end, 2)
	part2(grid, start, end, 20)
}

func part1(grid [][]string, start point, end point, cheatDistance int) {
	// need to find initial shortest path with BFS
	// then need to calculate shortcuts

	originalRoute := bfs(grid, start, end)
	cheatRoute := findCheatPaths(originalRoute, cheatDistance)

	totalSaved := 0
	for distance, count := range cheatRoute {
		//fmt.Println("Distance: ", distance, "Count: ", count)
		if distance >= 100 {
			totalSaved += count
		}
	}

	fmt.Println("(part 1) Total saved distance above 100: ", totalSaved)
}

func part2(grid [][]string, start point, end point, cheatDistance int) {
	// need to find initial shortest path with BFS
	// then need to calculate shortcuts

	originalRoute := bfs(grid, start, end)
	cheatRoute := findCheatPaths(originalRoute, cheatDistance)

	totalSaved := 0
	for distance, count := range cheatRoute {
		//fmt.Println("Distance: ", distance, "Count: ", count)
		if distance >= 100 {
			totalSaved += count
		}
	}

	fmt.Println("(part 2) Total saved distance above 100: ", totalSaved)
}

// findCheatPaths identifies potential shortcuts in a given route and calculates the distance saved by each shortcut.
//
// Parameters:
// - route: A map where the keys are points representing the route and the values are the distances from the start point.
// - cheatDistance: The maximum distance to consider for potential shortcuts.
//
// Returns:
// - A map where the keys are the distances saved by the shortcuts and the values are the counts of how many times each distance is saved.
func findCheatPaths(route map[point]int, cheatDistance int) map[int]int {
	cheats := make(map[cheat]int)

	for currentStep, routeDistance := range route {
		potentialMoves := getMoves(currentStep, cheatDistance)
		for _, possibleMove := range potentialMoves {
			currentRouteDistance, stillInRoute := route[possibleMove]
			if stillInRoute {
				distanceSaved := currentRouteDistance - routeDistance - calculateDistance(currentStep, possibleMove)
				if distanceSaved > 0 {
					cheats[cheat{currentStep, possibleMove}] = distanceSaved
				}
			}
		}
	}

	cheatRoute := make(map[int]int)
	for _, distance := range cheats {
		cheatRoute[distance]++
	}
	return cheatRoute
}

// bfs performs a breadth-first search (BFS) on a 2D grid to find the shortest path from the start point to the end point.
// It returns a map where the keys are points visited during the search and the values are the order in which they were visited.
//
// Parameters:
// - grid: A 2D slice of strings representing the grid.
// - start: The starting point as a point struct.
// - end: The ending point as a point struct.
//
// Returns:
// - A map where the keys are points visited during the search and the values are the order in which they were visited.
func bfs(grid [][]string, start point, end point) map[point]int {
	visited := make(map[point]int)
	queue := list.New()
	queue.PushBack(start)

	for queue.Len() > 0 {
		current := queue.Remove(queue.Front()).(point)

		visited[current] = len(visited)

		if current == end {
			return visited
		}

		for _, possibleMove := range getMoves(current, 1) {
			if _, ok := visited[possibleMove]; ok {
				continue
			}
			if grid[possibleMove.x][possibleMove.y] == "#" {
				continue
			}
			queue.PushBack(possibleMove)
		}
	}
	return nil
}

// getMoves generates a list of valid neighboring points within a specified cheat distance.
//
// Parameters:
// - current: The current point as a point struct.
// - cheatDistance: The maximum distance to consider for neighboring points.
//
// Returns:
// - A slice of point structs representing the valid neighboring points.
func getMoves(current point, cheatDistance int) []point {
	validNeighbors := make([]point, 0)
	for y := cheatDistance * -1; y <= cheatDistance; y++ {
		for x := cheatDistance * -1; x <= cheatDistance; x++ {
			neighbor := point{current.x + x, current.y + y}
			distance := calculateDistance(current, neighbor)
			if distance > 0 && distance <= cheatDistance {
				validNeighbors = append(validNeighbors, neighbor)
			}
		}
	}
	return validNeighbors
}

// calculateDistance computes the Manhattan distance between two points.
//
// Parameters:
// - current: The current point as a point struct.
// - neighbor: The neighboring point as a point struct.
//
// Returns:
// - An integer representing the Manhattan distance between the current point and the neighboring point.
func calculateDistance(current point, neighbor point) int {
	x := math.Abs(float64(current.x - neighbor.x))
	y := math.Abs(float64(current.y - neighbor.y))
	return int(x + y)
}

// find locates the first occurrence of a specified character in a 2D grid.
// It returns the coordinates of the character as a point struct.
// If the character is not found, it returns a point with coordinates (-1, -1).
//
// Parameters:
// - char: The character to search for in the grid.
// - grid: A 2D slice of strings representing the grid.
//
// Returns:
// - A point struct containing the coordinates of the character if found, or (-1, -1) if not found.
func find(char string, grid [][]string) point {
	rows, cols := len(grid), len(grid[0])
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == char {
				return point{r, c}
			}
		}
	}
	return point{-1, -1}
}

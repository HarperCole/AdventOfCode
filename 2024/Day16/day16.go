package main

import (
	"2024/Day16/datastructure"
	"2024/util"
	"container/heap"
	"fmt"
	"math"
	"slices"
	"strings"
)

const (
	wall     = "#"
	start    = "S"
	end      = "E"
	turnCost = 1000
)

var directions = map[int][2]int{
	0: {0, 1},  // East
	1: {1, 0},  // South
	2: {0, -1}, // West
	3: {-1, 0}, // North
}

var rotations = [2]int{1, -1}

type setItem struct {
	pos [2]int
	dir int
}

type costInfo struct {
	cost  int
	paths int
}

func main() {
	grid := util.TransformStringSliceInto2DMatrix(util.ReadInput(util.Parameter()))
	startPOS, endPOS := findPOS(grid, start), findPOS(grid, end)
	part1(grid, startPOS, endPOS)
	part2(grid, startPOS, endPOS)
}

func part1(grid [][]string, start [2]int, end [2]int) {
	cost := dijkstra(grid, start, end)
	fmt.Println("(part 1) Ans:", cost)
}

func part2(grid [][]string, start [2]int, end [2]int) {
	spaces := findAllMinPathsAndSpaces(grid, start, end)
	fmt.Println("(part 2) Ans:", spaces.Size())
	addVisitedSpots(grid, spaces)
	printGrid(grid)
}

func findAllMinPathsAndSpaces(grid [][]string, start [2]int, end [2]int) util.HashSet {
	rows, cols := len(grid), len(grid[0])
	queue := datastructure.NewPriorityQueue()
	heap.Push(queue, &datastructure.PriorityQueueItem{Position: start, Direction: 0, Cost: 0})

	cost := make(map[setItem]costInfo)
	pred := make(map[setItem][]setItem)
	for queue.Len() > 0 {
		item := heap.Pop(queue).(*datastructure.PriorityQueueItem)

		if cosItem, ok := cost[setItem{item.Position, item.Direction}]; ok && item.Cost > cosItem.cost {
			continue
		}

		newPos := [2]int{item.Position[0] + directions[item.Direction][0], item.Position[1] + directions[item.Direction][1]}
		if 0 <= newPos[0] && newPos[0] < rows && 0 <= newPos[1] && newPos[1] < cols && grid[newPos[0]][newPos[1]] != wall {
			newCost := item.Cost + 1

			if cosInfo, ok := cost[setItem{newPos, item.Direction}]; !ok || newCost < cosInfo.cost {
				cost[setItem{newPos, item.Direction}] = costInfo{newCost, 1}
				pred[setItem{newPos, item.Direction}] = []setItem{{item.Position, item.Direction}}
				heap.Push(queue, &datastructure.PriorityQueueItem{Position: newPos, Direction: item.Direction, Cost: newCost})
			} else if newCost == cosInfo.cost && !slices.Contains(pred[setItem{newPos, item.Direction}], setItem{item.Position, item.Direction}) {
				cost[setItem{newPos, item.Direction}] = costInfo{newCost, cosInfo.paths + 1}
				pred[setItem{newPos, item.Direction}] = append(pred[setItem{newPos, item.Direction}], setItem{item.Position, item.Direction})
			}
		}

		for _, r := range rotations {
			newFacing := (r + item.Direction + 4) % 4
			newCost := item.Cost + turnCost

			if cosInfo, ok := cost[setItem{item.Position, newFacing}]; !ok || newCost < cosInfo.cost {
				cost[setItem{item.Position, newFacing}] = costInfo{newCost, 1}
				pred[setItem{item.Position, newFacing}] = []setItem{{item.Position, item.Direction}}
				heap.Push(queue, &datastructure.PriorityQueueItem{Position: item.Position, Direction: newFacing, Cost: newCost})
			} else if newCost == cosInfo.cost && !slices.Contains(pred[setItem{item.Position, newFacing}], setItem{item.Position, item.Direction}) {
				cost[setItem{item.Position, newFacing}] = costInfo{newCost, cosInfo.paths + 1}
				pred[setItem{item.Position, newFacing}] = append(pred[setItem{item.Position, newFacing}], setItem{item.Position, item.Direction})
			}
		}

	}

	return backtrack(cost, pred, end)
}

func backtrack(cost map[setItem]costInfo, pred map[setItem][]setItem, end [2]int) util.HashSet {
	visited := util.NewHashSet()

	minCost := math.MaxInt

	for direction := 0; direction < 4; direction++ {
		endState := setItem{pos: end, dir: direction}
		if cInfo, ok := cost[endState]; ok && cInfo.cost < minCost {
			minCost = cInfo.cost
		}
	}

	for direction := 0; direction < 4; direction++ {
		endState := setItem{pos: end, dir: direction}
		if cInfo, ok := cost[endState]; ok && cInfo.cost == minCost {
			backtrackHelper(endState, visited, pred, cost, minCost)
		}
	}

	toReturn := util.NewHashSet()

	for _, seen := range visited.ToSlice() {
		node := seen.(setItem)
		toReturn.Add(node.pos)
	}

	return *toReturn
}

func backtrackHelper(state setItem, visited *util.HashSet, pred map[setItem][]setItem, cost map[setItem]costInfo, minCost int) {

	if cost[state].cost > minCost {
		return
	}

	if visited.Contains(state) {
		return
	}
	visited.Add(state)

	if predecessors, ok := pred[state]; ok {
		for _, predecessor := range predecessors {
			backtrackHelper(predecessor, visited, pred, cost, minCost)
		}
	}
}

func dijkstra(grid [][]string, start [2]int, end [2]int) int {
	rows, cols := len(grid), len(grid[0])

	queue := datastructure.NewPriorityQueue()
	visited := util.NewHashSet()

	heap.Push(queue, &datastructure.PriorityQueueItem{Position: start, Direction: 0, Cost: 0})

	for queue.Len() > 0 {
		item := heap.Pop(queue).(*datastructure.PriorityQueueItem)

		if item.Position == end {
			return item.Cost
		}

		if visited.Contains(setItem{item.Position, item.Direction}) {
			continue
		}

		visited.Add(setItem{item.Position, item.Direction})
		newPos := [2]int{item.Position[0] + directions[item.Direction][0], item.Position[1] + directions[item.Direction][1]}

		if 0 <= newPos[0] && newPos[0] < rows && 0 <= newPos[1] && newPos[1] < cols && grid[newPos[0]][newPos[1]] != wall {
			heap.Push(queue, &datastructure.PriorityQueueItem{Position: newPos, Direction: item.Direction, Cost: item.Cost + 1})
		}

		for _, r := range rotations {
			newFacing := (item.Direction + r + len(directions)) % len(directions)
			heap.Push(queue, &datastructure.PriorityQueueItem{Position: item.Position, Direction: newFacing, Cost: item.Cost + turnCost})
		}

	}
	return -1
}

func findPOS(grid [][]string, s string) [2]int {
	rows := len(grid)
	cols := len(grid[0])

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == s {
				return [2]int{r, c}
			}
		}
	}
	return [2]int{-1, -1}
}

func addVisitedSpots(grid [][]string, seen util.HashSet) {
	for _, s := range seen.ToSlice() {
		step := s.([2]int)
		grid[step[0]][step[1]] = "O"
	}
}

func printGrid(grid [][]string) {
	for _, row := range grid {
		fmt.Println(strings.Join(row, ""))
	}
}

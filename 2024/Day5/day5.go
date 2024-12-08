package main

import (
	util "2024/until"
	"container/list"
	"fmt"
	"strings"
	"sync"
)

/*
	A further optimization would be to process the updates in one go, then pass the sorted lists to part1 and two where
	the first list is valid updates, then the second is the topological sort since those are what's needed for part2

	Advent of Code Day 5: The solution was fun to do, finally got to implement a topological sort which was fun to do
*/

func main() {
	rawInput := util.ReadInput(util.Parameter())
	order, updates := separateData(rawInput)
	graph := makeGraph(order)
	convertedUpdates := convertUpdates(updates)
	part1(graph, convertedUpdates)
	part2(graph, convertedUpdates)
}

// separateData separates the order (page ranks) from the updates
func separateData(rawInput []string) (order, updates []string) {

	index := 0
	for index < len(rawInput) {
		if strings.Contains(rawInput[index], "|") {
			order = append(order, rawInput[index])
		} else {
			break
		}
		index++
	}

	updates = append(updates, rawInput[index+1:]...)
	return
}

// part1 Take an update and perform the topological sort, if its a valid update (matches sort results)
// then keep the update
func part1(graph map[int][]int, updates [][]int) {
	var wg sync.WaitGroup
	results := make(chan []int, len(updates))

	for _, update := range updates {
		wg.Add(1)
		go func(update []int) {
			defer wg.Done()
			topSort := topologicalSort(update, graph)
			if validateUpdate(update, topSort) {
				results <- update
			}
		}(update)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for result := range results {
		sum += result[len(result)/2]
	}

	fmt.Println("(part 1) Ans: ", sum)

}

// part2 does a very similar process to part1, expect in situations where the update is not valid, then keep the sort results
func part2(graph map[int][]int, updates [][]int) {

	var wg sync.WaitGroup
	results := make(chan []int, len(updates))

	for _, update := range updates {
		wg.Add(1)
		go func(update []int) {
			defer wg.Done()
			topSort := topologicalSort(update, graph)
			if !validateUpdate(update, topSort) {
				results <- topSort
			}
		}(update)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for result := range results {
		sum += result[len(result)/2]
	}

	fmt.Println("(part 2) Ans: ", sum)

}

// makeGraph creates the overall graph out of the provided order
func makeGraph(order []string) map[int][]int {
	graph := make(map[int][]int)
	for _, vertex := range order {
		parts := strings.Split(vertex, "|")
		from := util.ParseInt(parts[0])
		to := util.ParseInt(parts[1])
		graph[from] = append(graph[from], to)
	}
	return graph
}

// topologicalSort performs a topological sort on by creating a subgraph of the provided update
// sort uses Khan's Algorithm
func topologicalSort(update []int, graph map[int][]int) []int {
	// Build subgraph for current update
	subgraph := make(map[int][]int)
	inDegree := make(map[int]int)
	nodesInUpdate := make(map[int]bool)

	// Initialize in-degree and subgraph with nodes that are in the update
	for _, node := range update {
		nodesInUpdate[node] = true
		subgraph[node] = make([]int, 0)
		inDegree[node] = 0
	}

	// Populate the subgraph and in-degree map
	for node := range nodesInUpdate {
		for _, neighbor := range graph[node] {
			if nodesInUpdate[neighbor] {
				subgraph[node] = append(subgraph[node], neighbor)
				inDegree[neighbor]++
			}
		}
	}

	// Perform topological sort
	var topOrder []int
	queue := list.New()
	// add any nodes where the degree was initially 0, indicates they should be first in the order
	for node, degree := range inDegree {
		if degree == 0 {
			queue.PushBack(node)
		}
	}

	// Looping through the queue, remove element from queue, and add to topOrder, this means that it's degrees remaining is zero
	for queue.Len() > 0 {
		element := queue.Front()
		node := element.Value.(int)
		queue.Remove(element)
		topOrder = append(topOrder, node)

		// decrement the neighbors degrees, and then add to queue if degree hits zero
		for _, neighbor := range subgraph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue.PushBack(neighbor)
			}
		}
	}

	return topOrder
}

// validateUpdate checks the update against its calculated topological sort
// If the position of value is out of order compared to the sort results, then we know its incorrect
func validateUpdate(update []int, topOrder []int) bool {
	// After sort, if the len of topological order doesn't match update, then cycle is present and it is not valid
	if len(topOrder) != len(update) {
		return false
	}

	// Map pages to their position in the topological order
	position := make(map[int]int)
	for i, page := range topOrder {
		position[page] = i
	}

	// Check if the update order respects the topological order
	for i := 0; i < len(update)-1; i++ {
		if position[update[i]] > position[update[i+1]] {
			return false
		}
	}
	return true
}

// convertUpdates converts the string input into a usable format
func convertUpdates(updateStr []string) [][]int {
	toReturn := make([][]int, 0)
	for _, update := range updateStr {
		converted := make([]int, 0)
		split := strings.Split(update, ",")
		for _, value := range split {
			converted = append(converted, util.ParseInt(value))
		}
		toReturn = append(toReturn, converted)
	}
	return toReturn
}

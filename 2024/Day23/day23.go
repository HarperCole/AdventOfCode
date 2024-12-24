package main

import (
	"2024/util"
	"fmt"
	"slices"
	"strings"
)

func main() {
	graph := makeGraph(util.ReadInput(util.Parameter()))
	part1(graph)
	part2(graph)

}

// part1 processes the graph to find all triangles (three nodes that are all connected to each other)
// and counts how many of these triangles contain at least one node that starts with the letter 't'.
// It prints the count of such triangles.
func part1(graph map[string][]string) {

	// Create a map to store triangles, where the key is a sorted string of the three nodes
	// and the value is a slice of the three nodes.
	triangles := make(map[string][]string)

	// Iterate over each computer and its neighbors in the graph.
	for computer, neighbors := range graph {
		// For each neighbor of the current computer.
		for _, neighbor := range neighbors {
			// Check for a third node that forms a triangle with the current computer and neighbor.
			for _, n := range graph[computer] {
				// Ensure the third node is not the same as the neighbor and is connected to the neighbor.
				if n != neighbor && slices.Contains(graph[neighbor], n) {
					// Add the triangle to the map.
					triangles[triangleKey([]string{computer, neighbor, n})] = []string{computer, neighbor, n}
				}
			}
		}
	}

	// Initialize a counter for triangles that contain at least one node starting with 't'.
	countThatContainLetterT := 0

	// Iterate over each triangle in the map.
	for _, t := range triangles {
		// Check each node in the triangle.
		for _, computer := range t {
			// If a node starts with 't', increment the counter and break out of the loop.
			if computer[0] == 't' {
				countThatContainLetterT++
				break
			}
		}
	}

	// Print the count of triangles that contain at least one node starting with 't'.
	fmt.Println("(part 1) Ans: ", countThatContainLetterT)
}

// part2 finds the largest clique (a subset of nodes where every two nodes are connected) in the graph
// and prints the nodes in the largest clique in sorted order.
//
// Parameters:
//
//	graph - a map where the key is a node and the value is a slice of nodes connected to the key node.
func part2(graph map[string][]string) {
	var largestClique []string

	// Iterate over each computer and its neighbors in the graph.
	for computer, neighbors := range graph {
		clique := []string{computer}
		// Check each neighbor to see if it can be added to the current clique.
		for _, neighbor := range neighbors {
			if isFullyConnected(graph, clique, neighbor) {
				clique = append(clique, neighbor)
			}
		}

		// Update the largest clique if the current clique is larger.
		if len(clique) > len(largestClique) {
			largestClique = clique
		}
	}

	// Sort the largest clique and print the result.
	slices.Sort(largestClique)
	fmt.Println("(part 2) Ans: ", strings.Join(largestClique, ","))
}

// isFullyConnected checks if a candidate node can be added to a clique such that all nodes in the clique
// remain fully connected to each other.
//
// Parameters:
//
//	graph - a map where the key is a node and the value is a slice of nodes connected to the key node.
//	clique - a slice of nodes representing the current clique.
//	candidate - a node to be checked for full connectivity with the clique.
//
// Returns:
//
//	true if the candidate node is fully connected to all nodes in the clique, false otherwise.
func isFullyConnected(graph map[string][]string, clique []string, candidate string) bool {
	for _, member := range clique {
		if !slices.Contains(graph[member], candidate) {
			return false
		}
	}
	return true
}

func makeGraph(input []string) map[string][]string {
	graph := make(map[string][]string)

	for _, line := range input {
		computers := strings.Split(line, "-")
		computerOne := computers[0]
		computerTwo := computers[1]
		graph[computerOne] = append(graph[computerOne], computerTwo)
		graph[computerTwo] = append(graph[computerTwo], computerOne)

	}
	return graph

}

func triangleKey(key []string) string {
	slices.Sort(key)
	return strings.Join(key, ",")
}

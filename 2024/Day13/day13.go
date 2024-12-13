package main

import (
	util "2024/until"
	"fmt"
	"math"
	"regexp"
	"strings"
)

const (
	patternA = `Button A: X\+(\d+), Y\+(\d+)\s*`
	patternB = `Button B: X\+(\d+), Y\+(\d+)\s*`
	prize    = `Prize: X=(\d+), Y=(\d+)`
	attempts = 100
	bigPrize = 10000000000000
)

func main() {
	input := util.ReadInput(util.Parameter())
	machines := captureMachines(strings.Join(input, "\n"))
	part1(machines)
	part2(machines)
}

// part1 solves the problem for the first part where attempts are limited to 100.
// It calculates the minimum total cost for all machines if solutions are possible.
func part1(machines [][6]int) {
	minCost := 0
	solutions := make([][2]int, 0)
	for _, m := range machines {
		minCostFound, solution := solveLinearCombinations(m[0], m[1], m[2], m[3], m[4], m[5])
		if minCostFound != math.MaxInt64 {
			minCost += minCostFound
			solutions = append(solutions, solution)
		}
	}

	fmt.Println("(part 1) Ans :", minCost)
}

// solveLinearCombinations attempts to solve the claw machine problem for a single machine
// by iterating through a range of attempts to find a valid combination of button presses.
// It returns the minimum cost and the solution (A presses, B presses).
func solveLinearCombinations(xA int, yA int, xB int, yB int, xP int, yP int) (int, [2]int) {
	minCost := math.MaxInt64
	bestSolution := [2]int{-1, -1}

	for attempt := 0; attempt <= attempts; attempt++ {

		if (xP-xA*attempt)%xB != 0 {
			continue
		}

		bX := (xP - xA*attempt) / xB

		if (yP-yA*attempt)%yB != 0 {
			continue
		}

		bY := (yP - yA*attempt) / yB

		if bX == bY && bX >= 0 {
			cost := 3*attempt + bX
			if cost < minCost {
				minCost = cost
				bestSolution = [2]int{attempt, bX}
			}
		}

	}
	return minCost, bestSolution
}

// part2 solves the problem for the second part where prize positions are greatly increased.
// It calculates the minimum total cost for all machines with no restriction on attempts.
func part2(machines [][6]int) {
	var minCost int64 = 0

	for _, m := range machines {
		cost := solveLargerPrize(m[0], m[1], m[2], m[3], int64(m[4]+bigPrize), int64(m[5]+bigPrize))
		if cost != -1 {
			minCost += cost

		}
	}

	fmt.Println("(part 2) Ans :", minCost)
}

// solveLargerPrize solves the claw machine problem for a single machine with larger prize positions.
// It calculates the cost directly using linear algebra and floating-point arithmetic.
// Returns the minimum cost or -1 if no solution exists.
func solveLargerPrize(xA, yA, xB, yB int, pX, pY int64) int64 {
	fax := float64(xA)
	fay := float64(yA)
	fbx := float64(xB)
	fby := float64(yB)
	fpx := float64(pX)
	fpy := float64(pY)

	n := (fpx*fby - fpy*fbx) / (fax*fby - fay*fbx)
	m := (fpy - n*fay) / fby

	if n == math.Trunc(n) && m == math.Trunc(m) {
		return int64(3*n + m)
	}

	return int64(-1)
}

// captureMachines parses the raw input string to extract machine configurations.
// Each machine's configuration is represented as a 6-element array of integers
// containing coefficients for A, B, and prize positions.
func captureMachines(input string) [][6]int {
	machines := make([][6]int, 0)

	pattern := regexp.MustCompile(patternA + patternB + prize)

	matches := pattern.FindAllStringSubmatch(input, -1)

	for _, m := range matches {
		var values [6]int
		for i, v := range m[1:] {
			values[i] = util.ParseInt(v)
		}
		machines = append(machines, values)
	}

	return machines
}

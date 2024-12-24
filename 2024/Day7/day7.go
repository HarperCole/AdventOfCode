package main

import (
	"2024/util"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// equations is struct where:
//   - target: target number for values to possibly equal
//   - values: slice of numbers that calculations will be performed on to see if it matches target
type equations struct {
	target int
	values []int
}

/*
	Advent of Code Day 7
	Part 1: Looking at this problem it was either implement recursive backtracking, or because it's a situation where you can break it down into sub-problems, use Dynamic programming, which is what I did.
			Very simple sub-problem as well, using a table where each number in the equation list gets an entry, then for the current table entry calculate and add the sum and product of each previous entry with the current number
			to the table. Then, once the entire table is calculated, the last entry of the dp table will have all possible combinations of sum and products for each number in equation list and target will be there if it was possible with the
			numbers available.
	Part 2: Took code from part 1  (It would be a lot of work to refactor to compartmentalize reusable code, and there is enough of a distinction from the part 1 code that I felt it was fine). Added a third operation to dp, which was concatenation
			Which is to put the two number together like 5 || 6 => 56. Just converted both elements to a string, concatenated them, and parsed the int out.

	Since I chose to use dynamic programming, I was able to keep it non-recursive so felt much better about making this solution concurrent, where each possible equation is calculated on its own goroutine
*/

func main() {
	inpput := util.ReadInput(util.Parameter())
	equationList := parseEquations(inpput)
	part1(equationList)
	part2(equationList)
}

// part2 is the dynamic programming solution, uses the same flow as part1 but adds the third operator
func part2(eqs []equations) {
	var wg sync.WaitGroup
	results := make(chan int, len(eqs))

	for _, equation := range eqs {
		wg.Add(1)
		go func(equation equations, results chan int, wg *sync.WaitGroup) {
			defer wg.Done()
			target := equation.target

			dp := make([]util.HashSet, len(equation.values))

			for index := 0; index < len(equation.values); index++ {
				dp[index] = *util.NewHashSet()
			}

			dp[0].Add(equation.values[0])

			for i := 1; i < len(equation.values); i++ {
				currentTable := dp[i-1].ToSlice()
				for _, val := range currentTable {
					dp[i].Add(val.(int) + equation.values[i])
					dp[i].Add(val.(int) * equation.values[i])
					// To concatenation
					leftPart := strconv.Itoa(val.(int))
					rightPart := strconv.Itoa(equation.values[i])
					dp[i].Add(util.ParseInt(leftPart + rightPart))
				}
			}

			if dp[len(equation.values)-1].Contains(target) {
				results <- target
			}
		}(equation, results, &wg)

	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for result := range results {
		sum += result
	}

	fmt.Println("(part 2) Ans: ", sum)
}

// part1 solves the prompt using dynamic programming, read above for thought process
func part1(eqs []equations) {

	var wg sync.WaitGroup
	results := make(chan int, len(eqs))

	for _, equation := range eqs {
		wg.Add(1)
		go func(equation equations, results chan int, wg *sync.WaitGroup) {
			defer wg.Done()
			target := equation.target

			// create dp table
			dp := make([]util.HashSet, len(equation.values))

			// go requires you to explicitly initialize every index
			for index := 0; index < len(equation.values); index++ {
				dp[index] = *util.NewHashSet()
			}

			// populate the first number into the equation
			dp[0].Add(equation.values[0])

			// for each number after 0, we are going to find the product and sum with every other number
			// that exists in the previous table entry
			for i := 1; i < len(equation.values); i++ {
				currentTable := dp[i-1].ToSlice()
				for _, val := range currentTable {
					dp[i].Add(val.(int) + equation.values[i])
					dp[i].Add(val.(int) * equation.values[i])
				}
			}

			// check last dp table entry to see if target is present, if it is, then send target on channel
			if dp[len(equation.values)-1].Contains(target) {
				results <- target
			}
		}(equation, results, &wg)

	}

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for result := range results {
		sum += result
	}

	fmt.Println("(part 1) Ans: ", sum)

}

// parseEquations extracts the target val and list of numbers that potentially equate to target
func parseEquations(input []string) []equations {

	toReturn := make([]equations, 0)
	for _, line := range input {
		split := strings.Split(line, " ")
		targetStr := split[0]
		targetStr = targetStr[:len(targetStr)-1]

		target := util.ParseInt(targetStr)
		values := make([]int, 0)
		for i := 1; i < len(split); i++ {
			values = append(values, util.ParseInt(split[i]))
		}
		toReturn = append(toReturn, equations{target: target, values: values})
	}
	return toReturn
}

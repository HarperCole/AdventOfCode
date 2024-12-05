package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
	Advent of Code Day 3
	This problem was a straight forward parsing problem where I got to practice using regex.
	Just created regex patterns and matched and filtered the results
	pattern matched for mul({1-3},{1-3}) in part 1
	pattern matched for mul({1-3},{1-3}) do() don't()
*/

const (
	regexPatternPart1 = `(mul\(\d{1,3},\d{1,3}\))`
	numbersPatter     = `\d{1,3},\d{1,3}`
	regexPatternPart2 = `mul\(\d{1,3},\d{1,3}\)|do\(\)|don't\(\)`
	do                = "do()"
	dont              = "don't()"
)

// readInput read in input, want this raw data as a string
func readInput(fname string) string {
	fmt.Println("Reading input file", fname)
	file, _ := os.Open(fname)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	input := make([]string, 0)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	return strings.Join(input, "")
}

func main() {
	filename := flag.String("file", "", "Input filename")
	flag.Parse()

	if *filename == "" {
		fmt.Println("No file name was provided")
		os.Exit(1)
	}

	input := readInput(*filename)
	part1(input)
	part2(input)

}

// matchPattern takes raw input, then uses the regex Pattern that was passed in to create matches, filter for do and dont if filter is true
// Finally extract the number pairs
//   - pattern - input data
//   - regexPattern - the regex pattern used to find groups
//   - filter - used to determine if we want to filter out some groups (has to do with part 2)
func matchPattern(pattern string, regexPattern string, filter bool) []string {
	re := regexp.MustCompile(regexPattern)
	numRe := regexp.MustCompile(numbersPatter)
	matches := re.FindAllString(pattern, -1)
	if filter {
		validResults := make([]string, 0)
		skip := false
		for _, m := range matches {
			if m == dont {
				skip = true
			} else if m == do {
				skip = false
			}
			if m != do && m != dont && !skip {
				validResults = append(validResults, m)
			}

		}
		matches = validResults
	}
	extractedNumbers := make([]string, 0)
	for _, match := range matches {
		extractedNumbers = append(extractedNumbers, numRe.FindAllString(match, -1)...)
	}
	return extractedNumbers
}

// extractNumbers transforms the extracted number pairs into integers
// Since Golang doesn't have tuples, just using a slice of ints of size 2
func extractNumbers(numsInStr []string) [][]int {
	results := make([][]int, 0)
	for _, numStr := range numsInStr {
		pair := make([]int, 2)
		numbers := strings.Split(numStr, ",")
		pair[0], _ = strconv.Atoi(numbers[0])
		pair[1], _ = strconv.Atoi(numbers[1])
		results = append(results, pair)
	}
	return results
}

// part1 is the solution for part1 of day 3 chains together calls to find the solution:
// match all instances of mul(X,Y) and sum these products, where X,Y are 1 - 3 digit numbers
func part1(input string) {
	extractedPairs := extractNumbers(matchPattern(input, regexPatternPart1, false))
	sum := findProductSum(extractedPairs)
	fmt.Println("(part 1) Ans: ", sum)
}

// part2 is the solution for part2 of day 3 chains together calls to find the solution:
// matches all instances of mul(X,Y) expect with these requirements:
//   - The do() instruction enables future mul instructions.
//   - The don't() instruction disables future mul instructions.
//   - Only the most recent do() or don't() instruction applies. At the beginning of the program, mul instructions are enabled.
func part2(input string) {
	extractedPairs := extractNumbers(matchPattern(input, regexPatternPart2, true))
	sum := findProductSum(extractedPairs)
	fmt.Println("(part 2) Ans: ", sum)
}

// findProductSum simple function to find sum of the product of the pairs
func findProductSum(extractedPairs [][]int) int {
	sum := 0
	for _, pair := range extractedPairs {
		sum += pair[0] * pair[1]
	}
	return sum
}

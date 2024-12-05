package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Day 1 Solutions to Advent of Code 2024
// Part 1 was to sort the two columns in ascending order, then sum the difference of each row
// Part 2 was to create a frequency map out of the right column, then using the left column values, search the frequency map for value then multiple left_value * frequency. If left value is not present in map, then set to 0.

// Wrote this using goroutines to practice writing concurrent go code, it is very overkill for this problem

func main() {
	rawData := readInData("Day1/part1Data.txt")
	leftValues, rightValues := convertInputData(rawData)
	part1(leftValues, rightValues)
	part2(leftValues, rightValues)
}

// readInData takes a file and reads each line in
func readInData(filename string) []string {
	fmt.Println("Reading file ", filename)
	data, _ := os.Open(filename)
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	data.Close()
	return fileLines
}

// convertInputData takes the raw data that was read into the file and turns it into a usable format for the problem
func convertInputData(input_data []string) (leftValues []int, rightValues []int) {
	fmt.Println("Converting raw Data")
	for _, line := range input_data {
		fields := strings.Fields(line)
		left, _ := strconv.Atoi(fields[0])
		right, _ := strconv.Atoi(fields[1])

		leftValues = append(leftValues, left)
		rightValues = append(rightValues, right)
	}
	return
}

// part1 Had to sort the two slices in ascending order then find the difference values at same indices and then sum the differences
// Example: [1,2,3] [2.3.4] -> 1 + 1 + 1 -> 3
// I wrote it concurrently to get practice writing concurrent go code (Very overkill)
func part1(leftValues []int, rightValues []int) {
	sort.Ints(leftValues)
	sort.Ints(rightValues)
	result := make(chan int, len(leftValues))
	var wg sync.WaitGroup

	for index, _ := range leftValues {
		wg.Add(1)
		go func(left, right int, result chan<- int, wg *sync.WaitGroup) {
			defer wg.Done()
			result <- int(math.Abs(float64(left - right)))
		}(leftValues[index], rightValues[index], result, &wg)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	answer := 0
	for resultValue := range result {
		answer += resultValue
	}

	fmt.Println("part1 answer is", answer)
}

// part2 multiple values in the leftValue slice by their frequency in the second map
// Example [1, 2, 3] [1,1,2] -> (1 * 2) + (2 * 1) + (3 * 0) -> 3
// Using a frequency map created out of the right values, iterate through left values and calculate "similarity score"
// with the frequency map
func part2(leftValues []int, rightValues []int) {
	frequencyMap := make(map[int]int)
	populateFrequencyMap(frequencyMap, rightValues)

	result := make(chan int, len(leftValues))
	var wg sync.WaitGroup

	for _, value := range leftValues {
		wg.Add(1)

		go func(value int, result chan<- int, wg *sync.WaitGroup) {
			defer wg.Done()

			toAdd := 0
			if frequency, ok := frequencyMap[value]; ok {
				toAdd = value * frequency
			}
			result <- toAdd

		}(value, result, &wg)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	ans := 0
	for resultValue := range result {
		ans += resultValue
	}

	fmt.Println("part2 answer is", ans)
}

// populateFrequencyMap used to create a frequency map out of a int slice
// value -> frequency of occurrence within the slice
func populateFrequencyMap(frequency_map map[int]int, values []int) {
	for _, value := range values {
		frequency_map[value]++
	}
}

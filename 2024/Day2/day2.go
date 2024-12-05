package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
	Advent of Code 2024 Day 2 Solutions
	Kept the same theme of trying to solve these concurrently. I felt like I did a good job with part 1, it makes sense.
	However, once I saw part 2 I realized that doing it concurrently might not be the best approach, seeing as you had
	to search for a possible permutation where the int slice would become valid. One concern would be too many goroutines
	being produced. I think this could have been simplified if I didn't use wait groups in part 2, however because
	checkSafetyLevels required a wait group I feel like I was forced to use them.
*/

func main() {
	filename := "2024/Day2/day2Data.txt"
	//filename := "2024/Day2/example_data.txt"
	data := convertData(readInData(filename))
	part1(data)
	part2(data)

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

// convertData takes in the raw data and turns it into a slice of int slices
// a list of reports
func convertData(rawData []string) [][]int {
	fmt.Println("Converting data ")
	toReturn := make([][]int, 0)
	for _, line := range rawData {
		row := make([]int, 0)
		for _, number := range strings.Fields(line) {
			convertedNumber, _ := strconv.Atoi(number)
			row = append(row, convertedNumber)
		}
		toReturn = append(toReturn, row)
	}
	return toReturn
}

// part1 Requirements:
// Must be all increasing or decreasing
// difference between adjacent levels are considered safe if differ by at least 1 or at most 3
func part1(lines [][]int) {
	safeChannel := make(chan bool)
	var wg sync.WaitGroup

	for _, line := range lines {
		wg.Add(1)
		go checkSafetyLevels(line, safeChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(safeChannel)
	}()

	safeReports := 0
	for safeReport := range safeChannel {
		if safeReport {
			safeReports++
		}
	}

	fmt.Println("(part1) Safe reports found: ", safeReports)
}

// part2 Requirements:
// Must be all increasing or decreasing
// difference between adjacent levels are considered safe if differ by at least 1 or at most 3
// can skip one bad level in each report
func part2(lines [][]int) {

	safeChannel := make(chan bool)
	var wg sync.WaitGroup

	for _, line := range lines {
		wg.Add(1)
		go func(line []int, safeChannel chan bool, wg *sync.WaitGroup) {
			defer wg.Done()
			var waitGroup sync.WaitGroup
			individualChannel := make(chan bool, 1)

			waitGroup.Add(1)
			go checkSafetyLevels(line, individualChannel, &waitGroup)
			waitGroup.Wait()
			close(individualChannel)
			for result := range individualChannel {
				if result {
					safeChannel <- true
					return
				}
			}

			// remove each index one by one until either we hit a state where it becomes true or just send a false

			for index, _ := range line {
				waitGroup.Add(1)
				dampeningChannel := make(chan bool, 1)
				go checkSafetyLevels(remove(line, index), dampeningChannel, &waitGroup)
				waitGroup.Wait()
				close(dampeningChannel)
				for result := range dampeningChannel {
					if result {
						safeChannel <- true
						return
					}
				}
			}

			safeChannel <- false

		}(line, safeChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(safeChannel)
	}()

	safeReports := 0
	for safeReport := range safeChannel {
		if safeReport {
			safeReports++
		}
	}

	fmt.Println("(part2) Safe reports found: ", safeReports)
}

// checkSafetyLevels is the function that checks each report based on the requirements of part1
// slice must be A) ascending or descending B) differences between each level (index) cannot be less than 1 or greater than 3

func checkSafetyLevels(line []int, safeChannel chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(line) < 2 {
		safeChannel <- true
		return
	}

	trailer := line[0]
	ascending := trailer < line[1]

	for i := 1; i < len(line); i++ {

		diff := trailer - line[i]
		if diff < 0 {
			diff = -diff
		}
		stillAscending := trailer < line[i]

		if diff < 1 || diff > 3 || ascending != stillAscending {
			safeChannel <- false
			return
		}
		trailer = line[i]
		ascending = stillAscending
	}
	safeChannel <- true
}

// remove is a simple function that makes a copy of the input slice and removes an index
func remove(slice []int, index int) []int {
	// Create a new slice to ensure we don't modify the original
	result := make([]int, 0, len(slice)-1)
	result = append(result, slice[:index]...)
	result = append(result, slice[index+1:]...)
	return result
}

package main

import (
	"2024/util"
	"fmt"
	"strings"
)

func main() {
	locks, keys := parseData(util.ReadInput(util.Parameter()))
	part1(locks, keys)
}

func part1(locks [][]int, keys [][]int) {
	combos := 0

	for _, lock := range locks {
		for _, key := range keys {
			if canUnlock(lock, key) {
				combos++
			}
		}
	}

	fmt.Println("(part 1) Ans: ", combos)
}

func canUnlock(lock []int, key []int) bool {

	for i := 0; i < len(lock); i++ {
		if lock[i]+key[i] > 5 {
			return false
		}
	}

	return true
}

func parseData(input []string) ([][]int, [][]int) {
	locks := make([][]int, 0)
	keys := make([][]int, 0)

	rawInputs := make([][]string, 0)
	rawInput := make([]string, 0)
	for _, line := range input {
		if line == "" {
			rawInputs = append(rawInputs, rawInput)
			rawInput = make([]string, 0)
			continue
		}
		rawInput = append(rawInput, line)
	}
	rawInputs = append(rawInputs, rawInput)

	for _, rawInput := range rawInputs {
		if strings.Join(strings.Fields(rawInput[0]), "") == "#####" {
			locks = append(locks, parseLockKey(rawInput, true))
		} else {
			keys = append(keys, parseLockKey(rawInput, false))
		}
	}
	return locks, keys
}

func parseLockKey(input []string, islock bool) []int {
	rows := len(input)
	cols := len(input[0])
	lock := make([]int, 0)

	for c := 0; c < cols; c++ {
		size := 0

		if islock {
			for r := 1; r < rows-1; r++ {
				if input[r][c] == '#' {
					size++
				} else if input[r][c] == '.' {
					break
				}
			}
		} else {
			for r := rows - 2; r > 0; r-- {
				if input[r][c] == '#' {
					size++
				} else if input[r][c] == '.' {
					break
				}
			}
		}

		lock = append(lock, size)
	}
	return lock
}

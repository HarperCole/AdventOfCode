package main

import (
	util "2024/until"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const (
	part1Rounds = 25
	part2Rounds = 75
	ruleThree   = 2024
	workers     = 10
)

// Task represents works
//   - stoneNum: starting number on stone
//   - rounds: number of rounds that need to be computed
type Task struct {
	stoneNum int
	rounds   int
}

// Result represents finished work
//   - stoneNum: original stone number
//	 - rounds: number of rounds that were computed
//   - count: number of stones created from original stone and rounds worked
type Result struct {
	stoneNum int
	rounds   int
	count    int
}

func main() {
	input := convertInt(util.ReadInput(util.Parameter()))
	part1(input)
	part2(input)
}

func convertInt(input []string) []int {
	toReturn := make([]int, 0)
	for _, number := range strings.Split(input[0], " ") {
		toReturn = append(toReturn, util.ParseInt(number))
	}
	return toReturn
}

// part1 solves part 1, creates 10 workers, and creates the tasks out of input data
// for each stone, part1Rounds are computed
func part1(input []int) {
	memo := make(map[string]int)
	var mu sync.Mutex

	// create channels to send tasks to workers, and then collect the results out of the works
	tasks := make(chan Task)
	results := make(chan Result, len(input))
	var wg sync.WaitGroup

	// add worker count to barrier then launch workers
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go worker(tasks, results, memo, &wg, &mu, part1Rounds)
	}

	// create and send tasks to workers
	for _, stone := range input {
		tasks <- Task{stone, part1Rounds}
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	// collect results
	count := 0
	for result := range results {
		count += result.count
	}

	fmt.Println("(Part 1) Ans:", count)
}

// part2 solves part 2, creates 10 workers, and creates the tasks out of input data
// for each stone, part2Rounds are computed
func part2(input []int) {
	memo := make(map[string]int)
	var mu sync.Mutex

	tasks := make(chan Task)
	results := make(chan Result, len(input))
	var wg sync.WaitGroup

	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go worker(tasks, results, memo, &wg, &mu, part2Rounds)
	}

	for _, stone := range input {
		tasks <- Task{stone, part2Rounds}
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for result := range results {
		count += result.count
	}

	fmt.Println("(Part 2) Ans:", count)
}

// countStones is the function that calculates the number of stones created a specific stone for an n number of rounds
func countStones(stone int, round int, memo map[string]int, mu *sync.Mutex) int {
	// terminate when round == 0, and return 1 because no rules are applied
	if round == 0 {
		return 1
	}

	// create key to check memo map (stone, current round) need to store the round because a 12 stone at round 2 will result in a different number than a stone 12 at round 1
	key := fmt.Sprintf("%d,%d", stone, round)
	// check map for key
	mu.Lock()
	if val, exists := memo[key]; exists {
		// if we've seen this calculation, return it
		mu.Unlock()
		return val
	}
	mu.Unlock()

	var result int

	// Rule 1: if a stone is 0, set it as 1 for the next round
	if stone == 0 {
		result = countStones(1, round-1, memo, mu)
		// Rule 2: If a stone is a even length, split in half and continue on the left and right side
	} else if len(strconv.Itoa(stone))%2 == 0 {
		digits := strconv.Itoa(stone)
		left := util.ParseInt(digits[:len(digits)/2])
		right := util.ParseInt(digits[len(digits)/2:])

		leftResult := countStones(left, round-1, memo, mu)
		rightResult := countStones(right, round-1, memo, mu)
		result += leftResult + rightResult
		// Rule 3: If stone length is odd, multiply the by ruleThree
	} else {
		result = countStones(stone*ruleThree, round-1, memo, mu)
	}

	mu.Lock()
	// add result calculated to the memo map
	memo[key] = result
	mu.Unlock()

	return result
}

// worker a function that is designed to complete a task, in this case work on a task
// the workers are created on a goroutine, and the workers share the channel of tasks to work on
func worker(tasks <-chan Task, results chan<- Result, memo map[string]int, wg *sync.WaitGroup, mu *sync.Mutex, rounds int) {
	defer wg.Done()
	for task := range tasks {
		count := countStones(task.stoneNum, task.rounds, memo, mu)
		results <- Result{stoneNum: task.stoneNum, rounds: rounds, count: count}
	}
}

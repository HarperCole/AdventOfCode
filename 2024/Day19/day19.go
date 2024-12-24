package main

import (
	"2024/util"
	"fmt"
	"regexp"
	"slices"
	"sync"
)

const (
	pattern = `([a-zA-Z]+)`
	workers = 10
)

type Task struct {
	word string
}

func main() {
	input := util.ReadInput(util.Parameter())
	splitIndex := slices.Index(input, "")
	wordBank, words := createWordBank(input[:splitIndex]), input[splitIndex+1:]

	part1(wordBank, words)
	part2(wordBank, words)
}

func part2(bank []string, words []string) {
	tasks := make(chan Task)
	results := make(chan int, len(words))
	var wg sync.WaitGroup

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(tasks, results, &wg, bank, longestWord(bank), false)
	}

	for _, word := range words {
		tasks <- Task{word}
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for result := range results {
		count += result
	}

	fmt.Println("(part 2) Ans:", count)
}

func part1(bank []string, words []string) {
	tasks := make(chan Task)
	results := make(chan int, len(words))
	var wg sync.WaitGroup
	longestWordInBank := longestWord(bank)

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(tasks, results, &wg, bank, longestWordInBank, true)
	}

	for _, word := range words {
		tasks <- Task{word}
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for result := range results {
		count += result
	}

	fmt.Println("(part 1) Ans:", count)
}

func longestWord(bank []string) int {
	longest := -1
	for _, word := range bank {
		longest = max(longest, len(word))
	}
	return longest
}

func createWordBank(strings []string) []string {
	reg := regexp.MustCompile(pattern)
	toReturn := make([]string, 0)
	matches := reg.FindAllStringSubmatch(strings[0], -1)
	for _, match := range matches {
		toReturn = append(toReturn, match[0])
	}
	return toReturn
}

// wordBreak checks if the given word can be segmented into a sequence of one or more dictionary words.
// It uses dynamic programming to determine if the word can be broken down using the words in the wordBank.
//
// Parameters:
//   - word: The word to be segmented.
//   - wordBank: A slice of strings representing the dictionary of words.
//   - longestWord: The length of the longest word in the wordBank.
//
// Returns:
//   - An integer (1 or 0) indicating whether the word can be segmented (1) or not (0).
func wordBreak(word string, wordBank []string, longestWord int) int {
	dp := make([]bool, len(word)+1)
	dp[0] = true
	for i := 1; i <= len(word); i++ {
		for j := max(0, i-longestWord); j < i; j++ {
			if dp[j] && slices.Contains(wordBank, word[j:i]) {
				dp[i] = true
				break
			}
		}
	}

	if dp[len(word)] {
		return 1
	} else {
		return 0
	}
}

// combinations calculates the number of ways the given word can be segmented into a sequence of one or more dictionary words.
// It uses dynamic programming to count the possible segmentations.
//
// Parameters:
//   - word: The word to be segmented.
//   - bank: A slice of strings representing the dictionary of words.
//
// Returns:
//   - An integer representing the number of ways the word can be segmented.
func combinations(word string, bank []string) int {
	dp := make([]int, len(word)+1)
	dp[0] = 1

	for i := 1; i <= len(word); i++ {
		for _, w := range bank {
			wordLen := len(w)
			if i >= wordLen && word[i-wordLen:i] == w {
				dp[i] += dp[i-wordLen]
			}
		}
	}

	return dp[len(word)]
}

func worker(tasks <-chan Task, results chan<- int, group *sync.WaitGroup, wordBank []string, longestWord int, partOne bool) {
	defer group.Done()
	for task := range tasks {
		if partOne {
			results <- wordBreak(task.word, wordBank, longestWord)
		} else {
			results <- combinations(task.word, wordBank)
		}
	}
}

package main

import (
	util "2024/until"
	"fmt"
	"sync"
)

const (
	pruneNum = 16777216
)

type sequence struct {
	zero  int
	one   int
	two   int
	three int
}

func main() {
	input := util.StringToIntSlice(util.ReadInput(util.Parameter()))
	part1(input)
	part2(input)
}

func part1(input []int) {

	results := make(chan int, len(input))
	var waitGroup sync.WaitGroup

	for _, secretNum := range input {
		waitGroup.Add(1)
		go calculateAllSecretNum(&waitGroup, secretNum, results, 2000)
	}

	go func() {
		waitGroup.Wait()
		close(results)
	}()

	sum := 0
	for result := range results {
		sum += result
	}

	fmt.Println("(part 1) Ans: ", sum)
}

func part2(input []int) {
	secMap := make(map[int]int)
	seqNumbers := make(map[sequence]int)

	for _, secretNum := range input {
		currSeqMap := determineSequences(secretNum, &secMap, 2000)
		for key, value := range currSeqMap {
			seqNumbers[key] += value
		}
	}

	maxVal := 0
	for _, val := range seqNumbers {
		if val > maxVal {
			maxVal = val
		}
	}

	fmt.Println("(part 2) Ans: ", maxVal)
}

func determineSequences(secretNum int, cache *map[int]int, times int) map[sequence]int {
	localMap := make(map[sequence]int)
	vals := []int{secretNum % 10}

	// Perform initial sequence calculation
	for i := 0; i < 4; i++ {
		secretNum = calculateSingleSecretNum(secretNum, 1, cache)
		vals = append(vals, secretNum%10)
		times--
	}
	localMap[sequence{vals[1] - vals[0], vals[2] - vals[1], vals[3] - vals[2], vals[4] - vals[3]}] = vals[4]

	// Process remaining times iteratively
	for times > 1 {
		secretNum = calculateSingleSecretNum(secretNum, 1, cache)
		vals = vals[1:] // Shift window
		vals = append(vals, secretNum%10)
		currentSequence := sequence{vals[1] - vals[0], vals[2] - vals[1], vals[3] - vals[2], vals[4] - vals[3]}
		_, ok := localMap[currentSequence]
		if !ok {
			localMap[currentSequence] = vals[4]
		}
		times--
	}
	return localMap
}

func calculateSingleSecretNum(secretNum int, times int, secMap *map[int]int) int {
	for times > 0 {
		startNum := secretNum
		if solution, ok := (*secMap)[startNum]; ok {
			secretNum = solution
		} else {
			secretNum = calculateSecretNum(startNum)
			(*secMap)[startNum] = secretNum
		}
		times--
	}
	return secretNum
}

func calculateAllSecretNum(waitGroup *sync.WaitGroup, secretNum int, results chan int, times int) {
	defer waitGroup.Done()
	sn := secretNum
	for i := 0; i < times; i++ {
		sn = calculateSecretNum(sn)
	}
	//fmt.Println("SecretNum: ", sn)
	results <- sn
}

func calculateSecretNum(secretNum int) int {
	sn := secretNum
	result := sn * 64
	sn = mix(sn, result)
	sn = prune(sn)
	result = sn / 32
	sn = mix(sn, result)
	sn = prune(sn)
	result = sn * 2048
	sn = mix(sn, result)
	sn = prune(sn)
	return sn
}

func prune(secretNum int) int {
	return secretNum % pruneNum
}

func mix(secretNum, result int) int {
	return secretNum ^ result
}

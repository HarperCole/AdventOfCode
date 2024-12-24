package main

import (
	"2024/util"
	"fmt"
	"strings"
)

// indexSpace is a struct representation of a decompressed position in the disk
//   - free indicates if this position in the disk is occupied or not
//   - fileNum indicates the file num of the file that occupies this space
type indexSpace struct {
	free    bool
	fileNum int
}

/*
	Advent of Code Day 9:
		Part 1: After compressing it, it was a two pointer solution to swap the left most free position with the right most filled spot
		Part 2: This took a moment for me to realize that the way I structured my decompressed file map, I first needed to group the file, then search for a span of free space. Once I realized that, it was pretty straight forward

   I think I could have done without using the indexSpace block and instead just had a int slice where -1 meant free and >0 meant file. This would have solved a lot of the issues I was running into with part 2
*/

func main() {
	input := util.ReadInput(util.Parameter())
	decompressedInput, fileNum := decompress(input)
	copy1 := make([]indexSpace, len(decompressedInput))
	copy2 := make([]indexSpace, len(decompressedInput))
	copy(copy1, decompressedInput)
	copy(copy2, decompressedInput)
	part1(copy1)
	part2(copy2, fileNum)
}

// part1 is the solution to part1, will swap right most filled spot with left most free
func part1(input []indexSpace) {

	left := findNextEmpty(input, 0)
	right := findNextFilled(input, len(input)-1)

	for left < right {
		input[left], input[right] = input[right], input[left]
		left = findNextEmpty(input, left+1)
		right = findNextFilled(input, right-1)
	}

	sum := 0

	for index, file := range input {
		if file.free {
			break
		}
		sum += index * file.fileNum
	}

	fmt.Println("(part1) Ans: ", sum)

}

// part2 solves part 2 by group files on the right and attempting to find suitable free spans starting from the left
func part2(input []indexSpace, fileNum int) {

	// Group indices by fileNum to avoid repeated iteration
	fileBlocks := make(map[int][]int)
	for i, space := range input {
		if space.fileNum > 0 {
			fileBlocks[space.fileNum] = append(fileBlocks[space.fileNum], i)
		}
	}

	// Iterate through file numbers in descending order
	for fileID := fileNum; fileID > 0; fileID-- {
		block, exists := fileBlocks[fileID]
		if !exists || len(block) == 0 {
			continue
		}

		fileStart, fileEnd := block[0], block[len(block)-1]
		fileLen := fileEnd - fileStart + 1

		// Find a suitable free block
		freeStart, currentFreeLen := -1, 0
		for i := 0; i < len(input); i++ {
			if input[i].free {
				currentFreeLen++
				if currentFreeLen == fileLen {
					freeStart = i - fileLen + 1
					break
				}
			} else {
				currentFreeLen = 0
			}
		}

		// Relocate the file if a free block is found
		if freeStart != -1 && freeStart < fileStart {
			for i := 0; i < fileLen; i++ {
				input[freeStart+i].free = false
				input[freeStart+i].fileNum = fileID
				input[fileStart+i].free = true
			}
		}
	}

	sum := 0
	for index, file := range input {
		if !file.free {
			sum += index * file.fileNum
		}
	}

	fmt.Println("(part2) Ans: ", sum)

}

// findNextFilled is a part 1 helper function, finds the next filled space from the given index (starting from the end of the slice)
func findNextFilled(input []indexSpace, start int) int {
	for index := start; index >= 0; index-- {
		if !input[index].free {
			return index
		}
	}
	return -1
}

// findNextEmpty is a part 1 helper function that finds the next empty space from the given index (starting from the beginning of the slice)
func findNextEmpty(input []indexSpace, start int) int {
	for index := start; index < len(input); index++ {
		if input[index].free {
			return index
		}
	}
	return -1
}

// decompress unpacks the initial disk map, it has a pattern of being (filled, unfilled) where the number represents the size
// Then each filled file is given an file ID
// So 123 => 0..111
func decompress(input []string) ([]indexSpace, int) {

	decompressed := make([]indexSpace, 0)
	split := strings.Split(input[0], "")
	fileIndex := 0
	used := 1
	for index := 0; index < len(split); index += 1 {

		space := util.ParseInt(split[index])

		if used%2 != 0 {
			toAdd := indexSpace{false, fileIndex}
			for i := 0; i < space; i++ {
				decompressed = append(decompressed, toAdd)
			}
			fileIndex++
		} else {
			toAdd := indexSpace{true, -1}
			for i := 0; i < space; i++ {
				decompressed = append(decompressed, toAdd)
			}
		}
		used += 1
	}
	return decompressed, fileIndex
}

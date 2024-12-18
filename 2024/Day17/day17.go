package main

import (
	computer "2024/Day17/threebitcomputer"
	util "2024/until"
	"fmt"
	"regexp"
	"slices"
)

const pattern = `-?\d+`

func main() {
	input := util.ReadInput(util.Parameter())
	a, b, c, program := parseInput(input)
	part1(a, b, c, program)
	part2(program)
}

func part1(a int, b int, c int, program []int) {
	comp := computer.NewComputer(a, b, c, false)
	comp.Run(program)
	comp.PrintOutput()

}

func part2(program []int) {
	results := make([]int, 0)
	values := make([]int, len(program))
	copy(values, program)
	findSolutions(0, program, values, &results, 1)
	slices.Sort(results)
	fmt.Println("Part 2: Ans", results[0])
}

func findSolutions(a int, program []int, values []int, results *[]int, level int) {
	if len(values) == 0 {
		return
	}

	val := values[len(values)-1]
	values = values[:len(values)-1]

	candidates := util.NewHashSet()

	for i := 0; i < 8; i++ {
		compt := computer.NewComputer(a+i, 0, 0, false)
		compt.OneTime = true
		compt.Run(program)
		result := compt.GetOutput()
		firstVal := util.ParseInt(result[0])
		if firstVal == val {
			candidates.Add(i)
			if level == len(program) {
				fmt.Println("valid a", a+i)
				*results = append(*results, a+i)
			}
		}
	}

	for _, can := range candidates.ToSlice() {
		candidate := can.(int)
		newValues := make([]int, len(values))
		copy(newValues, values)
		findSolutions((a+candidate)*8, program, newValues, results, level+1)
	}

}

func parseInput(input []string) (int, int, int, []int) {
	var a int
	var b int
	var c int
	program := make([]int, 0)
	reg := regexp.MustCompile(pattern)
	match1 := reg.FindAllStringSubmatch(input[0], -1)
	match2 := reg.FindAllStringSubmatch(input[1], -1)
	match3 := reg.FindAllStringSubmatch(input[2], -1)

	a = util.ParseInt(match1[0][0])
	b = util.ParseInt(match2[0][0])
	c = util.ParseInt(match3[0][0])

	for i := 4; i < len(input); i++ {
		match := reg.FindAllStringSubmatch(input[i], -1)
		for _, m := range match {
			for _, num := range m {
				program = append(program, util.ParseInt(num))
			}
		}
	}

	return a, b, c, program
}

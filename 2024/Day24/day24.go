package main

import (
	"2024/util"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	wirePattern = `(\w+)\s*:\s*(\d+)`
	gatePattern = `(\w+)\s+(XOR|OR|AND)\s+(\w+)\s+->\s+(\w+)`
	AND         = "AND"
	OR          = "OR"
	XOR         = "XOR"
)

type gate struct {
	inputOne  string
	inputTwo  string
	output    string
	operation string
}

var wires = map[string]int{}

func main() {
	gates := parseData(util.ReadInput("Day24/day24Data.txt"))
	part1(gates)
	part2(gates)
}

func part2(gates []gate) {
	swaps := checkParallelAdders(gates)
	slices.Sort(swaps)
	fmt.Println("(part 2) Ans: ", strings.Join(swaps, ","))

}

// checkParallelAdders identifies and swaps parallel adders in a list of logic gates.
// It returns a list of swapped wire names.
//
// Parameters:
//
//	gates []gate - A slice of gate structures representing the logic gates.
//
// Returns:
//
//	[]string - A slice of swapped wire names.
func checkParallelAdders(gates []gate) []string {
	var currentCarryWire string
	var swaps []string
	bit := 0

	for {
		// Generate wire names for the current bit position.
		xWire := fmt.Sprintf("x%02d", bit)
		yWire := fmt.Sprintf("y%02d", bit)
		zWire := fmt.Sprintf("z%02d", bit)

		if bit == 0 {
			// For the first bit, find the carry wire using the AND operation.
			currentCarryWire = findGate(xWire, yWire, AND, gates)
		} else {
			// For subsequent bits, find the XOR and AND gates for the current bit position.
			abXorGate := findGate(xWire, yWire, XOR, gates)
			abAndGate := findGate(xWire, yWire, AND, gates)

			// Find the XOR gate between the result of the previous XOR gate and the current carry wire.
			cinAbXorGate := findGate(abXorGate, currentCarryWire, XOR, gates)

			if cinAbXorGate == "" {
				// If the XOR gate is not found, swap the output wires of the XOR and AND gates,
				// reset the bit counter, and continue the loop.
				swaps = append(swaps, abXorGate, abAndGate)
				gates = swapOutputWires(abXorGate, abAndGate, gates)
				bit = 0
				continue
			}

			if cinAbXorGate != zWire {
				// If the XOR gate's output does not match the expected zWire, swap the output wires,
				// reset the bit counter, and continue the loop.
				swaps = append(swaps, cinAbXorGate, zWire)
				gates = swapOutputWires(cinAbXorGate, zWire, gates)
				bit = 0
				continue
			}

			// Find the AND gate between the previous XOR gate and the current carry wire,
			// and the OR gate for the carry wire.
			cinAbAndGate := findGate(abXorGate, currentCarryWire, AND, gates)
			carryWire := findGate(abAndGate, cinAbAndGate, OR, gates)
			currentCarryWire = carryWire
		}
		bit++
		if bit >= 45 {
			break
		}
	}
	return swaps
}

func swapOutputWires(wireA string, wireB string, gates []gate) []gate {
	newGates := make([]gate, 0)
	for _, g := range gates {
		if g.output == wireA {
			newGates = append(newGates, gate{inputOne: g.inputOne, inputTwo: g.inputTwo, output: wireB, operation: g.operation})
		} else if g.output == wireB {
			newGates = append(newGates, gate{inputOne: g.inputOne, inputTwo: g.inputTwo, output: wireA, operation: g.operation})
		} else {
			newGates = append(newGates, g)
		}
	}
	return newGates
}

func findGate(xWire string, yWire string, gate string, gates []gate) string {
	for _, g := range gates {
		if g.inputOne == xWire && g.inputTwo == yWire && g.operation == gate || g.inputOne == yWire && g.inputTwo == xWire && g.operation == gate {
			return g.output
		}
	}
	return ""
}

func part1(gates []gate) {

	knownInputs := make([]gate, 0)
	unknownInputs := make([]gate, 0)

	for _, gate := range gates {
		if wires[gate.inputOne] != -1 && wires[gate.inputTwo] != -1 {
			knownInputs = append(knownInputs, gate)
		} else {
			unknownInputs = append(unknownInputs, gate)
		}
	}

	processKnownGates(knownInputs)
	processGates(unknownInputs)

	zPosition := make([]string, 0)
	for k, _ := range wires {
		if k[0] == 'z' {
			zPosition = append(zPosition, k)
		}
	}

	slices.Sort(zPosition)
	binary := ""
	for z := len(zPosition) - 1; z >= 0; z-- {
		if wires[zPosition[z]] == 1 {
			binary += "1"
		} else {
			binary += "0"
		}
	}

	decimal, _ := strconv.ParseInt(binary, 2, 64)
	fmt.Println("(part 1) Ans: ", decimal, binary)

}

func processGates(inputs []gate) {
	reattempt := make([]gate, 0)

	for len(inputs) > 0 {
		for _, g := range inputs {
			if wires[g.inputOne] != -1 && wires[g.inputTwo] != -1 {
				processGate(g)
			} else {
				reattempt = append(reattempt, g)
			}
		}

		inputs = reattempt
		reattempt = make([]gate, 0)
	}
}

func processKnownGates(inputs []gate) {
	for _, g := range inputs {
		processGate(g)
	}
}

func processGate(g gate) {
	switch g.operation {
	case "AND":
		wires[g.output] = wires[g.inputOne] & wires[g.inputTwo]
	case "OR":
		wires[g.output] = wires[g.inputOne] | wires[g.inputTwo]
	case "XOR":
		wires[g.output] = wires[g.inputOne] ^ wires[g.inputTwo]
	}
}

func parseData(input []string) []gate {
	toReturn := make([]gate, 0)

	wireReg := regexp.MustCompile(wirePattern)
	gateReg := regexp.MustCompile(gatePattern)

	splitIndex := slices.Index(input, "")
	initialWires, gates := input[:splitIndex], input[splitIndex+1:]

	for _, wire := range initialWires {
		wireGroups := wireReg.FindStringSubmatch(wire)
		wireName := wireGroups[1]
		wireVal := util.ParseInt(wireGroups[2])
		wires[wireName] = wireVal
	}

	for _, rawGate := range gates {
		gateGroups := gateReg.FindStringSubmatch(rawGate)
		inputOne := gateGroups[1]
		inputTwo := gateGroups[3]
		logicGate := gateGroups[2]
		output := gateGroups[4]

		if _, ok := wires[inputOne]; !ok {
			wires[inputOne] = -1
		}

		if _, ok := wires[inputTwo]; !ok {
			wires[inputTwo] = -1
		}

		if _, ok := wires[output]; !ok {
			wires[output] = -1
		}

		toReturn = append(toReturn, gate{inputOne: inputOne, inputTwo: inputTwo, output: output, operation: logicGate})
	}

	return toReturn
}

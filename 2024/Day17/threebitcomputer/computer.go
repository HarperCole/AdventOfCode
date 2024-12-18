package threebitcomputer

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Computer struct {
	a          int      // A register
	b          int      // B register
	c          int      // C register
	ip         int      // instruction pointer
	jump       bool     // signals if we have jumped instructions
	outResults []string // results of out command
	dump       bool
	done       bool
	OneTime    bool
}

const (
	eight = 0x8
	move  = 2
)

// NewComputer creates a new Computer and sets the registers to the passed in values
func NewComputer(a, b, c int, dump bool) *Computer {
	cptr := new(Computer)
	cptr.a = a
	cptr.b = b
	cptr.c = c
	cptr.ip = 0
	cptr.dump = dump

	return cptr
}

// Run runs the program and prints out logs if dump is set
func (comp *Computer) Run(program []int) {
	ip := 0

	for ip < len(program) {
		opcode := program[ip]
		operand := program[ip+1]
		result := comp.execute(opcode, operand)
		if result == -1 {
			fmt.Println("ERROR SOMETHING BAD HAPPENED!")
			comp.debug(opcode, operand)
			os.Exit(1)
		}
		if comp.dump {
			comp.debug(opcode, operand)
			fmt.Println()
		}
		ip = result
		if comp.done {
			break
		}
	}
}

// getComboOperand retrieves the combo operand of the passed in operand
func (comp *Computer) getComboOperand(operand int) int {
	switch operand {
	case 4:
		return comp.a
	case 5:
		return comp.b
	case 6:
		return comp.c
	default:
		if operand <= 3 {
			return operand
		} else {
			return 7
		}
	}
}

// division calculates the division (x / 2^y) where x is the numerator and y is the denominator
func (comp *Computer) division(numerator, denominator int) float64 {
	return float64(numerator) / math.Pow(2.0, float64(denominator))
}

// execute executes the opcode with the operand passed in
func (comp *Computer) execute(opcode, operand int) int {
	// this is where i'll call the specific instructions

	switch opcode {
	case 0:
		return comp.adv(operand)
	case 1:
		return comp.bxl(operand)
	case 2:
		return comp.bst(operand)
	case 3:
		return comp.jnz(operand)
	case 4:
		return comp.bxc(operand)
	case 5:
		return comp.out(operand)
	case 6:
		return comp.bdv(operand)
	case 7:
		return comp.cdv(operand)
	default:
		return -1
	}
}

// moveInstructionPointer is will move the ip to the next instruction
// if the jump flag has been set, then increment ip by 1, otherwise increase by move
func (comp *Computer) moveInstructionPointer() int {

	comp.ip += move

	if comp.dump {
		fmt.Println(fmt.Sprintf("Moving IP %d, IP %d", move, comp.ip))
	}

	return comp.ip
}

// adv instruction performs division using the value stored in register a and the combo operand
// opcode 0
func (comp *Computer) adv(operand int) int {
	result := comp.division(comp.a, comp.getComboOperand(operand))
	comp.a = int(math.Trunc(result))
	if comp.dump {
		fmt.Println(fmt.Sprintf("adv - optcode 0 with operand %d, divsion result %f, register A: %d", operand, result, comp.a))
	}
	return comp.moveInstructionPointer()
}

// bxl takes the XOR of register b with the literal operand
// opcode 1
func (comp *Computer) bxl(operand int) int {
	result := comp.b ^ operand
	comp.b = result

	if comp.dump {
		fmt.Println(fmt.Sprintf("bxl optcode 1 - operand %d, result %d, register B: %d", operand, result, comp.b))
	}

	return comp.moveInstructionPointer()
}

// bst calculates the combo operand and then modulo with 8 and writes the value to b register
// opcode 2
func (comp *Computer) bst(operand int) int {
	result := comp.getComboOperand(operand) % eight
	comp.b = result
	if comp.dump {
		fmt.Println(fmt.Sprintf("bst optcode 2 - operand %d, result %d, register B: %d", operand, result, comp.b))
	}
	return comp.moveInstructionPointer()
}

// jnz moves the instruction pointer to the literal operand value, sets the jump flag
// opcode 3
func (comp *Computer) jnz(operand int) int {
	if comp.OneTime || comp.a == 0 {
		comp.done = true
	} else {
		comp.ip = operand
	}

	return comp.ip

}

// bxc finds the XOR of register b and c then stores it in register b, operand is ignored
// opcode 4
func (comp *Computer) bxc(operand int) int {
	result := comp.b ^ comp.c
	comp.b = result
	if comp.dump {
		fmt.Println(fmt.Sprintf("bxc optcode 4 - register B: %d, register C: %d, result %d", comp.b, comp.c, result))
	}
	return comp.moveInstructionPointer()
}

// out calculates the value of combo operand modulo eight, and then outputs the value
// opcode 5
func (comp *Computer) out(operand int) int {
	comp.outResults = append(comp.outResults, strconv.Itoa(comp.getComboOperand(operand)%eight))
	if comp.dump {
		fmt.Println(fmt.Sprintf("out optcode 5 - output %v", comp.outResults))
	}
	return comp.moveInstructionPointer()
}

// bdv works the same as adqv except the results are stored into the b register
// opcode 6
func (comp *Computer) bdv(operand int) int {
	results := comp.division(comp.a, comp.getComboOperand(operand))
	comp.b = int(math.Trunc(results))
	if comp.dump {
		fmt.Println(fmt.Sprintf("bdv - optcode 6 with operand %d, divsion result %f, register B: %d", operand, results, comp.b))
	}
	return comp.moveInstructionPointer()
}

// cdv works the same as adv except the results are stored in the c register
// opcode 7
func (comp *Computer) cdv(operand int) int {
	results := comp.division(comp.a, comp.getComboOperand(operand))
	comp.c = int(math.Trunc(results))

	if comp.dump {
		fmt.Println(fmt.Sprintf("cdv - optcode 7 with operand %d, divsion result %f, register A: %d", operand, results, comp.c))
	}
	return comp.moveInstructionPointer()
}

// PrintOutput prints out the results of the out command thus far
func (comp *Computer) PrintOutput() {
	fmt.Println(fmt.Sprintf("output results %s\n", strings.Join(comp.outResults, ",")))
}

// DumpRegisters is a debugger helper function for dumping the current state of the registers
func (comp *Computer) DumpRegisters() {
	fmt.Printf("Register A: %d, B: %d, C: %d\n", comp.a, comp.b, comp.c)
}

func (comp *Computer) debug(opcode, operand int) {
	fmt.Printf("Commands just ran: opcode %d, operand %d\n", opcode, operand)
	comp.DumpRegisters()
}

func (comp *Computer) GetOutput() []string {
	return comp.outResults
}

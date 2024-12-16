package main

import (
	util "2024/until"
	"fmt"
	"regexp"
	"strings"
)

const (
	pattern = `-?\d+`
	seconds = 100
)

/*
	Advent of Code Day 14:
		Part 1: Just run the simulation of the robots moving for 100 seconds, then check the quadrants
		Part 2: I literally printed out the grid 10000 times and looked for the first instance where there was enough # (character used to indicate that a robot where there) in a row to see the Christmas tree
*/

// robot is a struct representation
//   - position holds the x,y coordinates of the robot
//   - velocity holds the dx,dy of the robot
type robot struct {
	position [2]int
	velocity [2]int
}

func main() {
	robots := grabRobots(util.ReadInput(util.Parameter()))
	wide, tall := 101, 103
	//part1(robots, wide, tall, seconds)
	part2(robots, wide, tall, 100000)
}

func part2(robots []robot, wide int, tall int, seconds int) {
	for second := 0; second < seconds; second++ {
		for index, rob := range robots {
			robots[index] = moveRobot(rob, wide, tall)
		}

		fmt.Println()
		fmt.Println("Second: ", second)
		printGrid(robots, wide, tall)
		fmt.Println()
	}
}

// part1 runs the simulation for a specified number of seconds, moves robots
// based on their velocities, and calculates a safety factor by counting
// robots in each quadrant. The results are printed to the console.
func part1(robots []robot, wide, tall, seconds int) {
	for second := 0; second < seconds; second++ {
		for index, rob := range robots {
			robots[index] = moveRobot(rob, wide, tall)
		}
		printGrid(robots, wide, tall)
		fmt.Println()
	}
	q1, q2, q3, q4 := findQuadrant(robots, wide, tall)

	//printGrid(robots, wide, tall)

	fmt.Println("(part 1) Ans: ", q1*q2*q3*q4)
}

// printGrid renders the current positions of the robots on a grid and
// displays it in the console. Empty cells are marked with '.', and cells
// with robots are marked with '#'.
func printGrid(robots []robot, wide int, tall int) {
	grid := make([][]string, tall)

	for i := 0; i < tall; i++ {
		grid[i] = make([]string, wide)
		for j := 0; j < wide; j++ {
			grid[i][j] = "."
		}
	}

	for _, rob := range robots {
		x, y := rob.position[0], rob.position[1]
		if grid[y][x] == "." {
			grid[y][x] = "#"
		}
	}

	for _, row := range grid {
		fmt.Println(strings.Join(row, ""))
	}

}

// findQuadrant calculates the number of robots in each of the four quadrants
// of the grid. Quadrants are determined based on the midpoint of the grid.
func findQuadrant(robots []robot, wide int, tall int) (int, int, int, int) {
	hMiddle := wide / 2
	vMiddle := tall / 2
	q1, q2, q3, q4 := 0, 0, 0, 0

	for _, rob := range robots {
		x, y := rob.position[0], rob.position[1]

		if x >= 0 && x < hMiddle && y >= 0 && y < vMiddle {
			q1++
		} else if x > hMiddle && x <= wide-1 && y >= 0 && y < vMiddle {
			q2++
		} else if x >= 0 && x < hMiddle && y > vMiddle && y <= tall-1 {
			q3++
		} else if x > hMiddle && x <= wide-1 && y > vMiddle && y <= tall-1 {
			q4++
		}
	}

	return q1, q2, q3, q4

}

// moveRobot updates a robot's position based on its velocity and applies
// wrapping to ensure the robot stays within the grid boundaries.
func moveRobot(rob robot, wide int, tall int) robot {
	x := (rob.position[0] + rob.velocity[0] + wide) % wide
	y := (rob.position[1] + rob.velocity[1] + tall) % tall

	rob.position[0] = x
	rob.position[1] = y
	return rob
}

// grabRobots parses a list of input strings to create a slice of robots. Each
// string should contain position and velocity values in the format
// "p=x,y v=dx,dy"
func grabRobots(input []string) []robot {
	toReturn := make([]robot, 0)

	regex := regexp.MustCompile(pattern)

	for _, line := range input {
		matches := regex.FindAllStringSubmatch(line, -1)
		x, y := util.ParseInt(matches[0][0]), util.ParseInt(matches[1][0])
		dx, dy := util.ParseInt(matches[2][0]), util.ParseInt(matches[3][0])
		toReturn = append(toReturn, robot{[2]int{x, y}, [2]int{dx, dy}})
	}
	return toReturn
}

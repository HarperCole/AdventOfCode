package until

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadInput takes a filename and reads in the contents into a string slice for each line of the file
func ReadInput(filename string) []string {
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

// Parameter generic CLI arguments for a file name
func Parameter() string {
	filename := flag.String("file", "", "File to read")
	flag.Parse()

	if *filename == "" {
		fmt.Println("Please provide a file!")
	}
	return *filename
}

// TransformStringSliceInto2DMatrix transforms a string slice into a 2D matrix
//
//	[ABC ABC ABC] -> [[A B C] [A B C] [A B C]]
func TransformStringSliceInto2DMatrix(input []string) [][]string {
	toReturn := make([][]string, 0)

	for _, line := range input {
		toReturn = append(toReturn, strings.Split(line, ""))
	}
	return toReturn
}

// ParseInt is a helper function to convert a string into an int
func ParseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error trying to parse int from string: ", err)
		os.Exit(1)
	}
	return val
}

func StringToIntSlice(input []string) []int {
	toReturn := make([]int, 0)

	for _, line := range input {
		toReturn = append(toReturn, ParseInt(line))
	}
	return toReturn
}

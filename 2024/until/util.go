package until

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

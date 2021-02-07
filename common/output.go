package common

import (
	"fmt"
	"golang.org/x/term"
	"math"
	"os"
	"strings"
	"text/tabwriter"
)

// printFileDefault print files line by line
func printFileDefault(files []string) {
	for _, str := range files {
		fmt.Println(str)
	}
}

// validateRows validates rows by getting the sum of all the longest strings in
// each column and comparing the value with the space left
func validateRows(table [][]string, leftspace int) bool {
	m := make(map[int]int)
	for _, row := range table {
		for x := range row {
			length := len(row[x])
			if length > m[x] {
				m[x] = length
			}
		}
	}

	width := 0
	for _, val := range m {
		width += val
	}

	if width > leftspace {
		return false
	}

	return true
}

// PrintFiles print files in table format when applicable
// TODO: optimize algorithm
func PrintFiles(files []string) {
	// get terminal size of standard output
	width, _, err := term.GetSize(int(os.Stdout.Fd()))

	// if error, just print files line by line
	if err != nil {
		printFileDefault(files)
		return
	}

	// reverse iterate
	// the iterator i represents the number of columns to test
	// stop when the iterator reaches 1
	for i := len(files); i > 1; i-- {
		// compute whitespace required beforehand
		// i*2  -> padded sections * 2 spaces
		whitespace := i * 2
		leftspace := width - whitespace

		// a file must have at least one character
		// so if leftspace is less than number of columns
		// there is no way for files to be printed in table format
		if leftspace < i {
			continue
		}

		// calculate number of rows needed to fit the files
		// take the ceiling of quotient of len(files) / # of columns
		rows := int(math.Ceil(float64(len(files)) / float64(i)))

		// it is possible to use offsets, and some simple math to calculate
		// file ordering, but it is error prone and kind of a hassle
		// just create 2D slice instead for simplicity
		table := make([][]string, rows)
		for j := range table {
			table[j] = make([]string, 0)
		}

		// fill 2D slice
		// the order is defined as
		//
		// a  d  g
		// b  e  h
		// c  f  i
		for j := range files {
			y := j % rows
			table[y] = append(table[y], files[j])
		}

		// validate table so that it can fit in the calculated leftspace
		if !validateRows(table, leftspace) {
			continue
		}

		// use tabwriter making use of the Elastic Tabstops algorithm
		// the columns are padded with at least 2 spaces
		writer := tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', 0)
		for _, row := range table {
			fmt.Fprintln(writer, strings.Join(row, "\t")+"\t")
		}
		writer.Flush()

		return
	}

	// use default print if columned printing is not possible
	printFileDefault(files)
}

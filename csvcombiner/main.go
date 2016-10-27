package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/YoshikiShibata/tools/util/files"
)

// csvcombiner combines all rows in all csv file.
// The number of columns must be same in all csv file.
// Only the first column is considered to be same in all csv file.

type row []string // each row values
type csvContents struct {
	name  string // filename without ".csv"
	lines []row
}

func main() {
	csvFiles, err := files.ListFiles(".",
		func(f string) bool {
			return strings.HasSuffix(f, ".csv")
		})

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	var csvContentsList []*csvContents

	for _, csvFile := range csvFiles {
		csv, err := toCVSContents(csvFile)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		csvContentsList = append(csvContentsList, csv)
	}

	printHeader(csvContentsList)
	printEachLine(csvContentsList)

}

func toCVSContents(f string) (*csvContents, error) {
	lines, err := files.ReadAllLines(f)
	if err != nil {
		return nil, err
	}

	var csv = csvContents{f, nil}

	for _, line := range lines {
		row := strings.Split(line, ",")
		csv.lines = append(csv.lines, row)
	}
	return &csv, nil
}

func printHeader(csvContentsList []*csvContents) {
	for _, csvContents := range csvContentsList {
		name := csvContents.name
		fmt.Printf(",%s,T", name[:len(name)-len(".csv")])
		for i := 3; i < len(csvContents.lines[0]); i++ {
			fmt.Printf(",")
		}
	}
	fmt.Println()
}

func printEachLine(csvContentsList []*csvContents) {
	noOfRows := len(csvContentsList[0].lines)
	for row := 0; row < noOfRows; row++ {
		fmt.Printf("%s", csvContentsList[0].lines[row][0])

		for _, csvContents := range csvContentsList {
			for i, column := range csvContents.lines[row] {
				if i != 0 {
					fmt.Printf(",%s", column)
				}
			}
		}
		fmt.Println()
	}
}

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/YoshikiShibata/tools/util/files"
)

// cvscombiner combines all rows in all cvs file.
// The number of columns must be same in all cvs file.
// Only the first column is considered to be same in all cvs file.

type row []string // each row values
type cvsContents struct {
	name  string // filename without ".cvs"
	lines []row
}

func main() {
	cvsFiles, err := files.ListFiles(".",
		func(f string) bool {
			return strings.HasSuffix(f, ".cvs")
		})

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	var cvsContentsList []*cvsContents

	for _, cvsFile := range cvsFiles {
		cvs, err := toCVSContents(cvsFile)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		cvsContentsList = append(cvsContentsList, cvs)
	}
}

func toCVSContents(f string) (*cvsContents, error) {
	lines, err := files.ReadAllLines(f)
	if err != nil {
		return nil, err
	}

	var cvs = cvsContents{f, nil}

	for _, line := range lines {
		row := strings.Split(line, ",")
		cvs.lines = append(cvs.lines, row)
	}
	return &cvs, nil
}

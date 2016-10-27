package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/YoshikiShibata/tools/util/file"
)

// cvscombiner combines all rows in all cvs file.
// The number of columns must be same in all cvs file.
// Only the first column is considered to be same in all cvs file.

type row []string // each row values
type cvsFile struct {
	name  string // filename without ".cvs"
	lines []row
}

func main() {
	files, err := file.ListFiles(".",
		func(f string) bool {
			return strings.HasSuffix(f, ".cvs")
		})

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func toCVS(file string) (*cvsFile, error) {
	lines, err := file.ReadAllLines(file)
	if err != nil {
		return nil, err
	}

	var cvs = cvsFile{file, nil}

	for _, line := range lines {
		row := strings.Split(line, ",")
		cvs.lines = append(cvs.lines, row)
	}
	return &cvs, nil
}

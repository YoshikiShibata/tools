// Copyright © 2017 Yoshiki Shibata. All rights reserved.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/YoshikiShibata/tools/util/files"
)

// gpldriver just scan subdirectories, check whether there are
// runall.bash, testall.bash, runall.bat, testall.bat files, and then
// reports the result.
func main() {
	if len(os.Args) != 3 {
		showUsage()
	}
	cwd, err := os.Getwd()
	if err != nil {
		exit(err, 1)
	}
	fmt.Fprintf(os.Stderr, "CWD = %s\n", cwd)

	directies := readDirecties(os.Args[1])
	var buf bytes.Buffer
	for _, d := range directies {

		fmt.Fprintf(&buf, "%s,", d)
		fmt.Fprintf(os.Stderr, "%s(%s):\n", d, os.Args[2])

		dir := cwd + "/" + d

		scanFiles(&buf, dir, []string{"runall.bash", "runall.bat"})
		fmt.Fprintf(&buf, ",")
		scanFiles(&buf, dir, []string{"testall.bash", "testall.bat"})

		fmt.Fprintln(&buf)
	}

	fmt.Print(buf.String())
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "usage: gpldriver [directive csv file] [name]")
	os.Exit(1)
}

func exit(err error, code int) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(code)
}

func readDirecties(file string) []string {
	lines, err := files.ReadAllLines(file)
	if err != nil {
		exit(err, 1)
	}

	directives := make([]string, 0, len(lines))
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		directives = append(directives, line)
	}
	return directives
}

func scanFiles(buf io.Writer, directory string, targets []string) {
	files, err := files.ListFiles(directory,
		func(filename string) bool {
			name := strings.ToLower(filename)
			for _, t := range targets {
				if name == t {
					return true
				}
			}
			return false
		})

	if err != nil {
		fmt.Fprintf(buf, "N/A")
		return
	}

	if len(files) == 0 {
		fmt.Fprintf(buf, "0")
	} else {
		fmt.Fprintf(buf, "1")
	}
}

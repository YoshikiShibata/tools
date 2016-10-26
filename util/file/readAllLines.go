// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ReadAllLines reads all lines from a file.
// file must be encoded in UTF-8.
func ReadAllLines(filePath string) ([]string, error) {
	if f, err := os.Open(filePath); err != nil {
		return nil, err
	} else {
		defer f.Close()

		return readLines(f)
	}
}

func readLines(reader io.Reader) ([]string, error) {
	lines := []string{}
	r := bufio.NewReader(reader)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return lines, nil
			}
			fmt.Printf("%v\n", err)
			return lines, err
		}
		lines = append(lines, line)
	}
}

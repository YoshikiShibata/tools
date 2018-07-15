// Copyright © 2016, 2018 Yoshiki Shibata. All rights reserved.

package files

import (
	"bufio"
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
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import "fmt"

// directive defines the information in the directive csv file: each information is
// seperated by comma. The first field is the directory under "src" and "test" directories.
// runOptions are passed to "oak run" command as additional arguments
type directive struct {
	directory  string
	runOptions []string
}

// Exit code of oak.
// [TODO] this should be refactored not to be copied from main.go of the oak.
const (
	codeError            = 1 // general error
	codeCompileError     = 2 // compile error
	codeExecutionTimeout = 3 // execution timeout
	codeTestsFailed      = 4 // test failed
	codeNoMainMethod     = 5 // no main method
	codeMainFailed       = 6 // executing main failed
)

var runCodeMap = map[int]string{
	0:                    "2",
	codeError:            "err",
	codeCompileError:     "CE",
	codeExecutionTimeout: "fail",
	codeTestsFailed:      "fail",
	codeNoMainMethod:     "1",
	codeMainFailed:       "fail",
}

var testCodeMap = map[int]string{
	0:                    "2",
	codeError:            "err",
	codeCompileError:     "CE",
	codeExecutionTimeout: "fail",
	codeTestsFailed:      "fail",
	codeNoMainMethod:     "fail",
	codeMainFailed:       "fail",
}

// oakdriver runs oak command for Java / Java 8 programing courses.
// The output of the oakdriver is the os.Stdout.
func main() {
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "usage: oakdriver [directory] [directive csv file]")
	os.Exit(1)
}

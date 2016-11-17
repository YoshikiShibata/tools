// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/YoshikiShibata/tools/util/files"
)

// directive defines the information in the directive csv file: each information is
// seperated by comma. The first field is the directory under "src" and "test" directories.
// runOptions are passed to "oak run" command as additional arguments
// Any line starting with # will be treated as comment and ignored
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

var runCodeOutput = map[int]string{
	0:                    "2",
	codeError:            "Error",
	codeCompileError:     "CompileError",
	codeExecutionTimeout: "Timeout",
	codeTestsFailed:      "Failed",
	codeNoMainMethod:     "1",
	codeMainFailed:       "Failed",
}

var testCodeOutput = map[int]string{
	0:                    "2",
	codeError:            "Error",
	codeCompileError:     "CompileError",
	codeExecutionTimeout: "Timeout",
	codeTestsFailed:      "Failed",
	codeNoMainMethod:     "fail",
	codeMainFailed:       "fail",
}

// oakdriver runs oak command for Java / Java 8 programing courses.
// The output of the oakdriver is the os.Stdout.
func main() {
	if len(os.Args) != 2 {
		showUsage()
	}
	cwd, err := os.Getwd()
	fmt.Fprintf(os.Stderr, "CWD = %s\n", cwd)
	if err != nil {
		exit(err, 1)
	}

	for _, d := range readDirectives(os.Args[1]) {
		fmt.Printf("%s,", d.directory)
		fmt.Fprintf(os.Stderr, "%s:\n", d.directory)

		run(cwd, &d)
		fmt.Printf(",")
		test(cwd, &d)

		fmt.Println()
	}
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "usage: oakdriver [directive csv file]")
	os.Exit(1)
}

func exit(err error, code int) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(code)
}

func readDirectives(file string) []directive {
	lines, err := files.ReadAllLines(file)
	if err != nil {
		exit(err, 1)
	}
	var directives []directive
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		columns := strings.Split(line, ",")
		if len(columns) == 1 {
			d := directive{columns[0], nil}
			directives = append(directives, d)
			continue
		}

		if len(columns) >= 2 {
			d := directive{columns[0], columns[1:]}
			directives = append(directives, d)
			continue
		}

		exit(fmt.Errorf("Columns is zero"), 1)
	}

	return directives
}

func run(cwd string, d *directive) {
	if err := os.Chdir(cwd + "/src/" + d.directory); err != nil {
		fmt.Printf("N/A")
		return
	}

	args := []string{"run"}
	args = append(args, d.runOptions...)

	cmd := exec.Command("oak", args...)
	redirect(cmd)
	err := cmd.Run()
	if err == nil {
		fmt.Printf("%s", runCodeOutput[0])
	} else {
		exitCode := extractExitCode(err)
		fmt.Printf("%s", runCodeOutput[exitCode])
	}
}

func test(cwd string, d *directive) {
	if err := os.Chdir(cwd + "/test/" + d.directory); err != nil {
		fmt.Printf("N/A")
		return
	}

	cmd := exec.Command("oak", "test")
	redirect(cmd)
	err := cmd.Run()
	if err == nil {
		fmt.Printf("%s", testCodeOutput[0])
	} else {
		exitCode := extractExitCode(err)
		fmt.Printf("%s", testCodeOutput[exitCode])
	}
}

func redirect(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		exit(err, codeError)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		exit(err, codeError)
	}
	go func() {
		io.Copy(os.Stderr, stderr)
	}()
	go func() {
		io.Copy(os.Stderr, stdout)
	}()
}

func extractExitCode(cmdErr error) int {
	exitError, ok := cmdErr.(*exec.ExitError)
	if !ok {
		exit(cmdErr, 1)
	}

	if s, ok := exitError.Sys().(syscall.WaitStatus); ok {
		return s.ExitStatus()
	}
	panic(exitError)
}

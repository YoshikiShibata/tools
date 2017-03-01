// Copyright © 2016, 2017 Yoshiki Shibata. All rights reserved.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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

const concurrent = 3

// oakdriver runs oak command for Java / Java 8 programing courses.
// The output of the oakdriver is the os.Stdout.
func main() {
	if len(os.Args) != 3 {
		showUsage()
	}
	cwd, err := os.Getwd()
	fmt.Fprintf(os.Stderr, "CWD = %s\n", cwd)
	if err != nil {
		exit(err, 1)
	}

	directives := readDirectives(os.Args[1])
	results := make([]*bytes.Buffer, 0, len(directives))

	sem := make(chan struct{}, concurrent)

	var wg sync.WaitGroup

	for i, d := range directives {
		var buf bytes.Buffer
		results = append(results, &buf)

		wg.Add(1)
		go func(buf io.Writer, d directive, index int) {
			sem <- struct{}{}
			fmt.Fprintf(buf, "%s,", d.directory)
			fmt.Fprintf(os.Stderr, "%s(%s):\n", d.directory, os.Args[2])

			run(buf, cwd, &d, index)
			fmt.Fprintf(buf, ",")
			test(buf, cwd, &d, index)

			fmt.Fprintln(buf)
			<-sem
			wg.Done()
		}(&buf, d, i)
	}

	wg.Wait()

	for _, buf := range results {
		fmt.Print(buf.String())
	}
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "usage: oakdriver [directive csv file] [name]")
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

func tempOption(index int) string {
	return "-temp=/tmp/oak" + strconv.Itoa(index%concurrent)
}

func run(buf io.Writer, cwd string, d *directive, index int) {
	if err := os.Chdir(cwd + "/src/" + d.directory); err != nil {
		fmt.Fprintf(buf, "N/A")
		return
	}

	args := []string{"-l", tempOption(index), "run"}
	args = append(args, d.runOptions...)

	cmd := exec.Command("oak", args...)
	redirect(cmd)
	err := cmd.Run()
	if err == nil {
		fmt.Fprintf(buf, "%s", runCodeOutput[0])
	} else {
		exitCode := extractExitCode(err)
		fmt.Fprintf(buf, "%s", runCodeOutput[exitCode])
	}
}

func test(buf io.Writer, cwd string, d *directive, index int) {
	if err := os.Chdir(cwd + "/test/" + d.directory); err != nil {
		fmt.Fprintf(buf, "N/A")
		return
	}

	cmd := exec.Command("oak", "-l", tempOption(index), "test")
	redirect(cmd)
	err := cmd.Run()
	if err == nil {
		fmt.Fprintf(buf, "%s", testCodeOutput[0])
	} else {
		exitCode := extractExitCode(err)
		fmt.Fprintf(buf, "%s", testCodeOutput[exitCode])
	}
}

func redirect(cmd *exec.Cmd) {
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
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

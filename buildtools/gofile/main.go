package main

import (
	"flag"
	"os"
	"path/filepath"
	"github.com/spekary/goradd/util"
	"strings"
	"os/exec"
	"bufio"
)

var remove = flag.Bool("r", false, "Remove indicated files and directories.")
var generate = flag.Bool("g", false, "Run go generate on the given files.")

func main() {
	var files []string

	flag.Parse()

	// Replace all GOPATH strings with the current go path, expand paths, apply OS specific stuff, clean path
	// This stuff allows the tool to use forward slashes when specifying paths on WINDOWS, making it universal
	for _,f := range flag.Args() {
		f = filepath.FromSlash(f)
		f = strings.Replace(f, "GOPATH", util.GoPath(), 1)
		if files2,_ := filepath.Glob(f); files2 != nil {
			for _, file2 := range files2 {
				f2, err := filepath.Abs(file2)	// not sure this is necessary, but just in case
				if err == nil {
					files = append(files, f2)
				}
			}
		}
	}


	if *remove {
		for _,f := range files {
			os.RemoveAll(f)
		}
	}

	if *generate {
		for _,f := range files {
			executeCmd("go", "generate", f)
		}
	}
}

func executeCmd(command string, args ...string) (stdOutText string, stdErrText string, err2 error) {
	cmd := exec.Command(command, args...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		err2 = err
		return
	}

	defer stdOut.Close()

	scanner := bufio.NewScanner(stdOut)
	go func() {
		for scanner.Scan() {
			stdOutText += scanner.Text() + "\n"
		}
	}()

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		err2 = err
		return
	}

	defer stdErr.Close()

	stdErrScanner := bufio.NewScanner(stdErr)
	go func() {
		for stdErrScanner.Scan() {

			stdErrText += stdErrScanner.Text() + "\n"
		}
	}()

	err = cmd.Start()
	if err != nil {
		err2 = err
		return
	}

	err = cmd.Wait()

	if err != nil {
		err2 = err
	}
	return
}



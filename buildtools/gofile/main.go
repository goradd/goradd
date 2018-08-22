package main

import (
	"os"
	"path/filepath"
	"github.com/spekary/goradd/util"
	"strings"
	"os/exec"
	"bufio"
)

func main() {
	var files []string

	cmd := os.Args[1]

	// Replace all GOPATH strings with the current go path, expand paths, apply OS specific stuff, clean path
	// This stuff allows the tool to use forward slashes when specifying paths on WINDOWS, making it universal
	for _,f := range os.Args[2:] {
		f = filepath.FromSlash(f)
		f = strings.Replace(f, "GOPATH", util.GoPath(), 1)

		if cmd == "mkdir" {
			files = append(files, f)
		} else {
			if files2,_ := filepath.Glob(f); files2 != nil {
				for _, file2 := range files2 {
					_,fName := filepath.Split(file2)
					if file2 != "" && fName[0] != '.' {	// ignore dot files
						f2, err := filepath.Abs(file2)	// not sure this is necessary, but just in case
						if err == nil {
							files = append(files, f2)
						}
					}
				}
			}
		}
	}

	switch cmd {
	case "remove":
		for _,f := range files {
			os.RemoveAll(f)
		}
	case "generate":
		for _,f := range files {
			executeCmd("go", "generate", f)
		}
	case "copy":
		copyFiles(files)

	case "mkdir":
		makeDirectory(files)
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

func copyFiles(files []string) {
	if len(files) < 2 {
		panic("the copy command must have a source and destination")
	}

	dest := files[len(files) - 1]
	src := files[:len(files) - 1]

	var fInfo1, fInfo2 os.FileInfo
	var err, err2 error

	fInfo2,err2 = os.Stat(dest)

	if !fInfo2.IsDir() && len(files) != 2 {
		panic("cannot copy more than one item to a file destination")
	}

	for _,file := range src {

		if fInfo1, err = os.Stat(file); err != nil {
			panic (err)
		}

		if fInfo1.IsDir() {
			if err2 != nil {
				panic (err2)
			}
			if !fInfo2.IsDir() {
				panic ("Cannot copy a directory onto a file")
			} else {
				err = util.DirectoryCopy(file, dest)
				if err != nil {
					panic (err)
				}
			}
		} else {
			if err != nil && !os.IsNotExist(err) {
				panic(err)
			}
			destFile := dest
			if fInfo2.IsDir() {
				_,f := filepath.Split(file)
				destFile = filepath.Join(destFile, f)
			}
			err = util.FileCopy(file, destFile)
			if err != nil {
				panic (err)
			}
		}
	}
}

func makeDirectory(files []string) {
	for _,dir := range files {
		os.MkdirAll(dir, os.FileMode(0777))
	}
}
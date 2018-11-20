package main

import (
	"github.com/spekary/goradd/pkg/sys"
	"os"
	"path/filepath"
	"strings"
)


func main() {
	var files []string
	var curOption string
	var excludes = make(map[string]bool)

	modules, err  := sys.ModulePaths()

	if err != nil {
		panic(err)
	}

	cmd := os.Args[1]

	// Replace all GOPATH strings with the current go path, expand paths, apply OS specific stuff, clean path
	// This stuff allows the tool to use forward slashes when specifying paths on WINDOWS, making it universal


	for _,f := range os.Args[2:] {

		if curOption != "" {
			// in the process of getting a command option
			if curOption == "x" {
				for _,s := range strings.Split(f, string([]byte{os.PathListSeparator})) {
					excludes[s] = true
				}
			}
			curOption = ""
			continue
		}
		if f[0:1] == "-" {
			curOption = f[1:]
			continue
		}

		f, err = sys.GetModulePath(f, modules)

		if cmd == "mkdir" {
			files = append(files, f)
		} else {
			if files2,_ := filepath.Glob(f); files2 != nil {
				for _, file2 := range files2 {
					_,fName := filepath.Split(file2)
					if file2 != "" &&
						fName[0] != '.' && // ignore dot files
						!excludes[fName] { // exclude specific files and directories
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
			sys.ExecuteShellCommand("go generate " + f)
		}
	case "copy":
		copyFiles(files)

	case "mkdir":
		makeDirectory(files)
	}
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
				err = sys.DirectoryCopy(file, dest)
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
			err = sys.FileCopy(file, destFile)
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


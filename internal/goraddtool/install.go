package goraddtool

import (
	"bufio"
	"fmt"
	sys2 "github.com/goradd/gofile/pkg/sys"
	install2 "github.com/goradd/goradd/internal/install"
	"github.com/goradd/goradd/pkg/sys"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const initCmd = "go mod init goradd-project" // Create go.mod file

var commands = []string{
	"go mod tidy", // Setup go.mod file
	"go install github.com/goradd/got@latest",                    // Template processor
	"go install github.com/goradd/moddoc@latest",                 // Documentation generator
	"go install golang.org/x/tools/cmd/goimports@latest",         // For auto-fixing of import declarations
	"go install github.com/goradd/gofile@latest",                 // For deployment
	"go install github.com/tdewolff/minify/v2/cmd/minify@latest", // for deployment
	"go mod tidy", // cleanup install
}

// install will copy the project and tmp directories to the cwd
func install(step int, overwrite bool, dependencies bool) {
	if dependencies {
		depInstall()
		return
	}

	loadCwd()

	switch step {
	case 0:
		copyInstall(overwrite)
		depInstall()
	case 1:
		copyInstall(overwrite)
	case 2:
		depInstall()
	}
}

// copyInstall copies the contents of the internal install directory to the current working directory.
func copyInstall(overwrite bool) {
	dirEntries, err := os.ReadDir(install2.InstallFolderLocation)
	if err != nil {
		log.Fatal(fmt.Errorf("could not read the install directory: %s", err.Error()))
	}

	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		dest := filepath.Join(cwd, dirEntry.Name())
		fmt.Println("Copying " + dirEntry.Name() + " ...")
		if sys.PathExists(dest) {
			if !overwrite {
				fmt.Printf("\n*** The %s directory already exists. Replace it? [y,n] ", dest)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				in := scanner.Text()
				if in != "y" {
					return
				}
			}
			err = os.RemoveAll(dest)
			if err != nil {
				log.Fatal("could not remove directory: " + err.Error())
			}
		}
		err = sys2.CopyDirectory(filepath.Join(install2.InstallFolderLocation, dirEntry.Name()), cwd, sys2.CopyDoNotOverwrite)
		if err != nil {
			log.Fatal("could not copy directory: " + err.Error())
		}
	}

	// When goradd is installed, all its files are read-only. Copying the directory will copy these files as
	// read-only as well, which is not what we want, since the project directory is something we want the user to edit.
	// So, we recursively make all the project files read-write.
	err = filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			_ = os.Chmod(path, 0644)
		}
		return nil
	})
	if err != nil {
		log.Fatal("could not walk install directory: " + err.Error())
	}

	err = os.Chdir(filepath.Join(cwd, "goradd-project"))
	if err != nil {
		log.Fatal("could not change to the goradd-project directory: " + err.Error())
		os.Exit(1)
	}

	c := initCmd
	fmt.Print("Executing " + c + " ")
	var res []byte
	res, err = sys2.ExecuteShellCommand(c)
	fmt.Println(string(res))
}

// depInstall installs the dependencies.
//
// Since some dependencies are contained in the go.mod file, this should be run from the same
// directory as the go.mod file.
func depInstall() {
	var err error

	// install binary commands
	for _, c := range commands {
		var res []byte
		fmt.Print("Executing " + c + " ")
		res, err = sys2.ExecuteShellCommand(c)
		fmt.Println(string(res))
		if err != nil {
			fmt.Println(string(err.(*exec.ExitError).Stderr))
			// the error message that was generated
			log.Fatal("could not execute " + c + " " + err.Error())
		}
	}
}

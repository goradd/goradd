package goraddtool

import (
	"bufio"
	"fmt"
	sys2 "github.com/goradd/gofile/pkg/sys"
	install2 "github.com/goradd/goradd/internal/install"
	"github.com/goradd/goradd/pkg/sys"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var dependencies = []string {
	"github.com/goradd/got/...",				// Template processor
	"golang.org/x/tools/cmd/goimports",		// For auto-fixing of import declarations
	"github.com/goradd/gofile/...",				// For deployment
	"github.com/goradd/gengen/...",				// For creation of generics
}

// install will copy the project and tmp directories to the cwd
func install(step int, overwrite bool) {
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

func copyInstall(overwrite bool) {
	var err error
	var fInfos []os.FileInfo

	if fInfos, err = ioutil.ReadDir(install2.InstallFolderLocation); err != nil {
		log.Fatal(fmt.Errorf("could not read the install directory: %s", err.Error()))
	}

	for _, fInfo := range fInfos {
		if !fInfo.IsDir() {
			continue
		}
		dest := filepath.Join(cwd, fInfo.Name())
		if sys.PathExists(dest) {
			if !overwrite {
				fmt.Printf("The %s directory already exists. Replace it? [y,n] ", dest)
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				in := scanner.Text()
				if in != "y" {
					os.Exit(0)
				}
			}
			err = os.RemoveAll(dest)
			if err != nil {
				log.Fatal("could not remove directory: " + err.Error())
			}
		}
		err = sys2.CopyDirectory(filepath.Join(install2.InstallFolderLocation, fInfo.Name()), cwd, sys2.CopyDoNotOverwrite)
		if err != nil {
			log.Fatal("could not copy directory: " + err.Error())
		}
	}
}

func depInstall() {
	var err error

	// install binary dependencies
	err = os.Chdir(filepath.Join(cwd, "goradd-project")) // For module-aware mode in go v1.11, we have to be in a directory with a go.mod file
	// Release notes for go v1.12 appear to indicate they have fixed this.
	if err != nil {
		log.Fatal("could not change to the goradd-project directory: " + err.Error())
	}

	for _,dep := range dependencies {
		var res []byte
		res, err = sys2.ExecuteShellCommand("go get " + dep)
		if err != nil {
			fmt.Print(string(res))
			fmt.Print(string(err.(*exec.ExitError).Stderr))
			// the error message that was generated
			log.Fatal("could not get " + dep + " " + err.Error())
		}
		fmt.Print(res)
	}
}


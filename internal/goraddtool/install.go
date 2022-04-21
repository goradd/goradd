package goraddtool

import (
	"bufio"
	"fmt"
	sys2 "github.com/goradd/gofile/pkg/sys"
	install2 "github.com/goradd/goradd/internal/install"
	"github.com/goradd/goradd/pkg/sys"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var dependencies = []string{
	"github.com/goradd/got/...",                // Template processor
	"golang.org/x/tools/cmd/goimports",         // For auto-fixing of import declarations
	"github.com/goradd/gofile/...",             // For deployment
	"github.com/tdewolff/minify/v2/cmd/minify", // for deployment
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

// copyInstall copies the contents of the internal install directory to the current working directory.
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
		fmt.Println("Copying " + fInfo.Name() + " ...")
		if sys.PathExists(dest) {
			if !overwrite {
				fmt.Printf("\n*** The %s directory already exists. Replace it? [y,n] ", dest)
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

		dest2 := filepath.Join(dest, "gomod.txt")
		dest3 := filepath.Join(dest, "go.mod")
		err = os.Rename(dest2, dest3)
		if err != nil {
			log.Fatal("could not rename go.mod: " + err.Error())
		}
	}

	// When goradd is installed, all its files are read-only. Copying the directory will copy these files as
	// read-only as well, which is not what we want, since the project directory is something we want the user to edit.
	// So, we recursively make all the project files read only.
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

}

func depInstall() {
	var err error

	// install binary dependencies
	err = os.Chdir(filepath.Join(cwd, "goradd-project"))
	if err != nil {
		log.Fatal("could not change to the goradd-project directory: " + err.Error())
	}

	if res, err := sys2.ExecuteShellCommand("go mod tidy"); err != nil {
		fmt.Print(string(res))
		fmt.Print(string(err.(*exec.ExitError).Stderr))
		// the error message that was generated
		log.Fatal("could not run go mod tidy")
	}

	for _, dep := range dependencies {
		var res []byte
		fmt.Print("Installing " + dep + " ")
		res, err = sys2.ExecuteShellCommand("go install " + dep + "@latest")
		if err != nil {
			fmt.Print(string(res))
			fmt.Print(string(err.(*exec.ExitError).Stderr))
			// the error message that was generated
			log.Fatal("could not get " + dep + " " + err.Error())
		}
		fmt.Println(string(res))
	}
}

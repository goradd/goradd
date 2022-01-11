package goraddtool

import (
	"fmt"
	sys2 "github.com/goradd/gofile/pkg/sys"
	"github.com/goradd/goradd/internal/ci"
	"github.com/goradd/goradd/pkg/sys"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func test(step int, browser bool, headless bool) {
	loadCwd()

	switch step {
	case 0:
		copyTestDir()
		testCodegen()
		pkgTest()
		dbTest()
		browserTest(browser)
	case 1:
		copyTestDir()
	case 2:
		testCodegen()
	case 3:
		pkgTest()
	case 4:
		dbTest()
	case 5:
		browserTest(browser)
	}
}

func copyTestDir() {
	var err error
	var fInfos []os.FileInfo

	if fInfos, err = ioutil.ReadDir(cwd); err != nil {
		log.Fatal(fmt.Errorf("could not read the current directory: %s", err.Error()))
	}

	// make sure we have a goradd-project directory, or the tests will not work
	var hasProj bool
	for _, fInfo := range fInfos {
		if fInfo.Name() == "goradd-project" && fInfo.IsDir() {
			hasProj = true
		}
	}
	if !hasProj {
		log.Fatal("Could not find a goradd-project directory in the current working directory")
	}

	dest := filepath.Join(cwd, "goradd-test")
	err = os.RemoveAll(dest)
	if err != nil {
		log.Fatal("could not remove goradd-test directory: " + err.Error())
	}

	err = sys2.CopyDirectory(filepath.Join(ci.TestFolderLocation, "goradd-test"), cwd, sys2.CopyDoNotOverwrite)
	if err != nil {
		log.Fatal("could not copy directory: " + err.Error())
	}
}

func testCodegen() {
	var err error

	err = os.Chdir(filepath.Join(cwd, "goradd-test", "codegen")) // For module-aware mode in go v1.11, we have to be in a directory with a go.mod file
	// Release notes for go v1.12 appear to indicate they have fixed this.
	if err != nil {
		log.Fatal("could not change to the goradd-test/codegen directory: " + err.Error())
	}

	var cmd = "go generate build.go"
	var generateResult []byte

	generateResult, err = sys2.ExecuteShellCommand(cmd)
	fmt.Print(string(generateResult))
	if err != nil {
		log.Fatal("could not generate code: " + err.Error())
	}
}

func pkgTest() {
	var err error

	err = os.Chdir(filepath.Join(cwd, "goradd-project")) // For module-aware mode in go v1.11, we have to be in a directory with a go.mod file
	if err != nil {
		log.Fatal("could not change to the goradd-project directory: " + err.Error())
	}

	cmd := "go test github.com/goradd/goradd/pkg/..."
	var testResult []byte
	testResult, err = sys2.ExecuteShellCommand(cmd)
	fmt.Print(string(testResult))
	if err != nil {
		log.Fatal("pkg unit test failed: " + err.Error())
	}
}

func dbTest() {
	var err error

	err = os.Chdir(filepath.Join(cwd, "goradd-project")) // For module-aware mode in go v1.11, we have to be in a directory with a go.mod file
	if err != nil {
		log.Fatal("could not change to the goradd-project directory: " + err.Error())
	}

	cmd := "go test github.com/goradd/goradd/test/dbtest"
	var testResult []byte
	testResult, err = sys2.ExecuteShellCommand(cmd)
	fmt.Print(string(testResult))
	if err != nil {
		log.Fatal("dbtest failed: " + err.Error())
	}

}

func browserTest(browser bool) {
	if browser {
		if err := sys.LaunchChrome("http://localhost:8000/goradd/Test.g?all=1"); err != nil {
			log.Fatal(err)
		}
	}

	if err := os.Chdir(filepath.Join(cwd, "goradd-test")); err != nil {
		log.Fatal("could not change to the goradd-test directory: " + err.Error())
	}

	cmd := "go run main.go"
	var result []byte
	var err error
	result, err = sys2.ExecuteShellCommand(cmd)
	fmt.Print(string(result))
	if err != nil {
		log.Fatal("browser test failed: " + err.Error())
	}

}

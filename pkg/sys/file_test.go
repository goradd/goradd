package sys

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

/*
func TestFileCopy(t *testing.T) {
	dir := os.TempDir()
	src := filepath.Join(dir, "test.txt")

	testContent := "I am a test \n and the 2nd line"

	defer os.Remove(src)
	err := ioutil.WriteFile(src, []byte(testContent), 0755)
	if err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "test2.txt")
	defer os.Remove(dst)
	err = FileCopy(src,dst)
	if err != nil {
		t.Fatal(err)
	}

	bytes,err := ioutil.ReadFile(dst)

	if string(bytes) != testContent {
		t.Error("File content does not match")
	}

	// test grow
	testContent = "I am a test \n and the 2nd line\nand a 3rd"
	err = ioutil.WriteFile(src, []byte(testContent), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = FileCopy(src,dst)
	if err != nil {
		t.Fatal(err)
	}

	bytes,err = ioutil.ReadFile(dst)

	if string(bytes) != testContent {
		t.Error("File content does not match")
	}

	// test shrink
	testContent = "I am a different test"
	err = ioutil.WriteFile(src, []byte(testContent), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = FileCopy(src,dst)
	if err != nil {
		t.Fatal(err)
	}

	bytes,err = ioutil.ReadFile(dst)

	if string(bytes) != testContent {
		t.Error("File content does not match")
	}
}
*/
func TestFileModDateCompare(t *testing.T) {
	f1,err := ioutil.TempFile("","test")
	if err != nil {
		t.Fatal(err)
	}
	name1 := f1.Name()
	defer os.Remove(name1)
	testContent := "I am a test \n and the 2nd line"
	f1.WriteString(testContent)
	f1.Close()

	time.Sleep(2 * time.Second)

	f2,err := ioutil.TempFile("","test")
	if err != nil {
		t.Fatal(err)
	}
	name2 := f2.Name()
	defer os.Remove(name2)
	f2.WriteString(testContent)
	f2.Close()

	v,err := FileModDateCompare(name1, name2)
	if err != nil {
		t.Fatal(err)
	}
	if v !=-1 {
		t.Error("First file is not earlier than second file.")
	}

}

/*
func TestDirectoryCopy(t *testing.T) {
	dir1 := filepath.Join(os.TempDir(), "dir1")
	dir2 := filepath.Join(os.TempDir(), "dir2")


	os.Mkdir(dir1, 0777)
	os.Mkdir(dir2, 0777)

	defer os.RemoveAll(dir1)
	defer os.RemoveAll(dir2)

	subdir := filepath.Join(dir1, "subdir1")
	if err := os.Mkdir(subdir, 0777); err != nil {
		t.Fatal(err)
	}

	// set up the test directory
	testContent := "I am a test"
	if err := ioutil.WriteFile(filepath.Join(dir1, "test1"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir1, "test2"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(subdir, "test3"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}

	if err := DirectoryCopy(dir1, dir2); err != nil {
		t.Fatal(err)
	}

	items,err := ioutil.ReadDir(dir2)
	if err != nil {
		t.Fatal(err)
	}
	if items[0].Name() != "dir1" {
		t.Fatal("First item in directory is not dir1")
	}
	items,err = ioutil.ReadDir(filepath.Join(dir2, "dir1"))
	if err != nil {
		t.Fatal(err)
	}

	if items[0].Name() != "subdir1" {
		t.Fatal("First item in directory is not subdir1, but rather: " + items[0].Name())
	}
	if items[1].Name() != "test1" {
		t.Fatal("Second item in directory is not test1, but rather: " + items[1].Name())
	}
}
*/

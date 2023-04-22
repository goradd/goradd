package sys

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// FileModDateCompare compares the modification date of two files, and returns -1 if the first is older than the second,
// 0 if they are the same, and 1 if the 2nd is older than the first. Returns an error if either is not a file.
func FileModDateCompare(file1, file2 string) (diff int, err error) {
	file1Info, err := os.Stat(file1)
	if err != nil {
		return
	}
	if file1Info.Mode().IsDir() {
		err = fmt.Errorf("%s is a directory, not a file. \n", file1)
		return
	}

	file2Info, err := os.Stat(file2)
	if err != nil {
		return
	}
	if file2Info.Mode().IsDir() {
		err = fmt.Errorf("%s is a directory, not a file. \n", file2)
		return
	}

	modTime1 := file1Info.ModTime()
	modTime2 := file2Info.ModTime()

	diff2 := modTime1.Sub(modTime2)

	if diff2 == (time.Duration(0) * time.Second) {
		diff = 0
	} else if diff2 < (time.Duration(0) * time.Second) {
		diff = -1
	} else {
		diff = 1
	}
	return
}

// IsDir returns true if the given path exists and is a directory
func IsDir(path string) bool {
	dstInfo, err := os.Stat(path)
	if err == nil { // file exists
		if dstInfo.Mode().IsDir() {
			return true
		}
	}
	return false
}

/* FileCopyIfNewer has been replaced by code in gofile
func FileCopyIfNewer(src, dst string) (err error) {
	var diff int

	dstInfo, err := os.Stat(dst)
	if err == nil { // file exists
		if dstInfo.Mode().IsDir() {
			return fmt.Errorf("%s is a directory, not a file. \n", dst)
		}

		if diff, err = FileModDateCompare(src, dst); diff != 1 || err != nil {
			// either src is older or the same as dst
			return
		}
	}

	return FileCopy(src, dst)
}
*/

// PathExists returns true if the given path exists in the OS. This does not necessarily mean that the path is
// usable. It may be write or read protected. But at least you know its there.
func PathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil || !os.IsNotExist(err)
}

// SourcePath is a utility function that returns the path to the file that called it.
func SourcePath() (file string) {
	_, file, _, _ = runtime.Caller(1)
	return
}

// SourceDirectory is a utility function that returns the path to the directory of the file that called it.
func SourceDirectory() (dir string) {
	var file string
	_, file, _, _ = runtime.Caller(1)
	return filepath.Dir(file)
}

// FileRoot returns the name of a file without any extensions.
func FileRoot(filename string) string {
	offset := strings.Index(filename, ".")
	if offset < 0 {
		return filename
	} else {
		return filename[0:offset]
	}
}

/* DirectoryCopy has been replaced by code in gofile
func DirectoryCopy(src, dst string) (err error) {
	dstInfo, dstErr := os.Stat(dst)
	srcInfo, srcErr := os.Stat(src)

	if srcErr != nil {
		return fmt.Errorf("source directory error: %s", srcErr.Error())
	}

	if dstErr != nil {
		return fmt.Errorf("destination directory error: %s", dstErr.Error())
	}

	if len(src) <= len(dst) && dst[:len(src)] == src { // does dst start with src?
		return fmt.Errorf("destination directory is not allowed to be in the src directory")
	}

	if !dstInfo.Mode().IsDir() {
		return fmt.Errorf("source %s is a file, not a directory", dst)
	}

	// create destination if needed
	newPath := filepath.Join(dst, filepath.Base(src))

	if !PathExists(newPath) {
		perm := srcInfo.Mode().Perm()	// copy the permission
		err = os.Mkdir(newPath, perm)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %s", newPath, err.Error())
		}
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	list, err := f.Readdir(-1)
	f.Close()

	for _,item := range list {
		itemName := item.Name()
		itemPath := filepath.Join(src, itemName)
		if item.IsDir() {
			if itemName != "." && itemName != ".." {
				DirectoryCopy(itemPath, newPath)
			}
		} else {
			newItemPath := filepath.Join(newPath, itemName)
			dstFileInfo, dstFileErr := os.Stat(newItemPath)
			if dstFileErr != nil {
				if os.IsNotExist(dstFileErr) {
					err = FileCopy(itemPath, newItemPath)
					if err != nil {
						return
					}
				} else {
					return dstFileErr
				}
			} else {
				if dstFileInfo.IsDir() {
					return fmt.Errorf("Path %s is a file in the source, but %s is a directory in the destination.", itemPath, newItemPath)
				}
				// otherwise do no copying since the file already exists
			}
		}
	}

	return
}
*/

/* DirectoryClear has been replaced by code in gofile
func DirectoryClear(dir string) error {
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _,item := range items {
		os.RemoveAll(filepath.Join(dir, item.Name()))
	}
	return nil
}
*/

/* FileCopy has been replaced by the code in gofile
func FileCopy(src, dst string) (err error) {
	var count int64

	srcInfo, srcErr := os.Stat(src)
	if srcErr != nil {
		return srcErr
	}
	perm := srcInfo.Mode() & os.ModePerm

	from, err := os.Open(src)
	if err != nil {
		return
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, perm) // copy source permissions
	if err != nil {
		return
	}

	defer to.Close()

	count, err = io.Copy(to, from)
	if err != nil {
		to.Close()
		return err
	}
	to.Truncate(count) // chop end of file in case file gets smaller

	return to.Close()
}
*/

type WriterToFunc func(w io.Writer) (int64, error)

func (f WriterToFunc) WriteTo(w io.Writer) (n int64, err error) {
	return f(w)
}

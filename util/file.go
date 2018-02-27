package util

import (
	"os"
	"io"
	"fmt"
	"time"
)

// Copy copies the src file to the destination. The destination file must either exist, or the directory in the file's
// path must exist.
func FileCopy(src, dst string) (err error) {
	from, err := os.Open(src)
	if err != nil {
		return
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}

	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		to.Close()
		return err
	}

	return to.Close()
}

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

// CopyIfNewer performs a copy to the destination if the src is newer than the destination, or the destination does
// not exist.
func FileCopyIfNewer(src, dst string) (err error) {
	var diff int

	dstInfo, err := os.Stat(dst)
	if err == nil { // file exists
		if dstInfo.Mode().IsDir() {
			return fmt.Errorf("%s is a directory, not a file. \n", dst)
		}

		if diff,err = FileModDateCompare(src, dst); diff != 1 || err != nil {
			// either src is older or the same as dst
			return
		}
	}

	return FileCopy(src, dst)
}
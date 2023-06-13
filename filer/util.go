package filer

import (
	"fmt"
	"os"
	"time"
)

// --------------------------------------util function.-----------------------------
func isFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}

	return false, err
}

func getFileModifyTime(file_path string) time.Time {
	finfo, err := os.Stat(file_path)
	if err != nil {
		fmt.Println("file is not exist:", file_path)
		return time.Unix(0, 0)
	}
	return finfo.ModTime()
}

func isSameDate(t time.Time, nt time.Time) bool {
	if t.Day() != nt.Day() {
		return false
	}

	if t.Month() != nt.Month() {
		return false
	}

	if t.Year() != nt.Year() {
		return false
	}
	return true
}

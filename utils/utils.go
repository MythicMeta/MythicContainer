package utils

import (
	"log"
	"os"
	"path/filepath"
)

func StringSliceContains(source []string, str string) bool {
	for _, v := range source {
		if str == v {
			return true
		}
	}
	return false
}

func RemoveStringFromSliceNoOrder(source []string, str string) []string {
	for index, value := range source {
		if str == value {
			source[index] = source[len(source)-1]
			source[len(source)-1] = ""
			source = source[:len(source)-1]
			return source
		}
	}
	// we didn't find the element to remove
	return source
}

func GetCwdFromExe() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("[-] Failed to get path to current executable: %v", err)
	}
	return filepath.Dir(exe)
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
}

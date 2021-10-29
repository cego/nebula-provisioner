package main

import "os"

func fileExists(filepath string) (bool, os.FileInfo) {

	fileinfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false, fileinfo
	}
	// Return false if the fileinfo says the file path is a directory.
	return !fileinfo.IsDir(), fileinfo
}

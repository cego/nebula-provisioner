package main

import (
	"os"
	"path/filepath"
	"sort"
)

func fileExists(filepath string) (bool, os.FileInfo) {

	fileinfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false, fileinfo
	}
	// Return false if the fileinfo says the file path is a directory.
	return !fileinfo.IsDir(), fileinfo
}

func dirExists(path string) (bool, os.FileInfo) {

	fileinfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, fileinfo
	}
	// Return false if the fileinfo says the path is not a directory.
	return fileinfo.IsDir(), fileinfo
}

func getConfigPath() string {
	ok, _ := dirExists("/etc/nebula-provisioner")
	if ok {
		return "/etc/nebula-provisioner"
	}

	ok, _ = dirExists("/opt/nebula-provisioner/etc")
	if ok {
		return "/opt/nebula-provisioner/etc"
	}

	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, "agent.yml")
}

func resolvePath(path string) string {
	if filepath.Dir(path) == "." {
		path = filepath.Join(configDir, path)
	}
	return path
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

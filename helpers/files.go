package helpers

import (
	"os"
	"path/filepath"
)

func GetFilenameWithoutExtension(path string) string {
	return filepath.Base(path)[0 : len(filepath.Base(path))-len(filepath.Ext(path))]
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		// We could check with os.IsNotExist(err) here, but since os.Stat threw an error, we likely can't use the file anyway.
		return false
	}
	return true
}

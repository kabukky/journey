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

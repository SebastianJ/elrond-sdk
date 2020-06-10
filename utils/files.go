package utils

import (
	"fmt"
	"path/filepath"
)

// IdentifyPemFiles - identify pem files from a specified path
func IdentifyPemFiles(path string) ([]string, error) {
	pattern := fmt.Sprintf("%s/*.pem", path)

	fmt.Println("Pem pattern is now: ", pattern)

	keys, err := globFiles(pattern)

	if err != nil {
		return nil, err
	}

	return keys, err
}

func globFiles(pattern string) ([]string, error) {
	files, err := filepath.Glob(pattern)

	if err != nil {
		return nil, err
	}

	return files, nil
}

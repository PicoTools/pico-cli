package utils

import "path/filepath"

// GetAbsPath returns absolute path from relative path
func GetAbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

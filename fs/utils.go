package fs

import (
	"fmt"
	"os"
)

func ValidateDirectory(dirPath string) error {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist at %v", dirPath)
	}
	if err != nil {
		return fmt.Errorf("directory error at %v: %w", dirPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("directory is actually a file at %v", dirPath)
	}
	return nil
}

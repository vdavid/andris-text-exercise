package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// findRootDir walks up from the current directory looking for a marker file
// that indicates the project root (e.g., package.json, .git, AGENTS.md).
func findRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory: %w", err)
	}

	for {
		// Check for project root markers
		markers := []string{"AGENTS.md", "package.json", ".git"}
		for _, marker := range markers {
			if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
				return dir, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fall back: go.mod is in scripts/check/, so root is two levels up
	execDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory: %w", err)
	}
	return filepath.Clean(filepath.Join(execDir, "..", "..")), nil
}

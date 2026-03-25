package checks

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxFileLines = 300

// ignoredDirs contains directories to skip during the file length check.
var ignoredDirs = map[string]bool{
	"node_modules": true,
	".next":        true,
	".git":         true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
}

// ignoredFiles contains specific filenames to skip during the file length check.
var ignoredFiles = map[string]bool{
	"package-lock.json": true,
}

// checkedExtensions contains file extensions to check.
var checkedExtensions = map[string]bool{
	".ts":   true,
	".tsx":  true,
	".js":   true,
	".jsx":  true,
	".go":   true,
	".css":  true,
	".json": true,
	".md":   true,
}

// FileLength checks that no source file exceeds the maximum line count.
var FileLength = Check{
	Name: "file-length",
	App:  "",
	Run: func(rootDir string) error {
		var violations []string

		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip ignored directories
			if info.IsDir() {
				if ignoredDirs[info.Name()] {
					return filepath.SkipDir
				}
				return nil
			}

			// Skip ignored files
			if ignoredFiles[info.Name()] {
				return nil
			}

			// Only check known extensions
			ext := filepath.Ext(path)
			if !checkedExtensions[ext] {
				return nil
			}

			lines, countErr := countLines(path)
			if countErr != nil {
				return nil // skip files we can't read
			}

			if lines > maxFileLines {
				relPath, _ := filepath.Rel(rootDir, path)
				violations = append(violations, fmt.Sprintf("  %s: %d lines (max %d)", relPath, lines, maxFileLines))
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		if len(violations) > 0 {
			return fmt.Errorf("files exceeding %d lines:\n%s", maxFileLines, strings.Join(violations, "\n"))
		}

		return nil
	},
}

func countLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}

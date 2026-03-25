package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const statsFileName = ".check-stats.csv"

// LogCheckResult appends a check result to the CSV stats file.
func LogCheckResult(rootDir string, checkName string, passed bool, duration time.Duration) {
	statsPath := filepath.Join(rootDir, statsFileName)

	file, err := os.OpenFile(statsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return // silently ignore stats errors
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	status := "pass"
	if !passed {
		status = "fail"
	}

	_ = w.Write([]string{
		time.Now().UTC().Format(time.RFC3339),
		checkName,
		status,
		fmt.Sprintf("%.3f", duration.Seconds()),
	})
}

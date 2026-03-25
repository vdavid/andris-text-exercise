package checks

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ScriptsGoGofmt checks that all Go files in scripts/ are properly formatted.
var ScriptsGoGofmt = Check{
	Name: "scripts-go-gofmt",
	App:  "scripts",
	Run: func(rootDir string) error {
		scriptsDir := filepath.Join(rootDir, "scripts", "check")
		result, err := RunCommand(scriptsDir, "gofmt", "-l", ".")
		if err != nil {
			return err
		}

		unformatted := strings.TrimSpace(result.Stdout)
		if unformatted != "" {
			return fmt.Errorf("unformatted Go files:\n%s\nRun: gofmt -w scripts/check/", unformatted)
		}

		return nil
	},
}

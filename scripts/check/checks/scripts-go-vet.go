package checks

import (
	"path/filepath"
)

// ScriptsGoVet runs go vet on the check scripts.
var ScriptsGoVet = Check{
	Name:      "scripts-go-vet",
	App:       "scripts",
	DependsOn: []string{"scripts-go-gofmt"},
	Run: func(rootDir string) error {
		scriptsDir := filepath.Join(rootDir, "scripts", "check")
		_, err := RunCommand(scriptsDir, "go", "vet", "./...")
		return err
	},
}

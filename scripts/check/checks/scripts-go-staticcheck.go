package checks

import (
	"path/filepath"
)

// ScriptsGoStaticcheck runs staticcheck on the check scripts.
var ScriptsGoStaticcheck = Check{
	Name:      "scripts-go-staticcheck",
	App:       "scripts",
	DependsOn: []string{"scripts-go-vet"},
	Run: func(rootDir string) error {
		scriptsDir := filepath.Join(rootDir, "scripts", "check")
		_, err := RunCommand(scriptsDir, "staticcheck", "./...")
		return err
	},
}

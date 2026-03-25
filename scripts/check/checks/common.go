package checks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Check represents a single check to run.
type Check struct {
	Name      string
	App       string   // "app" for the Next.js app, "scripts" for Go check scripts
	DependsOn []string // Names of checks that must pass before this one runs
	Run       func(rootDir string) error
}

// CommandResult holds the output of a command execution.
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// RunCommand executes a command in the given directory and returns the result.
func RunCommand(dir string, name string, args ...string) (CommandResult, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
		return result, fmt.Errorf("command %q exited with code %d:\n%s%s",
			strings.Join(append([]string{name}, args...), " "),
			result.ExitCode,
			result.Stdout,
			result.Stderr,
		)
	}

	if err != nil {
		return result, fmt.Errorf("command %q failed: %w", name, err)
	}

	return result, nil
}

var (
	npmInstallOnce sync.Once
	npmInstallErr  error
)

// EnsureNpmDependencies ensures that npm dependencies are installed.
func EnsureNpmDependencies(rootDir string) error {
	npmInstallOnce.Do(func() {
		// Check if node_modules exists
		nodeModules := filepath.Join(rootDir, "node_modules")
		if _, err := os.Stat(nodeModules); err == nil {
			return
		}

		// Check if package.json exists
		packageJSON := filepath.Join(rootDir, "package.json")
		if _, err := os.Stat(packageJSON); err != nil {
			npmInstallErr = fmt.Errorf("package.json not found at %s", packageJSON)
			return
		}

		_, npmInstallErr = RunCommand(rootDir, "npm", "install")
	})
	return npmInstallErr
}

// runPrettierCheck runs Prettier in check mode on the given paths.
func RunPrettierCheck(rootDir string, paths ...string) error {
	if err := EnsureNpmDependencies(rootDir); err != nil {
		return fmt.Errorf("failed to ensure dependencies: %w", err)
	}

	args := append([]string{"prettier", "--check"}, paths...)
	_, err := RunCommand(rootDir, "npx", args...)
	return err
}

// RunESLintCheck runs ESLint on the given paths.
func RunESLintCheck(rootDir string, paths ...string) error {
	if err := EnsureNpmDependencies(rootDir); err != nil {
		return fmt.Errorf("failed to ensure dependencies: %w", err)
	}

	args := append([]string{"eslint"}, paths...)
	_, err := RunCommand(rootDir, "npx", args...)
	return err
}

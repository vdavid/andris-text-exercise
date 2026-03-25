# Check System Architecture

## Overview

A Go-based parallel check runner that validates code quality across the project. Invoked via `./scripts/check.sh` from the project root.

## Structure

```
scripts/check/
├── main.go           # CLI entry point: flag parsing, check selection, summary
├── runner.go         # Parallel executor with dependency gating and fail-fast
├── colors.go         # ANSI color constants
├── stats.go          # CSV stats logging (writes to .check-stats.csv at root)
├── utils.go          # findRootDir() helper
├── go.mod            # Go module definition
├── CLAUDE.md         # This file
└── checks/
    ├── common.go     # Check type, RunCommand, EnsureNpmDependencies, lint/format helpers
    ├── registry.go   # AllChecks list and lookup functions
    └── *.go          # Individual check definitions (one per file)
```

## How It Works

1. `main.go` parses flags and selects which checks to run.
2. `runner.go` runs selected checks in parallel (goroutine pool), respecting dependency ordering.
3. Each check is a `Check` struct with a `Run func(rootDir string) error`. Return `nil` for pass, an error for fail.
4. Results are printed with color-coded pass/fail/skip indicators and duration.

## Adding a New Check

1. Create a new file in `checks/` (e.g., `checks/app-tests.go`).
2. Define a package-level `var` of type `Check`:
   ```go
   var AppTests = Check{
       Name:      "app-tests",
       App:       "app",           // "app" or "scripts"
       DependsOn: []string{"app-typecheck"}, // optional
       Run: func(rootDir string) error {
           _, err := RunCommand(rootDir, "npx", "jest")
           return err
       },
   }
   ```
3. Add it to `AllChecks` in `checks/registry.go`.

## Apps

- `app` - The Next.js application. Checks: ESLint, Prettier, TypeScript type-check.
- `scripts` - The Go check runner itself. Checks: gofmt, go vet, staticcheck.

## Key Flags

- `-check <name>` - Run a specific check (comma-separated for multiple)
- `-app <name>` - Run all checks for an app
- `-j <n>` - Parallelism (defaults to CPU count)
- `-fail-fast` - Stop after the first failure
- `-list` - List all checks

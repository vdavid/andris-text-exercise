package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"andris-text-exercise/scripts/check/checks"
)

// CheckResult holds the outcome of a single check.
type CheckResult struct {
	Name     string
	Passed   bool
	Skipped  bool
	Duration time.Duration
	Error    error
}

// Runner executes checks with parallelism, dependency gating, and fail-fast.
type Runner struct {
	Checks   []checks.Check
	Jobs     int
	FailFast bool
	RootDir  string

	mu        sync.Mutex
	started   map[string]bool         // checks that have been dispatched
	completed map[string]*CheckResult // checks that have finished
	failed    bool
}

// NewRunner creates a runner for the given checks.
func NewRunner(selectedChecks []checks.Check, jobs int, failFast bool, rootDir string) *Runner {
	return &Runner{
		Checks:    selectedChecks,
		Jobs:      jobs,
		FailFast:  failFast,
		RootDir:   rootDir,
		started:   make(map[string]bool),
		completed: make(map[string]*CheckResult),
	}
}

// Run executes all checks and returns the results.
func (r *Runner) Run() []CheckResult {
	// Build a set of check names we're running
	checkSet := make(map[string]bool)
	for _, c := range r.Checks {
		checkSet[c.Name] = true
	}

	// Channel for checks that are ready to run
	ready := make(chan checks.Check, len(r.Checks))

	// enqueue dispatches checks whose dependencies are all completed and passed.
	enqueue := func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		for _, c := range r.Checks {
			if r.started[c.Name] {
				continue
			}

			depsOk := true
			for _, dep := range c.DependsOn {
				if !checkSet[dep] {
					continue // dependency not in our run set, ignore
				}
				depResult, done := r.completed[dep]
				if !done {
					depsOk = false
					break
				}
				if !depResult.Passed {
					depsOk = false
					break
				}
			}

			if depsOk {
				r.started[c.Name] = true
				ready <- c
			}
		}
	}

	// skipDependents marks checks as skipped if any of their dependencies failed.
	skipDependents := func() []CheckResult {
		r.mu.Lock()
		defer r.mu.Unlock()

		var skipped []CheckResult
		for _, c := range r.Checks {
			if r.started[c.Name] {
				continue
			}
			for _, dep := range c.DependsOn {
				if !checkSet[dep] {
					continue
				}
				depResult, done := r.completed[dep]
				if done && !depResult.Passed {
					result := CheckResult{Name: c.Name, Skipped: true}
					r.started[c.Name] = true
					r.completed[c.Name] = &result
					skipped = append(skipped, result)
					break
				}
			}
		}
		return skipped
	}

	// Start worker pool
	done := make(chan CheckResult, len(r.Checks))
	for i := 0; i < r.Jobs; i++ {
		go func() {
			for c := range ready {
				result := r.runCheck(c)
				done <- result
			}
		}()
	}

	// Initial enqueue
	enqueue()

	// Process results
	var allResults []CheckResult
	remaining := len(r.Checks)

	for remaining > 0 {
		select {
		case result := <-done:
			r.mu.Lock()
			r.completed[result.Name] = &result
			if !result.Passed && !result.Skipped {
				r.failed = true
			}
			r.mu.Unlock()

			r.printResult(&result)
			LogCheckResult(r.RootDir, result.Name, result.Passed, result.Duration)
			allResults = append(allResults, result)
			remaining--

			if r.FailFast && r.failed {
				// Skip all remaining checks
				r.mu.Lock()
				for _, c := range r.Checks {
					if !r.started[c.Name] {
						skipped := CheckResult{Name: c.Name, Skipped: true}
						r.started[c.Name] = true
						r.completed[c.Name] = &skipped
						r.printResult(&skipped)
						allResults = append(allResults, skipped)
						remaining--
					}
				}
				r.mu.Unlock()
				break
			}

			// Skip checks whose dependencies failed
			for _, s := range skipDependents() {
				r.printResult(&s)
				allResults = append(allResults, s)
				remaining--
			}

			// Enqueue newly unblocked checks
			enqueue()

		case <-time.After(5 * time.Minute):
			fmt.Fprintf(os.Stderr, "%s%sTIMEOUT%s All checks timed out\n", colorBold, colorRed, colorReset)
			return allResults
		}
	}

	close(ready)
	return allResults
}

func (r *Runner) runCheck(c checks.Check) CheckResult {
	start := time.Now()
	err := c.Run(r.RootDir)
	duration := time.Since(start)

	return CheckResult{
		Name:     c.Name,
		Passed:   err == nil,
		Duration: duration,
		Error:    err,
	}
}

func (r *Runner) printResult(result *CheckResult) {
	if result.Skipped {
		fmt.Printf("  %s⊘ SKIP%s  %s %s(dependency failed)%s\n",
			colorYellow, colorReset, result.Name, colorDim, colorReset)
		return
	}

	if result.Passed {
		fmt.Printf("  %s✓ PASS%s  %s %s(%s)%s\n",
			colorGreen, colorReset, result.Name, colorDim, result.Duration.Round(time.Millisecond), colorReset)
	} else {
		fmt.Printf("  %s✗ FAIL%s  %s %s(%s)%s\n",
			colorRed, colorReset, result.Name, colorDim, result.Duration.Round(time.Millisecond), colorReset)
		if result.Error != nil {
			lines := strings.Split(result.Error.Error(), "\n")
			for _, line := range lines {
				fmt.Printf("         %s%s%s\n", colorDim, line, colorReset)
			}
		}
	}
}

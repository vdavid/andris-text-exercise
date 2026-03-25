package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"andris-text-exercise/scripts/check/checks"
)

func main() {
	var (
		checkFlag string
		appFlag   string
		jobsFlag  int
		failFast  bool
		listFlag  bool
	)

	flag.StringVar(&checkFlag, "check", "", "Run only the named check (comma-separated for multiple)")
	flag.StringVar(&appFlag, "app", "", "Run checks for the given app only (app, scripts)")
	flag.IntVar(&jobsFlag, "j", runtime.NumCPU(), "Number of parallel checks")
	flag.BoolVar(&failFast, "fail-fast", false, "Stop on first failure")
	flag.BoolVar(&listFlag, "list", false, "List all available checks and exit")
	flag.Parse()

	if listFlag {
		fmt.Println("Available checks:")
		for _, c := range checks.AllChecks {
			deps := ""
			if len(c.DependsOn) > 0 {
				deps = fmt.Sprintf(" (depends on: %s)", strings.Join(c.DependsOn, ", "))
			}
			fmt.Printf("  %-30s [%s]%s\n", c.Name, c.App, deps)
		}
		os.Exit(0)
	}

	rootDir, err := findRootDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not find project root: %v\n", err)
		os.Exit(1)
	}

	// Select checks to run
	selected := selectChecks(checkFlag, appFlag)
	if len(selected) == 0 {
		fmt.Fprintln(os.Stderr, "No checks matched the given filters.")
		os.Exit(1)
	}

	// Print header
	fmt.Printf("\n%s%sRunning %d check(s) with %d worker(s)%s\n\n",
		colorBold, colorCyan, len(selected), jobsFlag, colorReset)

	// Run
	start := time.Now()
	runner := NewRunner(selected, jobsFlag, failFast, rootDir)
	results := runner.Run()

	// Summary
	var passed, failed, skipped int
	for _, r := range results {
		switch {
		case r.Skipped:
			skipped++
		case r.Passed:
			passed++
		default:
			failed++
		}
	}

	duration := time.Since(start)
	fmt.Printf("\n%s%s──────────────────────────────────────%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s%d passed%s", colorBold, colorGreen, passed, colorReset)
	if failed > 0 {
		fmt.Printf(", %s%s%d failed%s", colorBold, colorRed, failed, colorReset)
	}
	if skipped > 0 {
		fmt.Printf(", %s%s%d skipped%s", colorBold, colorYellow, skipped, colorReset)
	}
	fmt.Printf(" %s(%s)%s\n\n", colorDim, duration.Round(time.Millisecond), colorReset)

	if failed > 0 {
		os.Exit(1)
	}
}

func selectChecks(checkFilter, appFilter string) []checks.Check {
	all := checks.AllChecks

	if checkFilter != "" {
		names := strings.Split(checkFilter, ",")
		nameSet := make(map[string]bool)
		for _, n := range names {
			nameSet[strings.TrimSpace(n)] = true
		}

		var selected []checks.Check
		for _, c := range all {
			if nameSet[c.Name] {
				selected = append(selected, c)
			}
		}
		return selected
	}

	if appFilter != "" {
		return checks.FindByApp(appFilter)
	}

	return all
}

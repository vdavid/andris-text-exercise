package checks

// AppESLint runs ESLint on the Next.js application source.
var AppESLint = Check{
	Name: "app-eslint",
	App:  "app",
	Run: func(rootDir string) error {
		return RunESLintCheck(rootDir, ".")
	},
}

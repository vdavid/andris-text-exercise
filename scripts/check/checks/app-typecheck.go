package checks

// AppTypecheck runs the TypeScript compiler in type-check-only mode.
var AppTypecheck = Check{
	Name: "app-typecheck",
	App:  "app",
	Run: func(rootDir string) error {
		if err := EnsureNpmDependencies(rootDir); err != nil {
			return err
		}
		_, err := RunCommand(rootDir, "npx", "tsc", "--noEmit")
		return err
	},
}

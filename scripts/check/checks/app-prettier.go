package checks

// AppPrettier runs Prettier in check mode on the project source files.
var AppPrettier = Check{
	Name: "app-prettier",
	App:  "app",
	Run: func(rootDir string) error {
		return RunPrettierCheck(rootDir, ".")
	},
}

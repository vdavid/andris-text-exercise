package checks

// AllChecks is the ordered list of all checks available in this project.
var AllChecks = []Check{
	AppESLint,
	AppPrettier,
	AppTypecheck,
	FileLength,
	ScriptsGoGofmt,
	ScriptsGoVet,
	ScriptsGoStaticcheck,
}

// FindByName returns the check with the given name, or nil if not found.
func FindByName(name string) *Check {
	for i := range AllChecks {
		if AllChecks[i].Name == name {
			return &AllChecks[i]
		}
	}
	return nil
}

// FindByApp returns all checks for the given app.
func FindByApp(app string) []Check {
	var result []Check
	for _, c := range AllChecks {
		if c.App == app {
			result = append(result, c)
		}
	}
	return result
}

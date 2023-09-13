package env

type attributes struct {
	Force              bool
	EnvironmentFiles   []string
	ErrorOnMissingFile bool
}

type Attribute func(*attributes)

func Force(force bool) func(*attributes) {
	return func(a *attributes) {
		a.Force = force
	}
}

func EnvironmentFiles(files ...string) Attribute {
	return func(a *attributes) {
		a.EnvironmentFiles = files
	}
}

func ErrorOnMissingFile(doErr bool) Attribute {
	return func(a *attributes) {
		a.ErrorOnMissingFile = doErr
	}
}

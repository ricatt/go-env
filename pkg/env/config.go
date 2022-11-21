package env

type Config struct {
	Force              bool
	EnvironmentFiles   []string
	ErrorOnMissingFile bool
}

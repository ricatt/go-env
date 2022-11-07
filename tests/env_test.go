package tests

import (
	"github.com/stretchr/testify/suite"
	"go-env/src/env"
	"os"
	"testing"
)

type Env struct {
	PackageName string `env:"PACKAGE_NAME"`
	LogLevel    string `env:"LOG_LEVEL"`
	Iterations  int    `env:"ITERATIONS"`
	BaseURL     string `env:"BASE_URL"`
	Message     string `env:"MESSAGE"`
}

type EnvOptional struct {
	PackageName string `env:"PACKAGE_NAME"`
	BaseURL     string `env:"BASE_URL"`
}

type EnvInvalidType struct {
	InvalidType any `env:"INVALID_TYPE"`
}

type EnvironmentTestSuite struct {
	suite.Suite
}

func (suite *EnvironmentTestSuite) TestAddValue() {
	var config Env
	err := env.Load(&config, env.Config{
		EnvironmentFile: ".env",
	})
	suite.NoError(err)
	suite.Empty(config.BaseURL)
	suite.Equal("env", config.PackageName)
	suite.Equal("debug", config.LogLevel)
	suite.Equal(10, config.Iterations)
}

func (suite *EnvironmentTestSuite) TestMissingEnvInFile() {
	var config Env
	err := env.Load(&config, env.Config{
		Force:           true,
		EnvironmentFile: ".env",
	})
	suite.Error(err)
	suite.Equal("missing value for BaseURL", err.Error())
}

func (suite *EnvironmentTestSuite) TestOptionalEnvStruct() {
	var (
		config      Env
		envOptional EnvOptional
		err         error
	)
	err = env.Load(&config, env.Config{
		EnvironmentFile: ".env",
	})
	suite.NoError(err)
	err = env.Load(&envOptional, env.Config{
		EnvironmentFile: ".env",
	})
	suite.NoError(err)
}

func (suite *EnvironmentTestSuite) TestInvalidTypeEnvStruct() {
	var (
		config EnvInvalidType
		err    error
	)
	err = env.Load(&config, env.Config{
		EnvironmentFile: ".env",
	})
	suite.Error(err)
	suite.Equal("env: type \"interface\" not supported", err.Error())
}

func (suite *EnvironmentTestSuite) TestInvalidPath() {
	var (
		config EnvInvalidType
		err    error
	)
	err = env.Load(&config, env.Config{
		EnvironmentFile: ".env-invalid-path",
	})
	suite.Error(err)
	suite.Equal("open .env-invalid-path: no such file or directory", err.Error())
}

func (suite *EnvironmentTestSuite) TestNoEnvFile() {
	var (
		config Env
		err    error
	)
	err = os.Setenv("MESSAGE", "Hello World")
	suite.NoError(err)
	err = env.Load(&config, env.Config{})
	suite.NoError(err)
	suite.Equal("Hello World", config.Message)
}

func TestEnvironmentTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}

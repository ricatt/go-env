package tests

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"go-env/src/env"
	"math"
	"os"
	"testing"
)

type Env struct {
	PackageName string `env:"PACKAGE_NAME"`
	LogLevel    string `env:"LOG_LEVEL"`
	Iterations  int    `env:"ITERATIONS"`
	BaseURL     string `env:"BASE_URL"`
	Message     string `env:"MESSAGE"`

	IsTrue bool `env:"IS_TRUE"`

	MaxInt   int   `env:"MAX_INT"`
	MaxUint  uint  `env:"MAX_UINT"`
	MaxInt64  int64  `env:"MAX_INT_64"`
	MaxUint64 uint64 `env:"MAX_UINT_64"`

	MaxFloat float64 `env:"MAX_FLOAT"`
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

func (suite *EnvironmentTestSuite) TestTypeInt() {
	os.Setenv("MAX_INT", fmt.Sprint(math.MaxInt))
	os.Setenv("MAX_UINT", "18446744073709551615")
	os.Setenv("MAX_INT_64", fmt.Sprint(math.MaxInt64))
	os.Setenv("MAX_UINT_64", "18446744073709551615")

	var config Env
	err := env.Load(&config, env.Config{})
	suite.NoError(err)
	suite.Equal(math.MaxInt, config.MaxInt)
	if math.MaxUint != config.MaxUint {
		suite.Fail("config does not contain max uint")
	}
	suite.Equal(int64(math.MaxInt64), config.MaxInt64)
	if math.MaxUint64 != config.MaxUint64 {
		suite.Fail("config does not contain max uint64")
	}

	os.Unsetenv("MAX_INT")
	os.Unsetenv("MAX_UINT")
	os.Unsetenv("MAX_INT_64")
	os.Unsetenv("MAX_UINT_64")
}

func (suite *EnvironmentTestSuite) TestTypeFloat() {
	os.Setenv("MAX_FLOAT", fmt.Sprint(math.MaxFloat64))

	var config Env
	err := env.Load(&config, env.Config{})
	suite.NoError(err)

	suite.Equal(math.MaxFloat64, config.MaxFloat)

	os.Unsetenv("MAX_FLOAT")
}

func (suite *EnvironmentTestSuite) TestTypeBool() {
	os.Setenv("IS_TRUE", "true")

	var config Env
	err := env.Load(&config, env.Config{})
	suite.NoError(err)

	suite.True(config.IsTrue)

	os.Unsetenv("IS_TRUE")

	err = env.Load(&config, env.Config{})
	suite.NoError(err)

	suite.False(config.IsTrue)
}

func TestEnvironmentTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}

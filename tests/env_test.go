package tests

import (
	"fmt"
	"github.com/ricatt/go-env"
	"github.com/stretchr/testify/suite"
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

	MaxInt    int    `env:"MAX_INT"`
	MaxUint   uint   `env:"MAX_UINT"`
	MaxInt64  int64  `env:"MAX_INT_64"`
	MaxUint64 uint64 `env:"MAX_UINT_64"`

	MaxFloat float64 `env:"MAX_FLOAT"`
}

type MultiLevelEnv struct {
	Host            string `env:"HOST"`
	ExternalService struct {
		Host     string `env:"EXTERNAL_SERVICE_HOST"`
		Username string `env:"EXTERNAL_SERVICE_USERNAME"`
		Password string `env:"EXTERNAL_SERVICE_PASSWORD"`
	}
}

type EnvOptional struct {
	PackageName string `env:"PACKAGE_NAME"`
	BaseURL     string `env:"BASE_URL"`
}

type EnvInvalidType struct {
	InvalidType any `env:"INVALID_TYPE"`
}

type EnvWithDefault struct {
	HasDefault     string `env:"HAS_DEFAULT" default:"Hello, World!"`
	NotOverwritten string `env:"NOT_OVERWRITTEN" default:"Dolor sit amet"`
}

type EnvForceValue struct {
    ThisIsForced string `env:"THIS_IS_FORCED" force-value:"true"`
}

type EnvironmentTestSuite struct {
	suite.Suite
}

func (suite *EnvironmentTestSuite) TestAddValue() {
	var config Env
	err := env.Load(&config, env.Attributes{
		EnvironmentFiles: []string{".env"},
	})
	suite.NoError(err)
	suite.Empty(config.BaseURL)
	suite.Equal("env", config.PackageName)
	suite.Equal("debug", config.LogLevel)
	suite.Equal(10, config.Iterations)
}

func (suite *EnvironmentTestSuite) TestMissingValueForce() {
	var config Env
	err := env.Load(&config, env.Attributes{
		Force: true,
	})
	suite.Error(err)
	suite.Empty(config.BaseURL)
}

func (suite *EnvironmentTestSuite) TestOptionalEnvStruct() {
	var (
		config      Env
		envOptional EnvOptional
		err         error
	)
	err = env.Load(&config, env.Attributes{
		EnvironmentFiles: []string{".env"},
	})
	suite.NoError(err)
	err = env.Load(&envOptional, env.Attributes{
		EnvironmentFiles: []string{".env"},
	})
	suite.NoError(err)
}

func (suite *EnvironmentTestSuite) TestInvalidTypeEnvStruct() {
	var (
		config EnvInvalidType
		err    error
	)
	err = env.Load(&config, env.Attributes{
		EnvironmentFiles: []string{".env"},
	})
	suite.Error(err)
	suite.Equal("env: type \"interface\" not supported", err.Error())
}

func (suite *EnvironmentTestSuite) TestErrorInvalidPath() {
	var (
		config  EnvInvalidType
		err     error
		envFile = ".env-invalid-path"
	)
	err = env.Load(&config, env.Attributes{
		EnvironmentFiles:   []string{envFile},
		ErrorOnMissingFile: true,
	})
	suite.Error(err)
	suite.Equal(fmt.Sprintf("open %s: no such file or directory", envFile), err.Error())
}

func (suite *EnvironmentTestSuite) TestSuccessInvalidPath() {
	var (
		config  Env
		err     error
		envFile = ".env-invalid-path"
	)
	_ = os.Setenv("MESSAGE", "TestSuccessInvalidPath")
	err = env.Load(&config, env.Attributes{
		EnvironmentFiles: []string{envFile},
	})
	suite.NoError(err)
	suite.Equal("TestSuccessInvalidPath", config.Message)
	_ = os.Unsetenv("MESSAGE")
}

func (suite *EnvironmentTestSuite) TestNoEnvFile() {
	var (
		config Env
		err    error
	)
	err = os.Setenv("MESSAGE", "Hello World")
	suite.NoError(err)
	err = env.Load(&config, env.Attributes{})
	suite.NoError(err)
	suite.Equal("Hello World", config.Message)
	err = os.Unsetenv("MESSAGE")
	suite.NoError(err)
}

func (suite *EnvironmentTestSuite) TestWithDefault() {
	var (
		config Env
		err    error
	)

	config.BaseURL = "http://localhost:8080"

	err = env.Load(&config, env.Attributes{})
	suite.NoError(err)
	suite.Equal("http://localhost:8080", config.BaseURL)
}

func (suite *EnvironmentTestSuite) TestTypeInt() {
	os.Setenv("MAX_INT", fmt.Sprint(math.MaxInt))
	os.Setenv("MAX_UINT", "18446744073709551615")
	os.Setenv("MAX_INT_64", fmt.Sprint(math.MaxInt64))
	os.Setenv("MAX_UINT_64", "18446744073709551615")

	var config Env
	err := env.Load(&config, env.Attributes{})
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
	err := env.Load(&config, env.Attributes{})
	suite.NoError(err)

	suite.Equal(math.MaxFloat64, config.MaxFloat)

	os.Unsetenv("MAX_FLOAT")
}

func (suite *EnvironmentTestSuite) TestTypeBool() {
	var config Env
	os.Setenv("IS_TRUE", "true")
	err := env.Load(&config, env.Attributes{})
	suite.NoError(err)
	suite.True(config.IsTrue)

	os.Setenv("IS_TRUE", "false")
	suite.NoError(err)
	err = env.Load(&config, env.Attributes{})
	suite.NoError(err)
	suite.False(config.IsTrue)
}

func (suite *EnvironmentTestSuite) TestDefaultValue() {
	var config EnvWithDefault
	err := env.Load(&config, env.Attributes{
		Force: true,
	})
	suite.NoError(err)
	suite.Equal("Hello, World!", config.HasDefault)
}

func (suite *EnvironmentTestSuite) TestDefaultNotOverwitten() {
	os.Setenv("NOT_OVERWRITTEN", "Lorem ipsum")
	var config EnvWithDefault
	err := env.Load(&config, env.Attributes{
		Force: true,
	})
	suite.NoError(err)
	suite.Equal("Lorem ipsum", config.NotOverwritten)
	os.Unsetenv("NOT_OVERWRITTEN")
}

func (suite *EnvironmentTestSuite) TestMultiLevelEnv() {
	var config MultiLevelEnv
	err := env.Load(&config, env.Attributes{
		EnvironmentFiles: []string{".env"},
	})
	suite.NoError(err)
	suite.Equal("http://localhost", config.Host)
	suite.Equal("http://example.com", config.ExternalService.Host)
	suite.Equal("username", config.ExternalService.Username)
	suite.Equal("password", config.ExternalService.Password)
}

func (suite *EnvironmentTestSuite) TestForceIndivdualError() {
    var config EnvForceValue
    err := env.Load(&config, env.Attributes{
        EnvironmentFiles: []string{".env"},
    })
    suite.Error(err)

    config.ThisIsForced = "lorem ipsum"
    err = env.Load(&config, env.Attributes{})
    suite.NoError(err)
}

func TestEnvironmentTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}

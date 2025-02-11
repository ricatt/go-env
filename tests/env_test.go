package tests

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/ricatt/go-env"
	"github.com/stretchr/testify/suite"
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
	MaxInt8   int8   `env:"MAX_INT_8"`
	MaxUint8  uint8  `env:"MAX_UINT_8"`
	MaxInt32  int32  `env:"MAX_INT_32"`
	MaxUint32 uint32 `env:"MAX_UINT_32"`
	MaxInt64  int64  `env:"MAX_INT_64"`
	MaxUint64 uint64 `env:"MAX_UINT_64"`

	MaxFloat float64 `env:"MAX_FLOAT"`

	StringSlice []string `env:"STRING_SLICE"`
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
	ThisIsForced     string `env:"FORCED_VALUE" force-value:"true"`
	ThisIsAlsoForced string `env:"FORCED_VALUE" force-env:"true"`
}

type EnvironmentTestSuite struct {
	suite.Suite
}

func (suite *EnvironmentTestSuite) TestAddValue() {
	var config Env
	err := env.Load(&config, env.EnvironmentFiles(".env"))

	suite.NoError(err)
	suite.Empty(config.BaseURL)
	suite.Equal("env", config.PackageName)
	suite.Equal("debug", config.LogLevel)
	suite.Equal(10, config.Iterations)
	suite.Equal([]string{"string1", "string2", "string3"}, config.StringSlice)
}

func (suite *EnvironmentTestSuite) TestMissingValueForce() {
	var config Env
	err := env.Load(&config, env.Force(true))
	suite.Error(err)
	suite.Empty(config.BaseURL)
}

func (suite *EnvironmentTestSuite) TestOptionalEnvStruct() {
	var (
		config      Env
		envOptional EnvOptional
		err         error
	)
	err = env.Load(&config, env.EnvironmentFiles(".env"))
	suite.NoError(err)
	err = env.Load(&envOptional, env.EnvironmentFiles(".env"))
	suite.NoError(err)
}

func (suite *EnvironmentTestSuite) TestErrorInvalidPath() {
	var (
		config  EnvInvalidType
		err     error
		envFile = ".env-invalid-path"
	)
	err = env.Load(&config, env.EnvironmentFiles(envFile), env.ErrorOnMissingFile(true))
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
	err = env.Load(&config, env.EnvironmentFiles(envFile))
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
	err = env.Load(&config)
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

	err = env.Load(&config)
	suite.NoError(err)
	suite.Equal("http://localhost:8080", config.BaseURL)
}

func (suite *EnvironmentTestSuite) TestTypeInt() {
	os.Setenv("MAX_INT", fmt.Sprint(math.MaxInt))
	os.Setenv("MAX_UINT", "18446744073709551615")
	os.Setenv("MAX_INT_8", fmt.Sprint(math.MaxInt8))
	os.Setenv("MAX_UINT_8", fmt.Sprint(math.MaxUint8))
	os.Setenv("MAX_INT_32", fmt.Sprint(math.MaxInt32))
	os.Setenv("MAX_UINT_32", fmt.Sprint(math.MaxUint32))
	os.Setenv("MAX_INT_64", fmt.Sprint(math.MaxInt64))
	os.Setenv("MAX_UINT_64", "18446744073709551615")

	var config Env
	err := env.Load(&config)
	suite.NoError(err)
	suite.Equal(math.MaxInt, config.MaxInt)
	if math.MaxUint != config.MaxUint {
		suite.Fail("config does not contain max uint")
	}
	suite.Equal(math.MaxInt32, int(config.MaxInt32))
	if math.MaxUint32 != config.MaxUint32 {
		suite.Fail("config does not contain max uint32")
	}
	suite.Equal(math.MaxInt8, int(config.MaxInt8))
	if math.MaxUint8 != config.MaxUint8 {
		suite.Fail("config does not contain max uint32")
	}
	suite.Equal(int64(math.MaxInt64), config.MaxInt64)
	if math.MaxUint64 != config.MaxUint64 {
		suite.Fail("config does not contain max uint64")
	}

	os.Unsetenv("MAX_INT")
	os.Unsetenv("MAX_UINT")
	os.Unsetenv("MAX_INT_32")
	os.Unsetenv("MAX_UINT_32")
	os.Unsetenv("MAX_INT_64")
	os.Unsetenv("MAX_UINT_64")
}

func (suite *EnvironmentTestSuite) TestTypeFloat() {
	os.Setenv("MAX_FLOAT", fmt.Sprint(math.MaxFloat64))

	var config Env
	err := env.Load(&config)
	suite.NoError(err)

	suite.Equal(math.MaxFloat64, config.MaxFloat)

	os.Unsetenv("MAX_FLOAT")
}

func (suite *EnvironmentTestSuite) TestTypeBool() {
	var config Env
	os.Setenv("IS_TRUE", "true")
	err := env.Load(&config)
	suite.NoError(err)
	suite.True(config.IsTrue)

	os.Setenv("IS_TRUE", "false")
	suite.NoError(err)
	err = env.Load(&config)
	suite.NoError(err)
	suite.False(config.IsTrue)
}

func (suite *EnvironmentTestSuite) TestDefaultValue() {
	var config EnvWithDefault
	err := env.Load(&config, env.Force(true))
	suite.NoError(err)
	suite.Equal("Hello, World!", config.HasDefault)
}

func (suite *EnvironmentTestSuite) TestDefaultNotOverwitten() {
	os.Setenv("NOT_OVERWRITTEN", "Lorem ipsum")
	var config EnvWithDefault
	err := env.Load(&config, env.Force(true))
	suite.NoError(err)
	suite.Equal("Lorem ipsum", config.NotOverwritten)
	os.Unsetenv("NOT_OVERWRITTEN")
}

func (suite *EnvironmentTestSuite) TestMultiLevelEnv() {
	var config MultiLevelEnv
	err := env.Load(&config, env.EnvironmentFiles(".env"))
	suite.NoError(err)
	suite.Equal("http://localhost", config.Host)
	suite.Equal("https://example.com", config.ExternalService.Host)
	suite.Equal("username", config.ExternalService.Username)
	suite.Equal("password", config.ExternalService.Password)
}

func (suite *EnvironmentTestSuite) TestForceIndividualError() {
	err := os.Unsetenv("FORCED_VALUE")
	suite.NoError(err)

	var config EnvForceValue
	err = env.Load(&config, env.EnvironmentFiles("faulty.env"))
	suite.Error(err)

	err = env.Load(&config, env.EnvironmentFiles(".env"))
	suite.NoError(err)
}

func (suite *EnvironmentTestSuite) TestSlices() {
	var config struct {
		Int    []int    `env:"INT"`
		String []string `env:"STRING"`
		Bool   []bool   `env:"BOOL"`
	}
	os.Setenv("INT", "1,2,3")
	os.Setenv("STRING", "str1,str2,str3")
	os.Setenv("BOOL", "true,false,true")

	err := env.Load(&config)
	suite.NoError(err)

	suite.Equal([]int{1, 2, 3}, config.Int)
	suite.Equal([]string{"str1", "str2", "str3"}, config.String)
	suite.Equal([]bool{true, false, true}, config.Bool)
}

func (suite *EnvironmentTestSuite) TestSlicesFormatError() {
	var config struct {
		Int    []int    `env:"INT"`
		String []string `env:"STRING"`
		Bool   []bool   `env:"BOOL"`
	}
	os.Setenv("INT", "1,2,3,")
	os.Setenv("STRING", "str1,str2,,str3")
	os.Setenv("BOOL", ",true,false,true")

	err := env.Load(&config)
	suite.NoError(err)

	suite.Equal([]int{1, 2, 3, 0}, config.Int)
	suite.Equal([]string{"str1", "str2", "", "str3"}, config.String)
	suite.Equal([]bool{false, true, false, true}, config.Bool)
}

func (suite *EnvironmentTestSuite) TestOverwritingSystemEnv() {
	type config struct {
		PackageName         string `env:"PACKAGE_NAME"`
		ExternalServiceHost string `env:"EXTERNAL_SERVICE_HOST"`
	}

	var (
		err error
		cnf config
	)

	// Sets initial value.
	err = os.Setenv("EXTERNAL_SERVICE_HOST", "https://example.se")
	suite.NoError(err)
	err = os.Unsetenv("PACKAGE_NAME")
	suite.NoError(err)

	cnf = config{}
	err = env.Load(&cnf)
	suite.NoError(err)
	suite.Equal("", cnf.PackageName)
	suite.Equal("https://example.se", cnf.ExternalServiceHost)

	cnf = config{}
	// Using two different environment-files, overwriting both initial value and those specified first in the list.
	err = env.Load(&cnf, env.EnvironmentFiles(".env"))
	suite.NoError(err)
	suite.Equal("env", cnf.PackageName)
	suite.Equal("https://example.com", cnf.ExternalServiceHost)

	cnf = config{}
	// Using two different environment-files, overwriting both initial value and those specified first in the list.
	err = env.Load(&cnf, env.EnvironmentFiles(".env", "overwrite.env"))
	suite.NoError(err)

	suite.Equal("env", cnf.PackageName)
	suite.Equal("http://127.0.0.1", cnf.ExternalServiceHost)
}

func (suite *EnvironmentTestSuite) TestMultiLineValue() {
	var cnf struct {
		MultiLineValue string `env:"MULTILINE_VALUE"`
	}

	err := env.Load(&cnf, env.EnvironmentFiles(".env"))
	suite.NoError(err)

	expected := "Lorem ipsum\ndolor sit amet.\n"

	suite.Equal(expected, cnf.MultiLineValue)
}

func TestEnvironmentTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}

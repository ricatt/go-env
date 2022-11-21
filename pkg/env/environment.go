package env

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "env"

func Load[T any](target *T, config Config) (err error) {
	for _, path := range config.EnvironmentFiles {
		err = parseEnvFile(path)
		if err != nil {
			if e, ok := err.(*os.PathError); ok {
				if config.ErrorOnMissingFile {
					return e
				}
				err = nil
				continue
			}
			return err
		}
	}
	tp := reflect.ValueOf(target)
	proxy := reflect.ValueOf(*target)
	for i := 0; i < proxy.NumField(); i++ {
		field := proxy.Type().Field(i)
		tag := field.Tag.Get(tagName)
		value := os.Getenv(tag)
		if value == "" && !config.Force {
			if config.Force {
				return fmt.Errorf("missing value for %s", tag)
			}
			continue
		}

		fieldType := field.Type.Kind()
		switch fieldType {
		case reflect.Bool:
			v, _ := strconv.ParseBool(value)
			tp.Elem().FieldByName(field.Name).SetBool(v)
		case reflect.Int:
			v, _ := strconv.ParseInt(value, 10, 0)
			tp.Elem().FieldByName(field.Name).SetInt(v)
		case reflect.Uint:
			v, _ := strconv.ParseUint(value, 10, 0)
			tp.Elem().FieldByName(field.Name).SetUint(v)
		case reflect.Int64:
			v, _ := strconv.ParseInt(value, 10, 0)
			tp.Elem().FieldByName(field.Name).SetInt(v)
		case reflect.Uint64:
			v, _ := strconv.ParseUint(value, 10, 0)
			tp.Elem().FieldByName(field.Name).SetUint(v)
		case reflect.Float64, reflect.Float32:
			v, _ := strconv.ParseFloat(value, 0)
			tp.Elem().FieldByName(field.Name).SetFloat(v)
		case reflect.String:
			tp.Elem().FieldByName(field.Name).SetString(value)
		default:
			return fmt.Errorf("env: type \"%s\" not supported", field.Type.Kind())
		}
	}
	return
}

func parseEnvFile(path string) (err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		var val string
		if len(line) == 2 {
			val = line[1]
		}
		if err = os.Setenv(line[0], val); err != nil {
			return err
		}
	}
	return scanner.Err()
}

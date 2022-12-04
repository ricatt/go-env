package env

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
    tagDefaultValue = "default"
    tagName = "env"
)

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
    v := reflect.ValueOf(target)
    v, err = parse(v, config)
	return
}

func parse(v reflect.Value, config Config) (reflect.Value, error) {
    var err error
    el := v.Elem()
    if el.Kind() == reflect.Struct {
        numField := el.NumField()
        for i := 0; i < numField; i++ {
            field := el.Field(i)
            tag := el.Type().Field(i).Tag.Get(tagName)
            if isEqual(field) {
                continue
            }
            if field.Kind() == reflect.Struct {
                field, err = parse(field.Addr(), config)
                if err != nil {
                    return el, err
                }
                el.Field(i).Set(field)
            } else {
                value := getValue(el.Type().Field(i))
                        if value == "" {
                            if config.Force {
                                return el, fmt.Errorf("missing value for %s", tag)
                            }
                            continue
                        }
                field, err = setData(field, el.Type().Field(i), value)
                if err != nil {
                    return el, err
                }
                el.Field(i).Set(field)
            }
        }
    }
    return el, nil
}

func isEqual(field reflect.Value) bool {
    if field.Kind() == reflect.Bool {
        return false
    }
    currentValue := field.Interface()
    zeroValue := reflect.Zero(field.Type()).Interface()
    return currentValue != zeroValue
}

func setData(target reflect.Value, field reflect.StructField, value string) (reflect.Value, error) {
    switch target.Type().Kind() {
	case reflect.Bool:
		v, _ := strconv.ParseBool(value)
        target.SetBool(v)
	case reflect.Int:
		v, _ := strconv.ParseInt(value, 10, 0)
        target.SetInt(v)
	case reflect.Uint:
		v, _ := strconv.ParseUint(value, 10, 0)
        target.SetUint(v)
	case reflect.Int64:
		v, _ := strconv.ParseInt(value, 10, 0)
        target.SetInt(v)
	case reflect.Uint64:
		v, _ := strconv.ParseUint(value, 10, 0)
        target.SetUint(v)
	case reflect.Float64, reflect.Float32:
		v, _ := strconv.ParseFloat(value, 0)
        target.SetFloat(v)
	case reflect.String:
        target.SetString(value)
	default:
        return target, fmt.Errorf("env: type \"%s\" not supported", target.Kind())
	}
    return target, nil
}

func getValue(field reflect.StructField) (value string) {
    var defaultValue string
    tag := field.Tag.Get(tagName)
    value = os.Getenv(tag)
    defaultValue = field.Tag.Get(tagDefaultValue)
    if value == "" {
        value = defaultValue
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
		var val string
		line := strings.Split(scanner.Text(), "=")
        if len(line) == 0 {
            continue
        }
		if len(line) == 2 {
			val = line[1]
		}
        if len(val) == 0 {
            continue
        }
		if err = os.Setenv(line[0], val); err != nil {
			return err
		}
	}
	return scanner.Err()
}

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
	tagName         = "env"
	tagForceValue   = "force-value" //deprecated
	tagForceEnv     = "force-env"
)

// Load The primary function to load the environment into the struct.
func Load[T any](target *T, attr Attributes) (err error) {
	var values map[string]string
	for _, path := range attr.EnvironmentFiles {
		values, err = parseEnvFile(path)
		if err != nil {
			if e, ok := err.(*os.PathError); ok {
				if attr.ErrorOnMissingFile {
					return e
				}
				err = nil
				continue
			}
			return err
		}
	}
	v := reflect.ValueOf(target)
	v, err = parse(v, attr, values)
	return
}

// parse Will loop through all struct-fields, will act recursivly upon multi-level structs.
func parse(v reflect.Value, config Attributes, values map[string]string) (reflect.Value, error) {
	var err error
	el := v.Elem()
	if el.Kind() == reflect.Struct {
		numField := el.NumField()
		for i := 0; i < numField; i++ {
			field := el.Field(i)
			tag := el.Type().Field(i).Tag.Get(tagName)
			forcedEnv := isForced(el.Type().Field(i))
			if isEqual(field) {
				continue
			}
			if field.Kind() == reflect.Struct {
				field, err = parse(field.Addr(), config, values)
				if err != nil {
					return el, err
				}
				el.Field(i).Set(field)
			} else {
				value := getValue(el.Type().Field(i), values)
				if value == "" {
					if config.Force || forcedEnv {
						return el, fmt.Errorf("missing value for %s", tag)
					}
					continue
				}
				field, err = setData(field, value)
				if err != nil {
					return el, err
				}
				el.Field(i).Set(field)
			}
		}
	}
	return el, nil
}

func isForced(field reflect.StructField) bool {
	val1 := field.Tag.Get(tagForceValue)
	val2 := field.Tag.Get(tagForceEnv)
	if val1 == "true" {
		return true
	}
	if val2 == "true" {
		return true
	}
	return false
}

// isEqual A small function to make sure we aren't overwriting any entries already provided before the struct
// is added to the Load-function.
func isEqual(field reflect.Value) bool {
	if field.Kind() == reflect.Bool {
		return false
	}
	currentValue := field.Interface()
	zeroValue := reflect.Zero(field.Type()).Interface()
	return !reflect.DeepEqual(currentValue, zeroValue)
}

// setData Will cast the value into the correct type and set it for the field.
func setData(target reflect.Value, value string) (reflect.Value, error) {
	switch target.Type().Kind() {
	case reflect.Slice:
		values := strings.Split(value, ",")
		s := reflect.New(target.Type())
		for _, v := range values {
			newField := reflect.New(target.Type().Elem())
			converted, _ := setData(newField.Elem(), v)
			s.Elem().Set(reflect.Append(s.Elem(), converted))
		}
		target.Set(s.Elem())
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

// getValue fetches the value from the environment and fetches potential default values from struct field.
// Default value will only be set if value is empty.
func getValue(field reflect.StructField, values map[string]string) (value string) {
	tag := field.Tag.Get(tagName)
	value = os.Getenv(tag)
	if value == "" {
		value = values[tag]
	}
	defaultValue := field.Tag.Get(tagDefaultValue)
	if value == "" {
		value = defaultValue
	}
	return
}

// parseEnvFile Will fetch all entries from provded file and set it in the environment.
func parseEnvFile(path string) (map[string]string, error) {
	var (
		err    error
		file   *os.File
		values = make(map[string]string)
	)
	file, err = os.Open(path)
	if err != nil {
		return values, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var val string
		line := strings.Split(scanner.Text(), "=")
		if len(line) == 0 {
			continue
		}
		val = strings.Join(line[1:], "=")
		if len(val) == 0 {
			continue
		}
		values[line[0]] = val
	}
	return values, scanner.Err()
}

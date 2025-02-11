package env

import (
	"bufio"
	"errors"
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
func Load[T any](target *T, opts ...Attribute) (err error) {
	var attr attributes
	for _, f := range opts {
		f(&attr)
	}

	values := make(map[string]string)
	for _, path := range attr.EnvironmentFiles {
		values, err = parseEnvFile(path, values)
		if err != nil {
			var e *os.PathError
			if errors.As(err, &e) {
				if attr.ErrorOnMissingFile {
					return e
				}
				continue
			}
			return err
		}
	}
	v := reflect.ValueOf(target)
	_, err = parse(v, attr, values)
	return
}

// parse Will loop through all struct-fields, will act recursivly upon multi-level structs.
func parse(v reflect.Value, config attributes, values map[string]string) (reflect.Value, error) {
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
	if field.Tag.Get(tagForceValue) == "true" || field.Tag.Get(tagForceEnv) == "true" {
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
	case reflect.Int8:
		v, _ := strconv.ParseInt(value, 10, 8)
		target.SetInt(v)
	case reflect.Uint8:
		v, _ := strconv.ParseUint(value, 10, 8)
		target.SetUint(v)
	case reflect.Int32:
		v, _ := strconv.ParseInt(value, 10, 32)
		target.SetInt(v)
	case reflect.Uint32:
		v, _ := strconv.ParseUint(value, 10, 32)
		target.SetUint(v)
	case reflect.Int64:
		v, _ := strconv.ParseInt(value, 10, 0)
		target.SetInt(v)
	case reflect.Uint64:
		v, _ := strconv.ParseUint(value, 10, 0)
		target.SetUint(v)
	case reflect.Float32:
		v, _ := strconv.ParseFloat(value, 32)
		target.SetFloat(v)
	case reflect.Float64:
		v, _ := strconv.ParseFloat(value, 64)
		target.SetFloat(v)
	case reflect.String:
		target.SetString(value)
	case reflect.Interface:
		target.Set(reflect.ValueOf(value))

	default:
		return target, fmt.Errorf("env: type \"%s\" not supported", target.Kind())
	}
	return target, nil
}

// getValue fetches the value from the environment and fetches potential default values from struct field.
// Default value will only be set if value is empty.
func getValue(field reflect.StructField, values map[string]string) string {
	tag := field.Tag.Get(tagName)
	value, ok := values[tag]
	if value == "" || !ok {
		value = os.Getenv(tag)
	}
	defaultValue := field.Tag.Get(tagDefaultValue)
	if value == "" {
		value = defaultValue
	}
	return value
}

// parseEnvFile Will fetch all entries from provded file and set it in the environment.
func parseEnvFile(filename string, envMap map[string]string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var key, value string
	inMultiLine := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Detect key-value pair
		if !inMultiLine {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue // Skip malformed lines
			}
			key = strings.TrimSpace(parts[0])
			value = strings.TrimSpace(parts[1])

			// Detect multi-line start (values enclosed in double quotes)
			if strings.HasPrefix(value, "\"") && !strings.HasSuffix(value, "\"") {
				inMultiLine = true
				value = strings.TrimPrefix(value, "\"") // Remove opening quote
				continue
			}

			// Store if it's a normal single-line key-value pair
			envMap[key] = strings.Trim(value, "\"")
		} else {
			// Accumulate multi-line values
			value += "\n" + line

			// Detect end of multi-line value
			if strings.HasSuffix(line, "\"") {
				inMultiLine = false
				value = strings.TrimSuffix(value, "\"") // Remove closing quote
				envMap[key] = value
			}
		}
	}
	return envMap, scanner.Err()
}


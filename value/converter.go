package value

import (
	"strconv"
	"strings"
	"time"

	"github.com/jit-b/go/errors"
)

const (
	errConvertIntegerTemplate      = "Can't convert value '%s' to integer"
	errConvertTimeDurationTemplate = "Can't convert value '%s' as time duration"
	errConvertStringSliceTemplate  = "Can't convert value '%s' with separator '%s' to slice of strings"
	errConvertStringMapTemplate    = "Can't convert value '%s' with separator '%s' and pair separator '%s' as map of strings"
)

type Converter string

var boolMap = map[string]bool{
	"true":  true,
	"1":     true,
	"false": false,
	"0":     false,
}

func NewConverter(v string) *Converter {
	result := Converter(v)

	return &result
}

func (c *Converter) AsDuration() (time.Duration, error) {
	value := c.AsString()
	result, err := time.ParseDuration(value)
	if err != nil {
		return 0, errors.New(errConvertTimeDurationTemplate, value).Cause(err)
	}

	return result, nil
}

func (c *Converter) AsStringMapOfStrings(pairSeparator string, separator string) (map[string]string, error) {
	data := c.AsString()
	result := map[string]string{}
	for _, keyValues := range strings.Split(data, separator) {
		pair := strings.Split(keyValues, pairSeparator)
		if len(pair) == 2 && pair[0] != "" && pair[1] != "" {
			result[pair[0]] = pair[1]
		}
	}

	if len(result) == 0 {
		return nil, errors.New(errConvertStringMapTemplate, data, separator, pairSeparator)
	}

	return result, nil
}

func (c *Converter) AsSliceOfString(separator string) ([]string, error) {
	values := c.AsString()
	var result []string
	for _, value := range strings.Split(values, ",") {
		result = append(result, strings.Trim(value, "\n\t\r "))
	}

	if len(result) == 1 && result[0] == values {
		return nil, errors.New(errConvertStringSliceTemplate, values, separator)
	}

	return result, nil
}

func (c *Converter) AsIntegerWithDefaultValue(defaultValue int64) int64 {
	value, err := c.AsInteger()
	if err != nil {
		return defaultValue
	}

	return value
}

func (c *Converter) AsInteger() (int64, error) {
	value := c.AsString()
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New(errConvertIntegerTemplate, value).Cause(err)
	}

	return int64(result), nil
}

func (c *Converter) AsBoolWithDefault(defaultValue bool) bool {
	value, isExist := boolMap[c.AsString()]
	if !isExist {
		return defaultValue
	}

	return value
}

func (c *Converter) AsString() string {
	return string(*c)
}

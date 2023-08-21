package environment

import (
	"testing"
	"time"

	env "github.com/jit-b/go/environment"
	"github.com/stretchr/testify/assert"
)

func Test_Parser_Get_AsRFC3339Time(test *testing.T) {
	test.Run(
		"Success getting value from env as RFC3339 time format",
		func(test *testing.T) {
			testTime, err := time.Parse(time.RFC3339, "2023-08-21T21:04:05+03:00")
			assert.NoError(test, err)

			test.Setenv("TEST_KEY", testTime.Format(time.RFC3339))

			result, err := env.Get("TEST_KEY").AsRFC3339Time()
			assert.NoError(test, err)
			assert.Equal(test, testTime, result)
		},
	)

	test.Run(
		"Success getting value with templating from env as duration",
		func(test *testing.T) {
			firstTestTime, err := time.Parse(time.RFC3339, "2023-08-21T21:04:05+03:00")
			assert.NoError(test, err)
			secondTestTime, err := time.Parse(time.RFC3339, "2023-08-09T12:30:00+03:00")
			assert.NoError(test, err)

			test.Setenv("FIRST_TEST_KEY", "2023-08-21T21:04:05+03:00")
			test.Setenv("SECOND_TEST_KEY", "1988-13-11T13:30:00+10:00")

			test.Setenv("FIRST_TEST_KEY", firstTestTime.Format(time.RFC3339))
			test.Setenv("SECOND_TEST_KEY", secondTestTime.Format(time.RFC3339))

			result, err := env.Get("%s_TEST_KEY", "FIRST").AsRFC3339Time()
			assert.NoError(test, err)
			assert.Equal(test, firstTestTime, result)

			result, err = env.Get("%s_TEST_KEY", "SECOND").AsRFC3339Time()
			assert.NoError(test, err)
			assert.Equal(test, secondTestTime, result)
		},
	)

	test.Run(
		"Error when value is invalid",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "2")

			_, err := env.Get("TEST_KEY").AsRFC3339Time()
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when env variable is empty",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "")

			_, err := env.Get("TEST_KEY").AsRFC3339Time()
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when env variable not exist",
		func(test *testing.T) {

			_, err := env.Get("TEST_KEY").AsDuration()
			assert.Error(test, err)
		},
	)
}

func Test_Parser_Get_AsDuration(test *testing.T) {
	test.Run(
		"Success getting value from env as duration",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "2h")

			result, err := env.Get("TEST_KEY").AsDuration()
			assert.NoError(test, err)
			assert.Equal(test, 2*time.Hour, result)
		},
	)

	test.Run(
		"Success getting value with templating from env as duration",
		func(test *testing.T) {
			test.Setenv("FIRST_TEST_KEY", "2h")
			test.Setenv("SECOND_TEST_KEY", "15s")

			result, err := env.Get("%s_TEST_KEY", "FIRST").AsDuration()
			assert.NoError(test, err)
			assert.Equal(test, 2*time.Hour, result)

			result, err = env.Get("%s_TEST_KEY", "SECOND").AsDuration()
			assert.NoError(test, err)
			assert.Equal(test, 15*time.Second, result)
		},
	)

	test.Run(
		"Error when value is invalid",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "2")

			_, err := env.Get("TEST_KEY").AsDuration()
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when env variable is empty",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "")

			_, err := env.Get("TEST_KEY").AsDuration()
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when env variable not exist",
		func(test *testing.T) {

			_, err := env.Get("TEST_KEY").AsDuration()
			assert.Error(test, err)
		},
	)
}

func Test_Parser_Get_AsStringMapOfStrings(test *testing.T) {
	test.Run(
		"Success getting value from env as string map of strings",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "firstSubKey=firstSubValue,secondSubKey=secondSubValue")

			result, err := env.Get("TEST_KEY").AsStringMapOfStrings("=", ",")
			assert.NoError(test, err)
			assert.Equal(
				test,
				"firstSubValue",
				result["firstSubKey"],
			)
			assert.Equal(
				test,
				"secondSubValue",
				result["secondSubKey"],
			)
		},
	)

	test.Run(
		"Error when can't parse value with given separators",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "firstSubKey->firstSubValue|secondSubKey->secondSubValue")

			_, err := env.Get("TEST_KEY").AsStringMapOfStrings("=", ",")
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when env variable is empty",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "")

			_, err := env.Get("TEST_KEY").AsStringMapOfStrings("=", ",")
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when env variable not exist",
		func(test *testing.T) {

			_, err := env.Get("TEST_KEY").AsStringMapOfStrings("=", ",")
			assert.Error(test, err)
		},
	)
}

func Test_Parser_Get_AsSliceOfString(test *testing.T) {
	test.Run(
		"Success getting value from env as slice of string",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "firstValue,secondValue")

			result, err := env.Get("TEST_KEY").AsSliceOfString(",")
			assert.NoError(test, err)
			assert.Equal(
				test,
				"firstValue",
				result[0],
			)
			assert.Equal(
				test,
				"secondValue",
				result[1],
			)
		},
	)

	test.Run(
		"Error when can't parse value with given separator",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "firstValue;secondValue")

			_, err := env.Get("TEST_KEY").AsSliceOfString(",")
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when can't parse value with cause it's empty",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "")

			_, err := env.Get("TEST_KEY").AsSliceOfString(",")
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when can't parse value with cause it's doesn't exist",
		func(test *testing.T) {
			_, err := env.Get("TEST_KEY").AsSliceOfString(",")
			assert.Error(test, err)
		},
	)
}

func Test_Parser_Get_AsIntegerWithDefaultValue(test *testing.T) {
	test.Run(
		"Success getting value from env",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "13")

			result := env.Get("TEST_KEY").AsIntegerWithDefaultValue(7)
			assert.Equal(test, int64(13), result)
		},
	)

	test.Run(
		"Success getting default value when env is doesn't exist",
		func(test *testing.T) {
			result := env.Get("TEST_KEY").AsIntegerWithDefaultValue(7)
			assert.Equal(test, int64(7), result)
		},
	)

	test.Run(
		"Success getting default value when env is empty",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "")

			result := env.Get("TEST_KEY").AsIntegerWithDefaultValue(7)
			assert.Equal(test, int64(7), result)
		},
	)
}

func Test_Parser_Get_AsInteger(test *testing.T) {
	test.Run(
		"Success getting value from env",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "13")

			result, err := env.Get("TEST_KEY").AsInteger()
			assert.NoError(test, err)
			assert.Equal(test, int64(13), result)
		},
	)

	test.Run(
		"Error when value is invalid",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "A")

			_, err := env.Get("TEST_KEY").AsInteger()
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when env variable is empty",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "")

			_, err := env.Get("TEST_KEY").AsInteger()
			assert.Error(test, err)
		},
	)

	test.Run(
		"Error when env variable not exist",
		func(test *testing.T) {

			_, err := env.Get("TEST_KEY").AsInteger()
			assert.Error(test, err)
		},
	)
}

func Test_Parser_Get_AsBoolWithDefault(test *testing.T) {
	test.Run(
		"Success getting true from env",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "false")

			result := env.Get("TEST_KEY").AsBoolWithDefault(true)
			assert.Equal(test, false, result)
		},
	)

	test.Run(
		"Success getting false from env",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "true")

			result := env.Get("TEST_KEY").AsBoolWithDefault(false)
			assert.Equal(test, true, result)
		},
	)

	test.Run(
		"Success getting 1 as true from env",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "1")

			result := env.Get("TEST_KEY").AsBoolWithDefault(false)
			assert.Equal(test, true, result)
		},
	)

	test.Run(
		"Success getting 0 as false from env",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "0")

			result := env.Get("TEST_KEY").AsBoolWithDefault(true)
			assert.Equal(test, false, result)
		},
	)

	test.Run(
		"Success getting default value when env is doesn't exist",
		func(test *testing.T) {
			result := env.Get("TEST_KEY").AsBoolWithDefault(true)
			assert.Equal(test, true, result)
		},
	)

	test.Run(
		"Success getting default value when env is empty",
		func(test *testing.T) {
			test.Setenv("TEST_KEY", "")

			result := env.Get("TEST_KEY").AsBoolWithDefault(true)
			assert.Equal(test, true, result)
		},
	)
}

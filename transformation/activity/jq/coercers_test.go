package jq_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/pmhofman-se/flogo-jq-activity/transformation/activity/jq/activity/jq"
	"github.com/stretchr/testify/assert"
)

func Test_StringCoercer_Coerce(t *testing.T) {
	// Special input values for the tests
	// empty_string_input := ""          // error expected for empty string, but it will also return int64(0)
	make_chan_input := make(chan int) // this will cause an error
	// Expected results for the tests
	no_error := "" // no error expected for int/int32, int64, float32, float64, nil, true/false, and valid number strings
	coerce_error_format := "unable to coerce (%T)(%v) to string"
	coerce_error := func(val interface{}) string {
		// Format the syntax error message for the illegal string input
		return fmt.Sprintf(coerce_error_format, val, val)
	}

	tests := []struct {
		name     string // description of this test case
		input    interface{}
		expected interface{}
		wantErr  bool
		errorMsg string
	}{
		{"String", "Hello, World!", "Hello, World!", false, no_error},
		{"Integer 123", 123, "123", false, no_error},
		{"Float 45.67", 45.67, "45.67", false, no_error},
		{"Boolean true", true, "true", false, no_error},
		{"Nil", nil, "", false, no_error},
		{"Slice []int{1,2,3}", []int{1, 2, 3}, "[1,2,3]", false, no_error},
		{"Map map[string]int{\"a\":1,\"b\":2}", map[string]int{"a": 1, "b": 2}, `{"a":1,"b":2}`, false, no_error},
		{"Struct struct{X int}{X:10}", "{10}", "{10}", false, no_error},
		{"Unsupported type chan int", make_chan_input, "", true, coerce_error(make_chan_input)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc = jq.NewStringCoercer(nil)
			got, err := sc.Coerce(tt.input)
			if tt.wantErr {
				assert.NotNil(t, err, "Expected error for input: %v", tt.input)
				assert.EqualErrorf(t, err, tt.errorMsg, "Error should be: %v, got: %v", tt.errorMsg, err)
				t.Logf("Error: %v", err)
			} else {
				assert.Nil(t, err, "Unexpected error for input: %v", tt.input)
			}
			// Check if the type of got matches the expected type
			if reflect.TypeOf(got) != reflect.TypeOf(tt.expected) {
				t.Errorf("Coerce() type mismatch: expected %T, got %T", tt.expected, got)
			}
			// Check if the value of got matches the expected value
			// Use DeepEqual to handle cases where the types might be different but values are equivalent
			assert.Equal(t, tt.expected, got, "Expected output: %v, got: %v", tt.expected, got)

			t.Logf("Input: %v (%v), Output: %v (%v), Expected: %v (%v)", tt.input, reflect.TypeOf(tt.input), got, reflect.TypeOf(got), tt.expected, reflect.TypeOf(tt.expected))
		})
	}
}

func Test_NumberCoercer(t *testing.T) {
	// Special input values for the tests
	f32_input := float32(2.718)       // special case for float32 to ensure it gets converted to float64 in expected result
	empty_string_input := ""          // error expected for empty string, but it will also return int64(0)
	illegal_string_input := "a"       // this will cause an error
	make_chan_input := make(chan int) // this will cause an error
	// Expected results for the tests
	no_error := "" // no error expected for int/int32, int64, float32, float64, nil, true/false, and valid number strings
	parse_error_format := "parse int error: strconv.ParseInt: parsing \"%v\": invalid syntax / parse float error: strconv.ParseFloat: parsing \"%v\": invalid syntax"
	parse_error := func(val string) string {
		// Format the syntax error message for the illegal string input
		return fmt.Sprintf(parse_error_format, val, val)
	}
	type_value_format := "(%T)(%v)"
	type_value_formated := func(val interface{}) string {
		return fmt.Sprintf(type_value_format, val, val)
	}
	coerce_error_format := "parse int error: unable to coerce %v to integer / parse float error: unable to coerce %v to float64"
	coerce_error := func(val interface{}) string {
		// Format the syntax error message for the illegal string input
		return fmt.Sprintf(coerce_error_format, val, val)
	}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
		wantErr  bool
		errorMsg string
	}{
		{"Int", 42, int64(42), false, no_error}, // always start with a carefully chosen number
		{"Int64", int64(1234567890), int64(1234567890), false, no_error},
		{"Int64", int64(math.MaxInt64), int64(math.MaxInt64), false, no_error},
		{"Int64", int64(math.MinInt64), int64(math.MinInt64), false, no_error},
		{"Float64", float64(3.1415), float64(3.1415), false, no_error},
		{"Float64", float64(3.1415), float64(3.1415), false, no_error},
		{"Float32", f32_input, float64(f32_input), false, no_error},
		{"Nil", nil, int64(0), false, no_error},
		{"Bool", true, int64(1), false, no_error},
		{"EmptyString", empty_string_input, int64(0), true, parse_error(empty_string_input)},
		{"StringInt", "1", int64(1), false, no_error},
		{"StringLargeInt", "1234567890123456789012345678901234567890", float64(1.2345678901234568e+39), false, no_error},
		{"StringFloat", "1.1", float64(1.1), false, no_error},
		{"StringInvalid", illegal_string_input, int64(0), true, parse_error(illegal_string_input)},
		{"Boolean true", true, int64(1), false, no_error},
		{"Nil", nil, int64(0), false, no_error},
		{"Slice []int{1,2,3}", []int{1, 2, 3}, int64(0), true, coerce_error("[]int{1, 2, 3}")},
		{"Map map[string]int{\"a\":1,\"b\":2}", map[string]int{"a": 1, "b": 2}, int64(0), true, coerce_error("map[string]int{\"a\":1, \"b\":2}")},
		{"Struct struct{X int}{X:10}", "{10}", int64(0), true, parse_error("{10}")},
		{"Unsupported type chan int", make_chan_input, int64(0), true, coerce_error(type_value_formated(make_chan_input))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := jq.NewNumberCoercer(nil)
			got, err := nc.Coerce(tt.input)
			if tt.wantErr {
				assert.NotNil(t, err, "Expected error for input: %v", tt.input)
				assert.EqualErrorf(t, err, tt.errorMsg, "Error should be: %v, got: %v", tt.errorMsg, err)
				t.Logf("Error: %v", err)
			} else {
				assert.Nil(t, err, "Unexpected error for input: %v", tt.input)
			}
			// Check if the type of got matches the expected type
			if reflect.TypeOf(got) != reflect.TypeOf(tt.expected) {
				t.Errorf("ParseNumber() type mismatch: expected %T, got %T", tt.expected, got)
			}
			// Check if the value of got matches the expected value
			// Use DeepEqual to handle cases where the types might be different but values are equivalent
			assert.Equal(t, tt.expected, got, "Expected output: %v, got: %v", tt.expected, got)

			t.Logf("Input: %v (%v), Output: %v (%v), Expected: %v (%v)", tt.input, reflect.TypeOf(tt.input), got, reflect.TypeOf(got), tt.expected, reflect.TypeOf(tt.expected))
		})
	}
}

func Test_NewStringCoercer(t *testing.T) {
	tests := []struct {
		name   string
		logger interface {
			Debugf(format string, args ...interface{})
		}
	}{
		{"With nil logger", nil},
		{"With actual logger", &minimalTestLogger{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coercer := jq.NewStringCoercer(tt.logger)

			// Basic assertions - just test that constructor works
			assert.NotNil(t, coercer)
			assert.IsType(t, &jq.StringCoercer{}, coercer)
		})
	}
}

func Test_NewNumberCoercer(t *testing.T) {
	tests := []struct {
		name   string
		logger interface {
			Debugf(format string, args ...interface{})
		}
	}{
		{"With nil logger", nil},
		{"With actual logger", &minimalTestLogger{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coercer := jq.NewNumberCoercer(tt.logger)

			// Basic assertions - just test that constructor works
			assert.NotNil(t, coercer)
			assert.IsType(t, &jq.NumberCoercer{}, coercer)
		})
	}
}

// Minimal test logger
type minimalTestLogger struct{}

func (tl *minimalTestLogger) Debugf(format string, args ...interface{}) {
	// Do nothing - just satisfy the interface
}

// Spying test logger to capture debug messages
type spyingTestLogger struct {
	messages []string
}

func (tl *spyingTestLogger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	tl.messages = append(tl.messages, msg)
}

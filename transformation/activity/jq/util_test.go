package jq

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUtil_ParseNumber tests the ParseNumber function with a variety of input types and values.
// It covers the following cases:
//   - Integer types (int, int64, math.MaxInt64, math.MinInt64)
//   - Floating point types (float32, float64)
//   - nil input
//   - Boolean input (true/false)
//   - String inputs representing integers, large integers, floats, empty strings, and invalid strings
//
// The test verifies:
//   - Correct conversion of input values to int64 or float64 as appropriate
//   - Proper error handling for invalid or empty string inputs
//   - Type correctness of the returned value
//   - Accurate error messages for parsing failures
//
// The test uses table-driven testing to iterate over various scenarios and asserts both the output value and error state.
func TestUtil_ParseNumber(t *testing.T) {
	// Special input values for the tests
	f32_input := float32(2.718) // special case for float32 to ensure it gets converted to float64 in expected result
	empty_string_input := ""    // error expected for empty string, but it will also return int64(0)
	illegal_string_input := "a" // this will cause an error
	// Expected results for the tests
	no_error := "" // no error expected for int/int32, int64, float32, float64, nil, true/false, and valid number strings
	syntax_error_format := "parse int error: strconv.ParseInt: parsing \"%v\": invalid syntax / parse float error: strconv.ParseFloat: parsing \"%v\": invalid syntax"
	syntax_error := func(val string) string {
		// Format the syntax error message for the illegal string input
		return fmt.Sprintf(syntax_error_format, val, val)
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
		{"EmptyString", empty_string_input, int64(0), true, syntax_error(empty_string_input)},
		{"StringInt", "1", int64(1), false, no_error},
		{"StringLargeInt", "1234567890123456789012345678901234567890", float64(1.2345678901234568e+39), false, no_error},
		{"StringFloat", "1.1", float64(1.1), false, no_error},
		{"StringInvalid", illegal_string_input, int64(0), true, syntax_error(illegal_string_input)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNumber(tt.input)
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

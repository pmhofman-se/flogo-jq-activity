package jq

import (
	"fmt"

	"github.com/project-flogo/core/data/coerce"
)

// ParseNumber converts various input types to a number (int64 or float64).
// It handles int, int64, float64, float32, nil, bool, and string representations of numbers.
// It returns an error if the input cannot be converted to a number.
func ParseNumber(inval interface{}) (outval interface{}, err error) {
	// make sure to support both int64 and float64 values
	// some special values, like "", nil, true/false are also handled by the coerce package
	// "" results into int64(0), but it will still give an error, which can be safely ignored
	// nil results in int64(0), but it will not give an error
	// true/false results in int64(1) / int64(0), but it will not give an error
	intVal, intErr := coerce.ToInt64(inval)
	fltVal, fltErr := coerce.ToFloat64(inval)
	if intErr != nil && fltErr != nil {
		err := fmt.Errorf("parse int error: %v / parse float error: %v", intErr.Error(), fltErr.Error())
		return intVal, err
	} else if intErr != nil && fltErr == nil {
		return fltVal, nil
	} else if intErr == nil && fltErr != nil {
		return intVal, nil
	} else {
		if float64(intVal) != fltVal {
			outval = fltVal
		} else {
			outval = intVal
		}
		return outval, err
	}
}

package jq

import (
	"fmt"

	"github.com/project-flogo/core/data/coerce"
)

// StringCoercer with proper error handling
type StringCoercer struct {
	logger interface {
		Debugf(format string, args ...interface{})
	}
}

func NewStringCoercer(logger interface {
	Debugf(format string, args ...interface{})
}) *StringCoercer {
	return &StringCoercer{logger: logger}
}

func (sc *StringCoercer) Coerce(value interface{}) (interface{}, error) {
	// coerce.ToString returns (string, error)
	return sc.parseString(value)
}

func (sc *StringCoercer) parseString(inval interface{}) (interface{}, error) {
	stringVal, stringErr := coerce.ToString(inval)

	if sc.logger != nil {
		sc.logger.Debugf("ParseString input: %v (%T)", inval, inval)
		sc.logger.Debugf("  coerce.ToString result: %v, error: %v", stringVal, stringErr)
	}

	return stringVal, stringErr
}

// NumberCoercer with proper error handling
type NumberCoercer struct {
	logger interface {
		Debugf(format string, args ...interface{})
	}
}

func NewNumberCoercer(logger interface {
	Debugf(format string, args ...interface{})
}) *NumberCoercer {
	return &NumberCoercer{logger: logger}
}

func (nc *NumberCoercer) Coerce(value interface{}) (interface{}, error) {
	return nc.parseNumber(value)
}

func (nc *NumberCoercer) parseNumber(inval interface{}) (interface{}, error) {
	// coerce.ToInt64 and coerce.ToFloat64 typically don't return errors in the same way,
	// but let's be consistent
	intVal, intErr := coerce.ToInt64(inval)
	fltVal, fltErr := coerce.ToFloat64(inval)

	if nc.logger != nil {
		nc.logger.Debugf("ParseNumber input: %v (%T)", inval, inval)
		nc.logger.Debugf("  coerce.ToInt64 result: %v, error: %v", intVal, intErr)
		nc.logger.Debugf("  coerce.ToFloat64 result: %v, error: %v", fltVal, fltErr)
	}

	// Original logic preserved
	if intErr != nil && fltErr != nil {
		err := fmt.Errorf("parse int error: %v / parse float error: %v", intErr, fltErr)
		return intVal, err
	} else if intErr != nil && fltErr == nil {
		return fltVal, nil
	} else if intErr == nil && fltErr != nil {
		return intVal, nil
	} else {
		if float64(intVal) != fltVal {
			return fltVal, nil
		} else {
			return intVal, nil
		}
	}
}

// ParseNumber maintains backward compatibility
func ParseNumber(inval interface{}) (outval interface{}, err error) {
	coercer := &NumberCoercer{} // logger is nil, so no logging
	return coercer.Coerce(inval)
}

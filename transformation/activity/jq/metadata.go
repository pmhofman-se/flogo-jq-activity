package jq

import (
	"github.com/project-flogo/core/data/coerce"
)

type JQActivity struct {
	settings *Settings
}

type Settings struct {
	Script string `md:"Script"`
}

type Input struct {
	InputData     map[string]interface{} `md:"InputData"`
	ArgumentNames []interface{}          `md:"ArgumentNames"`
	Arguments     map[string]interface{} `md:"Arguments"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"InputData":     i.InputData,
		"ArgumentNames": i.ArgumentNames,
		"Arguments":     i.Arguments,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error

	i.InputData, err = coerce.ToObject(values["InputData"])
	if err != nil {
		return err
	}

	i.ArgumentNames, err = coerce.ToArray(values["ArgumentNames"])
	if err != nil {
		return err
	}

	i.Arguments, err = coerce.ToObject(values["Arguments"])
	if err != nil {
		return err
	}

	return nil
}

type Output struct {
	Error        bool          `md:"Error"`
	ErrorMessage string        `md:"ErrorMessage"`
	OutputData   []interface{} `md:"OutputData"`
}

// ToMap conversion
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Error":        o.Error,
		"ErrorMessage": o.ErrorMessage,
		"OutputData":   o.OutputData,
	}
}

// FromMap conversion
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error

	o.Error, err = coerce.ToBool(values["Error"])
	if err != nil {
		return err
	}

	o.ErrorMessage, err = coerce.ToString(values["ErrorMessage"])
	if err != nil {
		return err
	}

	o.OutputData, err = coerce.ToArray(values["OutputData"])
	if err != nil {
		return err
	}
	return nil
}

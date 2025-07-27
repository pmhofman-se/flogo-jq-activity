package jq

import (
	"testing"

	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestMetaData(t *testing.T) {
	/*
		Set up activity Settings
	*/
	settingsScript := "."
	settingsJQ := &Settings{
		Script: settingsScript,
	}

	initContextJQ := test.NewActivityInitContext(settingsJQ, nil)
	actJQ, err := New(initContextJQ)
	assert.Nil(t, err)

	tcJQ := test.NewActivityContext(actJQ.Metadata())

	/*
		some tests only to increase coverage of metadata.go
		some are quite useless, and even dubious
	*/

	input := &Input{}
	err = tcJQ.GetInputObject(input)
	assert.NotNil(t, input)
	assert.Nil(t, err)

	v := input.ToMap()
	assert.NotNil(t, v)

	v["InputData"] = 1
	err = input.FromMap(v)
	assert.NotNil(t, err)

	v["Arguments"] = 2
	err = input.FromMap(v)
	assert.NotNil(t, err)

	output := &Output{}
	err = tcJQ.GetOutputObject(output)
	assert.NotNil(t, output)
	assert.Nil(t, err)

	v = output.ToMap()
	assert.NotNil(t, v)

	v["Error"] = 3
	err = output.FromMap(v)
	assert.Nil(t, err) // coerce to bool will just work

	v = output.ToMap()
	v["ErrorMessage"] = 4
	err = output.FromMap(v)
	assert.Nil(t, err) // coerce to string will just work

	v = output.ToMap()
	v["OutputData"] = 5
	err = output.FromMap(v)
	assert.Nil(t, err) // coerce to array will just work

	v = output.ToMap()
	assert.NotNil(t, v)
}

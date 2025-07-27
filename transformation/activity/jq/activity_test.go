package jq

import (
	"reflect"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"

	"github.com/stretchr/testify/assert"
)

const (
	KEYS_BASE_PATH = "../../../testdata/"
)

func CreateEmptyArgumentNames() []map[string]interface{} {
	argumentNames := make([]map[string]interface{}, 0)
	return argumentNames
}

func CreateEmptyArguments() map[string]interface{} {
	arguments := make(map[string]interface{})
	return arguments
}

func CreateArgumentNames() []map[string]interface{} {
	argNames := make([]map[string]interface{}, 5)
	argNames[0] = map[string]interface{}{"Name": "arg1", "Type": "String"}
	argNames[1] = map[string]interface{}{"Name": "arg2", "Type": "String"}
	argNames[2] = map[string]interface{}{"Name": "arg3", "Type": "String"}
	argNames[3] = map[string]interface{}{"Name": "arg4", "Type": "Number"}
	argNames[4] = map[string]interface{}{"Name": "arg5", "Type": "Number"}
	return argNames
}

func CreateArgumentNamesWithUnsupportedType() []map[string]interface{} {
	ArgNames := make([]map[string]interface{}, 1)
	ArgNames[0] = map[string]interface{}{"Name": "dummy", "Type": "Dummy"}
	return ArgNames
}

func CreateArguments() map[string]interface{} {
	args := make(map[string]interface{})
	args["arg1"] = "baz"
	args["arg2"] = "value2"
	args["arg3"] = "value3"
	args["arg4"] = 4
	args["arg5"] = 5
	return args
}

func CreatePayloadWithUnsupportedType() map[string]interface{} {
	args := make(map[string]interface{})
	args["dummy"] = "nonsense"
	return args
}

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&JQActivity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestJQ_NoScript_NoArguments_FooBar(t *testing.T) {

	/*
		Set up activity Settings
	*/
	settingsJQ := &Settings{
		Script: "",
	}

	initContextJQ := test.NewActivityInitContext(settingsJQ, nil)
	actJQ, err := New(initContextJQ)
	assert.Nil(t, err)

	tcJQ := test.NewActivityContext(actJQ.Metadata())

	argumentNames := CreateEmptyArgumentNames()
	tcJQ.SetInput("ArgumentNames", argumentNames)

	arguments := CreateEmptyArguments()
	tcJQ.SetInput("Arguments", arguments)

	tcJQ.SetInput("InputData", map[string]interface{}{"foo": "bar"})

	done, err := actJQ.Eval(tcJQ)
	assert.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, 1, len(tcJQ.GetOutput("OutputData").([]interface{})))
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, tcJQ.GetOutput("OutputData").([]interface{})[0])
	assert.Equal(t, false, tcJQ.GetOutput("Error").(bool))
	assert.Equal(t, "", tcJQ.GetOutput("ErrorMessage").(string))
}

func TestJQ_ScriptWithoutArguments_Arguments_FooBar(t *testing.T) {

	/*
		Set up activity Settings
	*/
	settingsJQ := &Settings{
		Script: ".",
	}

	initContextJQ := test.NewActivityInitContext(settingsJQ, nil)
	actJQ, err := New(initContextJQ)
	assert.Nil(t, err)

	tcJQ := test.NewActivityContext(actJQ.Metadata())

	argumentNames := CreateArgumentNames()
	tcJQ.SetInput("ArgumentNames", argumentNames)

	arguments := CreateArguments()
	tcJQ.SetInput("Arguments", arguments)

	tcJQ.SetInput("InputData", map[string]interface{}{"foo": "bar"})

	done, err := actJQ.Eval(tcJQ)
	assert.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, 1, len(tcJQ.GetOutput("OutputData").([]interface{})))
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, tcJQ.GetOutput("OutputData").([]interface{})[0])
	assert.Equal(t, false, tcJQ.GetOutput("Error").(bool))
	assert.Equal(t, "", tcJQ.GetOutput("ErrorMessage").(string))
}

func TestJQ_ScriptWithArgument_Arguments_FooBar(t *testing.T) {

	/*
		Set up activity Settings
	*/
	settingsJQ := &Settings{
		Script: ".foo |= $arg1",
	}

	initContextJQ := test.NewActivityInitContext(settingsJQ, nil)
	actJQ, err := New(initContextJQ)
	assert.Nil(t, err)

	tcJQ := test.NewActivityContext(actJQ.Metadata())

	argumentNames := CreateArgumentNames()
	tcJQ.SetInput("ArgumentNames", argumentNames)

	arguments := CreateArguments()
	tcJQ.SetInput("Arguments", arguments)

	tcJQ.SetInput("InputData", map[string]interface{}{"foo": "bar"})

	done, err := actJQ.Eval(tcJQ)
	assert.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, 1, len(tcJQ.GetOutput("OutputData").([]interface{})))
	//t.Logf("OutputData: %v", tcJQ.GetOutput("OutputData"))
	assert.Equal(t, map[string]interface{}{"foo": "baz"}, tcJQ.GetOutput("OutputData").([]interface{})[0])
	assert.Equal(t, false, tcJQ.GetOutput("Error").(bool))
	assert.Equal(t, "", tcJQ.GetOutput("ErrorMessage").(string))
}

func Test_executeJQ(t *testing.T) {
	type args struct {
		script    string
		inputData interface{}
		args      map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "simple identity",
			args: args{
				script:    ".",
				inputData: map[string]interface{}{"foo": "bar"},
				args:      nil,
			},
			want:    []interface{}{map[string]interface{}{"foo": "bar"}},
			wantErr: false,
		},
		{
			name: "extract field",
			args: args{
				script:    ".foo",
				inputData: map[string]interface{}{"foo": "bar"},
				args:      nil,
			},
			want:    []interface{}{"bar"},
			wantErr: false,
		},
		{
			name: "missing field",
			args: args{
				script:    ".baz",
				inputData: map[string]interface{}{"foo": "bar"},
				args:      nil,
			},
			want:    []interface{}{nil},
			wantErr: false,
		},
		{
			name: "use argument variable",
			args: args{
				script:    ". + $arg1",
				inputData: 1,
				args:      map[string]interface{}{"arg1": 2},
			},
			want:    []interface{}{3},
			wantErr: false,
		},
		{
			name: "don't use argument variable",
			args: args{
				script:    ".",
				inputData: 1,
				args:      map[string]interface{}{"arg1": 2},
			},
			want:    []interface{}{1},
			wantErr: false,
		},
		{
			name: "invalid jq script",
			args: args{
				script:    ".foo[",
				inputData: map[string]interface{}{"foo": "bar"},
				args:      nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "argument variable missing",
			args: args{
				script:    ". + $arg1",
				inputData: 1,
				args:      nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "array output",
			args: args{
				script:    "[.foo, .bar]",
				inputData: map[string]interface{}{"foo": 1, "bar": 2},
				args:      nil,
			},
			want:    []interface{}{[]interface{}{int(1), int(2)}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executeJQ(tt.args.script, tt.args.inputData, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeJQ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("executeJQ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecuteJQWithVariables(t *testing.T) {
	script := ".users[] | select(.age > $minAge) | .name"
	inputData := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 30},
			map[string]interface{}{"name": "Bob", "age": 25},
		},
	}
	args := map[string]interface{}{
		"minAge": 28,
	}

	results, err := executeJQ(script, inputData, args)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"Alice"}, results)
}

package jq

import (
	"fmt"

	"github.com/itchyny/gojq"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"gopkg.in/errgo.v2/fmt/errors"
)

func init() {
	_ = activity.Register(&JQActivity{}, New)
}

var activityLog = log.ChildLogger(log.RootLogger(), ACTIVITY_LOGGER)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)

	if err != nil {
		return nil, err
	}

	act := &JQActivity{settings: s}
	return act, nil
}

func (a *JQActivity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *JQActivity) Eval(context activity.Context) (done bool, err error) {
	// Initialize coercion manager
	coercionManager := NewCoercionManager(activityLog)

	script := a.settings.Script
	if script == "" {
		script = "."
	}

	input := &Input{}
	err = context.GetInputObject(input)
	if err != nil {
		return false, activity.NewError("Unable to get input object", "jq-activity-input-error", err)
	}

	activityLog.Debugf("Input Argument Names: %v", input.ArgumentNames)
	activityLog.Debugf("Input Arguments: %v", input.Arguments)

	// Use the coercion manager
	coercionResult, err := coercionManager.CoerceInputs(input.Arguments, input.ArgumentNames)
	if err != nil {
		return false, fmt.Errorf("failed to coerce inputs: %w", err)
	}

	if coercionResult.HasErrors() {
		return false, errors.Newf("Unable to coerce input ArgumentNames and Arguments: %v", coercionResult.Errors)
	}

	if coercionResult.Fields != nil {
		input.Arguments = coercionResult.Fields
	}

	activityLog.Debugf("Coerced input arguments: %v", input.Arguments)

	// ... rest of the function remains the same
	inputData, err := coerce.ToAny(input.InputData)
	if err != nil {
		return false, fmt.Errorf("unable to coerce input data to JSON: %w", err)
	}

	results, err := executeJQ(script, inputData, input.Arguments)
	if err != nil {
		return false, fmt.Errorf("unable to execute jq script: %w", err)
	}

	err = context.SetOutput("OutputData", results)
	if err != nil {
		return false, fmt.Errorf("unable to set output data: %w", err)
	}

	return true, nil
}

func executeJQ(script string, inputData interface{}, args map[string]interface{}) ([]interface{}, error) {
	// Parse the jq script
	query, err := gojq.Parse(script)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jq script: %w", err)
	}

	// Compile with or without arguments
	var code *gojq.Code
	var iter gojq.Iter

	if len(args) > 0 {
		// Extract variable names from the args map and prefix with "$"
		// Use a deterministic order by sorting the keys
		keys := make([]string, 0, len(args))
		for key := range args {
			keys = append(keys, key)
		}

		argVariables := make([]string, 0, len(keys))
		argValues := make([]interface{}, 0, len(keys))
		for _, key := range keys {
			argVariables = append(argVariables, "$"+key)
			argValues = append(argValues, args[key])
		}

		activityLog.Debugf("Compiling jq with variables: %v", argVariables)

		code, err = gojq.Compile(query, gojq.WithVariables(argVariables))
		if err != nil {
			return nil, fmt.Errorf("failed to compile jq script: %w", err)
		}

		activityLog.Debugf("Running jq with input data: %v and arguments: %v", inputData, argValues)

		iter = code.Run(inputData, argValues...)
	} else {
		activityLog.Debug("Compiling jq without variables")

		code, err = gojq.Compile(query)
		if err != nil {
			return nil, fmt.Errorf("failed to compile jq script: %w", err)
		}

		activityLog.Debugf("Running jq with input data: %v", inputData)

		iter = code.Run(inputData)
	}

	var results []interface{}

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, fmt.Errorf("jq execution error: %w", err)
		}
		results = append(results, v)
	}

	return results, nil
}

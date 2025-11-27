package jq

import (
	"fmt"

	"github.com/project-flogo/core/data/coerce"
)

// Coercer defines the interface for type coercion
type Coercer interface {
	Coerce(value interface{}) (interface{}, error)
}

// FieldDefinition defines the structure for input field metadata
type FieldDefinition struct {
	Name string
	Type string
}

// CoercionResult represents the result of coercion operations
type CoercionResult struct {
	Fields map[string]interface{}
	Errors []error
}

func (cr *CoercionResult) HasErrors() bool {
	return len(cr.Errors) > 0
}

func (cr *CoercionResult) AddError(err error) {
	cr.Errors = append(cr.Errors, err)
}

// CoercionManager handles type coercion with registry pattern
type CoercionManager struct {
	coercers map[string]Coercer
	logger   interface {
		Debugf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
	}
}

func NewCoercionManager(logger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}) *CoercionManager {
	cm := &CoercionManager{
		coercers: make(map[string]Coercer),
		logger:   logger,
	}

	// Register default coercers
	cm.RegisterCoercer("String", NewStringCoercer(logger))
	cm.RegisterCoercer("Number", NewNumberCoercer(logger))

	return cm
}

func (cm *CoercionManager) RegisterCoercer(typeName string, coercer Coercer) {
	cm.coercers[typeName] = coercer
}

// CoerceInputs performs type coercion on input arguments
func (cm *CoercionManager) CoerceInputs(
	inputArguments interface{},
	inputArgumentNames []interface{},
) (*CoercionResult, error) {

	if inputArgumentNames == nil || inputArguments == nil {
		return &CoercionResult{Fields: nil}, nil
	}

	// Handle coerce.ToObject which returns (map[string]interface{}, error)
	argumentFields, err := coerce.ToObject(inputArguments)
	if err != nil {
		cm.logger.Debugf("Failed to coerce input arguments to object: %v", err)
		return &CoercionResult{Fields: nil}, err
	}

	if argumentFields == nil {
		cm.logger.Debugf("Input arguments is nil after coercion")
		return &CoercionResult{Fields: nil}, nil
	}

	result := &CoercionResult{
		Fields: make(map[string]interface{}),
		Errors: []error{},
	}

	fieldDefs, err := cm.parseFieldDefinitions(inputArgumentNames)
	if err != nil {
		return nil, fmt.Errorf("failed to parse field definitions: %w", err)
	}

	for _, fieldDef := range fieldDefs {
		err := cm.coerceField(argumentFields, fieldDef, result)
		if err != nil {
			result.AddError(err)
		}
	}

	return result, nil
}

func (cm *CoercionManager) parseFieldDefinitions(inputArgumentNames []interface{}) ([]FieldDefinition, error) {
	var fieldDefs []FieldDefinition

	for _, v := range inputArgumentNames {
		field, err := coerce.ToObject(v)
		if err != nil || field == nil {
			cm.logger.Debugf("Failed to coerce field definition to object: %v", err)
			continue
		}

		// Handle coerce.ToString which returns (string, error)
		name, nameErr := coerce.ToString(field["Name"])
		if nameErr != nil {
			cm.logger.Debugf("Failed to coerce field name: %v", nameErr)
			continue
		}

		typeStr, typeErr := coerce.ToString(field["Type"])
		if typeErr != nil {
			cm.logger.Debugf("Failed to coerce field type: %v", typeErr)
			continue
		}

		if name == "" || typeStr == "" {
			cm.logger.Debugf("Field name or type is empty: name='%s', type='%s'", name, typeStr)
			continue
		}

		fieldDef := FieldDefinition{
			Name: name,
			Type: typeStr,
		}

		fieldDefs = append(fieldDefs, fieldDef)
		cm.logger.Debugf("Parsed field definition: %+v", fieldDef)
	}

	return fieldDefs, nil
}

// Also update the coerceField method to handle coerce.ToObject properly
func (cm *CoercionManager) coerceField(
	argumentFields map[string]interface{},
	fieldDef FieldDefinition,
	result *CoercionResult,
) error {

	cm.logger.Debugf("Processing field: %s (type: %s)", fieldDef.Name, fieldDef.Type)

	rawValue := argumentFields[fieldDef.Name]
	if rawValue == nil {
		result.Fields[fieldDef.Name] = nil
		return nil
	}

	coercer, exists := cm.coercers[fieldDef.Type]
	if !exists {
		err := fmt.Errorf("unsupported type %s for %s", fieldDef.Type, fieldDef.Name)
		cm.logger.Errorf("Coercion failed: %v", err)
		return err
	}

	fieldVal, err := coercer.Coerce(rawValue)
	if err != nil {
		cm.logger.Errorf("Coercion failed for %s (type %s) value %v: %v",
			fieldDef.Name, fieldDef.Type, rawValue, err)
		return err
	}

	cm.logger.Debugf("Coerced value: %v for field %s", fieldVal, fieldDef.Name)
	result.Fields[fieldDef.Name] = fieldVal
	return nil
}

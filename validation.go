package validation
// TODO: maybe rename to assert?

import (
	"errors"
	"fmt"
	"reflect"
)

type Subject interface {
	ValidationRules() RuleSet
}

// RuleSet key: fieldName value: ruleName
type RuleSet map[string][]Rule

// Go through all field rules defined by the user and run those rules against the field value
func (rs RuleSet) run(v reflect.Value) FailureSet{
	failures := make(FailureSet, len(rs))

	for fn, rules := range rs {
		field := v.FieldByName(fn)

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		for _, rule := range rules {
			if !rule.Supports(v.Interface()) {
				failures[fn] = append(failures[fn], fmt.Errorf("field '%s' does not support rule '%T'", field.Type().Name(), rule))
				continue
			}

			if !rule.IsValid(field.Interface()) {
				failures[fn] = append(failures[fn], fmt.Errorf("rule '%T' failed for field '%s'", rule, field.Type().Name()))
				continue
			}
		}

		if field.Kind() == reflect.Struct {
			nestedSubject, ok := field.Interface().(Subject)
			if !ok {
				continue
			}

			nestedFs := nestedSubject.ValidationRules().run(field)
			for nestedField, nestedF := range nestedFs {
				for _, err := range nestedF {
					fieldPath := fn + "." + nestedField
					failures[fieldPath] = append(failures[fieldPath], err)
				}
			}
		}

		if field.Kind() == reflect.Slice {
			for i := 0; i < field.Len(); i++ {
				nestedStruct := field.Index(i)
				nestedSubject, ok := nestedStruct.Interface().(Subject)
				if !ok {
					continue
				}

				nestedFs := nestedSubject.ValidationRules().run(nestedStruct)
				for nestedField, nestedF := range nestedFs {
					for _, err := range nestedF {
						fieldPath := fmt.Sprintf("%s.[%d].%s", fn, i, nestedField)
						failures[fieldPath] = append(failures[fieldPath], err)
					}
				}
			}
		}
	}

	return failures
}

type FailureSet map[string][]error

func Validate(input interface{}) (FailureSet, error) {
	var failSet FailureSet

	subject, ok := input.(Subject)
	if !ok {
		return failSet, errors.New("only instances of Subject can be validated")
	}

	inputValue := reflect.ValueOf(subject)

	if inputValue.Type().Kind() != reflect.Struct {
		return nil, errors.New("currently the validation package only supports validating structs")
	}

	return subject.ValidationRules().run(inputValue), nil
}

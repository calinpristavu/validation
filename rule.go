package validation

import (
	"reflect"
)

type Rule interface {
	Supports(i interface{}) bool
	IsValid(i interface{}) bool
}

type NotZeroValue struct {

}

func (NotZeroValue) Supports(_ interface{}) bool {
	return true
}

func (NotZeroValue) IsValid(i interface{}) bool {
	return !reflect.ValueOf(i).IsZero()
}

type NotNil struct {

}

func (NotNil) Supports(_ interface{}) bool {
	return true
}

func (NotNil) IsValid(i interface{}) bool {
	return reflect.ValueOf(i).IsNil()
}

type NestedValid struct {

}

func (NestedValid) Supports(i interface{}) bool {
	_, ok := i.(Subject)

	return ok
}

func (NestedValid) IsValid(_ interface{}) bool {
	return true
}

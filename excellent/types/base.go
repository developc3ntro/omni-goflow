package types

import (
	"fmt"
	"reflect"

	"github.com/nyaruka/goflow/utils"
)

// XValue is the base interface of all Excellent types
type XValue interface {
	fmt.Stringer

	Describe(env utils.Environment) string
	ToXText(env utils.Environment) XText
	ToXBoolean(env utils.Environment) XBoolean
	ToXJSON(env utils.Environment) XText
}

// XLengthable is the interface for types which have a length
type XLengthable interface {
	Length() int
}

// Equals checks for equality between the two give values. This is only used for testing as x = y
// specifically means text(x) == text(y)
func Equals(x1 XValue, x2 XValue) bool {
	// nil == nil
	if utils.IsNil(x1) && utils.IsNil(x2) {
		return true
	} else if utils.IsNil(x1) || utils.IsNil(x2) {
		return false
	}

	// different types aren't equal
	if reflect.TypeOf(x1) != reflect.TypeOf(x2) {
		return false
	}

	// common types, do real comparisons
	switch typed := x1.(type) {
	case *XArray:
		return typed.Equals(x2.(*XArray))
	case XBoolean:
		return typed.Equals(x2.(XBoolean))
	case XDate:
		return typed.Equals(x2.(XDate))
	case XDateTime:
		return typed.Equals(x2.(XDateTime))
	case *XDict:
		return typed.Equals(x2.(*XDict))
	case XError:
		return typed.Equals(x2.(XError))
	case XFunction:
		return typed.Equals(x2.(XFunction))
	case XNumber:
		return typed.Equals(x2.(XNumber))
	case XText:
		return typed.Equals(x2.(XText))
	case XTime:
		return typed.Equals(x2.(XTime))
	default:
		panic(fmt.Sprintf("can't compare equality of instances of %T", x1))
	}
}

// IsEmpty determines if the given value is empty
func IsEmpty(x XValue) bool {
	// nil is empty
	if utils.IsNil(x) {
		return true
	}

	// anything with length of zero is empty
	asLengthable, isLengthable := x.(XLengthable)
	if isLengthable && asLengthable.Length() == 0 {
		return true
	}

	return false
}

// String returns a representation of the given value for use in debugging
func String(x XValue) string {
	if utils.IsNil(x) {
		return "nil"
	}
	return x.String()
}

// Describe returns a representation of the given value for use in error messages
func Describe(env utils.Environment, x XValue) string {
	if utils.IsNil(x) {
		return "null"
	}
	return x.Describe(env)
}

// XRepresentable is the interface for any object which can be represented in an expression
type XRepresentable interface {
	ToXValue(env utils.Environment) XValue
}

// ToXValue is a utility to convert the given XRepresentable to an XValue
func ToXValue(env utils.Environment, obj XRepresentable) XValue {
	if utils.IsNil(obj) {
		return nil
	}
	return obj.ToXValue(env)
}

package scan

import (
	"fmt"
	"strconv"
)

const (
	loxValueStringType = iota
	loxValueFloatType
	loxValueBoolType
	loxValueNilType
	loxCallableType
	loxClassType
	loxClassInstanceType
)

type NilObject struct{}

type LoxValue struct {
	valueType      int
	stringValue    string
	floatValue     float64
	boolValue      bool
	nilObject      NilObject
	callableObject LoxCallable
	classObject    LoxClass
	classInstance  LoxClassInstance
}

func (l *LoxValue) IsNumber() bool {
	return l.valueType == loxValueFloatType
}

func (l *LoxValue) IsString() bool {
	return l.valueType == loxValueStringType
}

func (l *LoxValue) IsBoolean() bool {
	return l.valueType == loxValueBoolType
}

func (l *LoxValue) IsNil() bool {
	return l.valueType == loxValueNilType
}

func (l *LoxValue) IsCallable() bool {
	return l.valueType == loxCallableType
}

func (l *LoxValue) IsClass() bool {
	return l.valueType == loxClassType
}

func (l *LoxValue) IsClassInstance() bool {
	return l.valueType == loxClassInstanceType
}

func (l *LoxValue) ClassInstance() (LoxClassInstance, error) {
	switch l.valueType {
	case loxClassInstanceType:
		return l.classInstance, nil
	case loxCallableType:
		return nil, fmt.Errorf("function is not a class instance")
	case loxValueFloatType:
		return nil, fmt.Errorf("float is not a class instance")
	case loxValueStringType:
		return nil, fmt.Errorf("string is not a class instance")
	case loxValueBoolType:
		return nil, fmt.Errorf("bool is not a class instance")
	case loxValueNilType:
		return nil, fmt.Errorf("nil is not a class instance")
	case loxClassType:
		return nil, fmt.Errorf("class is not a class instance")
	}

	// unreachable
	panic("use not implemented value type")
}

func (l *LoxValue) Callable() (LoxCallable, error) {
	switch l.valueType {
	case loxCallableType:
		return l.callableObject, nil
	case loxClassType:
		return l.classObject, nil
	case loxValueFloatType:
		return nil, fmt.Errorf("float is not a function")
	case loxValueStringType:
		return nil, fmt.Errorf("string is not a function")
	case loxValueBoolType:
		return nil, fmt.Errorf("bool is not a function")
	case loxValueNilType:
		return nil, fmt.Errorf("nil is not a function")
	case loxClassInstanceType:
		return nil, fmt.Errorf("class is not a function")
	}

	// unreachable
	panic("use not implemented value type")
}

func (l *LoxValue) Number() (float64, error) {
	switch l.valueType {
	case loxValueFloatType:
		return l.floatValue, nil
	case loxValueStringType:
		return 0., fmt.Errorf("string is not a number")
	case loxValueBoolType:
		return 0., fmt.Errorf("bool is not a number")
	case loxValueNilType:
		return 0., fmt.Errorf("nil is not a number")
	case loxCallableType:
		return 0., fmt.Errorf("function is not a number")
	case loxClassInstanceType, loxClassType:
		return 0., fmt.Errorf("class is not a number")
	}

	// unreachable
	panic("use not implemented value type")
}

func (l *LoxValue) String() string {
	switch l.valueType {
	case loxValueStringType:
		return l.stringValue
	case loxValueBoolType:
		return strconv.FormatBool(l.boolValue)
	case loxValueFloatType:
		return strconv.FormatFloat(l.floatValue, 'f', -1, 64)
	case loxValueNilType:
		return "nil"
	case loxCallableType:
		return l.callableObject.String()
	case loxClassType:
		return l.classObject.String()
	case loxClassInstanceType:
		return l.classInstance.String()
	}

	// unreachable
	panic("use not implemented value type")
}

func (l *LoxValue) Bool() bool {
	switch l.valueType {
	case loxValueBoolType:
		return l.boolValue
	case loxValueFloatType:
		return l.floatValue > 0.
	default:
		return false
	}
}

func NewStringLoxValue(value string) *LoxValue {
	return &LoxValue{
		valueType:   loxValueStringType,
		stringValue: value,
	}
}

func NewFloatLoxValue(value float64) *LoxValue {
	return &LoxValue{
		valueType:  loxValueFloatType,
		floatValue: value,
	}
}

func NewBooleanLoxValue(value bool) *LoxValue {
	return &LoxValue{
		valueType: loxValueBoolType,
		boolValue: value,
	}
}

func NewNilLoxValue() *LoxValue {
	return &LoxValue{
		valueType: loxValueNilType,
		nilObject: NilObject{},
	}
}

func NewCallableLoxValue(callable LoxCallable) *LoxValue {
	return &LoxValue{
		valueType:      loxCallableType,
		callableObject: callable,
	}
}

func NewClassLoxValue(callable LoxClass) *LoxValue {
	return &LoxValue{
		valueType:   loxClassType,
		classObject: callable,
	}
}

func NewClassInstanceLoxValue(classInstance LoxClassInstance) *LoxValue {
	return &LoxValue{
		valueType:     loxClassInstanceType,
		classInstance: classInstance,
	}
}

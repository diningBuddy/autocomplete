package env

import (
	"fmt"
	"reflect"
)

// KeyNotFoundError is used when the key does not exist in environment variables
type KeyNotFoundError struct {
	Key string
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("key %s is not found in environment", e.Key)
}

// InvalidValueError is used when the value is invalid type so unable to transform to expected type
type InvalidValueError struct {
	Key        string
	Value      string
	ExpectType reflect.Kind
}

func (e InvalidValueError) Error() string {
	return fmt.Sprintf("value %v of key %s is an invalid value for type %v", e.Value, e.Key, e.ExpectType)
}

// InvalidProfileError is used when the value is invalid type so unable to transform to proper Profile type
type InvalidProfileError struct {
	Key   string
	Value string
}

func (e InvalidProfileError) Error() string {
	return fmt.Sprintf("value (%s) of key %s is invalid string for profile", e.Value, e.Key)
}

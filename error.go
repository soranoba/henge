package henge

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrInvalidValue is an error when the source is an invalid value.
	// Refer: reflect.IsValid
	ErrInvalidValue = errors.New("invalid value")
	// ErrUnsupportedType is an error when the source is an unsupported type.
	ErrUnsupportedType = errors.New("unsupported type")
	// ErrOverflow is an error if an overflow occurs during conversion.
	ErrOverflow = errors.New("overflows")
	// ErrNegativeNumber is an error if converting a negative number to an unsigned type.
	ErrNegativeNumber = errors.New("negative number")
)

// ConvertError is an error that shows where the error occurred during conversion.
type ConvertError struct {
	Field   string
	SrcType reflect.Type
	DstType reflect.Type
	Value   interface{}
	Err     error
}

func (e *ConvertError) Unwrap() error {
	return e.Err
}

func (e *ConvertError) Error() string {
	srcTypeString, dstTypeString := "nil", "nil"
	if e.SrcType != nil {
		srcTypeString = e.SrcType.String()
	}
	if e.DstType != nil {
		dstTypeString = e.DstType.String()
	}

	return fmt.Sprintf(
		"Failed to convert from %s to %s: fields=%s, value=%#v, error=%s",
		srcTypeString, dstTypeString, e.Field, e.Value, e.Err.Error(),
	)
}

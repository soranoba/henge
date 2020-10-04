package henge

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	invalidValueErr    = errors.New("invalid value")
	unsupportedTypeErr = errors.New("unsupported type")
)

type ConvertError struct {
	Field   string
	SrcType reflect.Type
	DstType reflect.Type
	Err     error
}

func (e *ConvertError) Unwrap() error {
	return e.Err
}

func (e *ConvertError) Error() string {
	return fmt.Sprintf(
		"Failed to convert from %s to %s: fields=%s, error=%s",
		e.SrcType.String(), e.DstType.String(), e.Field, e.Err.Error(),
	)
}

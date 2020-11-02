package henge

import (
	"reflect"
)

// Bool converts the input to bool type.
func (c *ValueConverter) Bool() *BoolConverter {
	var (
		value bool
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	if inV.IsValid() {
		switch inV.Type().Kind() {
		case reflect.Bool:
			value = inV.Interface().(bool)
		default:
			value = !inV.IsZero()
		}
	} else {
		err = ErrInvalidValue
	}

	if err != nil {
		var srcType reflect.Type
		if reflect.ValueOf(c.value).IsValid() {
			srcType = reflect.ValueOf(c.value).Type()
		}
		err = &ConvertError{
			Field:   c.field,
			SrcType: srcType,
			DstType: reflect.ValueOf((*bool)(nil)).Type().Elem(),
			Value:   c.value,
			Err:     err,
		}
	}
	return &BoolConverter{converter: c.converter, value: value, err: err}
}

// BoolPtr converts the input to pointer of bool type.
func (c *ValueConverter) BoolPtr() *BoolPtrConverter {
	return c.Bool().Ptr()
}

// BoolConverter is a converter that converts a bool type to another type.
type BoolConverter struct {
	converter
	value bool
	err   error
}

// Ptr converts the input to ptr type.
func (c *BoolConverter) Ptr() *BoolPtrConverter {
	if c.err != nil || c.isNil {
		return &BoolPtrConverter{converter: c.converter, value: nil, err: c.err}
	}
	return &BoolPtrConverter{converter: c.converter, value: &c.value, err: c.err}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *BoolConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}

	if c.err != nil {
		return c.err
	}
	if c.isNil {
		return nil
	}

	for outV.Kind() == reflect.Ptr {
		if outV.IsNil() {
			outV.Set(reflect.New(outV.Type().Elem()))
		}
		outV = outV.Elem()
	}

	switch outV.Kind() {
	case reflect.Bool:
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	default:
		return c.new(c.value, c.field).Convert(out)
	}
	return nil
}

// Result returns the conversion result and error.
func (c *BoolConverter) Result() (bool, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *BoolConverter) Value() bool {
	return c.value
}

// Error returns an error if the conversion fails
func (c *BoolConverter) Error() error {
	return c.err
}

// BoolPtrConverter is a converter that converts a pointer of bool type to another type.
type BoolPtrConverter struct {
	converter
	value *bool
	err   error
}

// Result returns the conversion result and error.
func (c *BoolPtrConverter) Result() (*bool, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *BoolPtrConverter) Value() *bool {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *BoolPtrConverter) Error() error {
	return c.err
}

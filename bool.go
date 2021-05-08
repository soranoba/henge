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
		err = c.wrapConvertError(c.value, reflect.ValueOf((*bool)(nil)).Type().Elem(), err)
	}
	return &BoolConverter{baseConverter: c.baseConverter, value: value, err: err}
}

// BoolPtr converts the input to pointer of bool type.
func (c *ValueConverter) BoolPtr() *BoolPtrConverter {
	return c.Bool().Ptr()
}

// BoolConverter is a converter that converts a bool type to another type.
type BoolConverter struct {
	baseConverter
	value bool
	err   error
}

// Ptr converts the input to ptr type.
func (c *BoolConverter) Ptr() *BoolPtrConverter {
	if c.err != nil || c.isNil {
		return &BoolPtrConverter{baseConverter: c.baseConverter, value: nil, err: c.err}
	}
	return &BoolPtrConverter{baseConverter: c.baseConverter, value: &c.value, err: c.err}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *BoolConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *BoolConverter) convert(outV reflect.Value) error {
	if c.err != nil {
		return c.wrapConvertError(c.value, outV.Type(), c.err)
	}
	if c.isNil {
		return nil
	}

	elemOutV := toInitializedNonPtrValue(outV)

	switch elemOutV.Kind() {
	case reflect.Bool:
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	default:
		return c.new(c.value, c.field).Convert(outV)
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

// Interface returns the conversion result of interface type.
func (c *BoolConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails
func (c *BoolConverter) Error() error {
	return c.err
}

// BoolPtrConverter is a converter that converts a pointer of bool type to another type.
type BoolPtrConverter struct {
	baseConverter
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

// Interface returns the conversion result of interface type.
func (c *BoolPtrConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *BoolPtrConverter) Error() error {
	return c.err
}

package henge

import (
	"reflect"
	"strconv"
)

// Float converts the input to float type.
func (c *ValueConverter) Float() *FloatConverter {
	var (
		value float64
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	if inV.IsValid() {
		inT := inV.Type()
		outT := reflect.TypeOf(value)
		if inT.ConvertibleTo(outT) {
			value = inV.Convert(outT).Interface().(float64)
		} else if inT.Kind() == reflect.String {
			value, err = strconv.ParseFloat(inV.Interface().(string), 64)
		} else if inT.Kind() == reflect.Bool {
			if inV.Interface().(bool) == true {
				value = 1
			}
		} else {
			err = ErrUnsupportedType
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
			DstType: reflect.ValueOf((*float64)(nil)).Type().Elem(),
			Value:   c.value,
			Err:     err,
		}
	}
	return &FloatConverter{converter: c.converter, value: value, err: err}
}

// FloatrPtr converts the input to pointer of float type.
func (c *ValueConverter) FloatrPtr() *FloatPtrConverter {
	return c.Float().Ptr()
}

// FloatConverter is a converter that converts a float type to another type.
type FloatConverter struct {
	converter
	value float64
	err   error
}

// Int converts the input to int type.
func (c *FloatConverter) Int() *IntegerConverter {
	if c.err != nil {
		return &IntegerConverter{converter: c.converter, value: 0, err: c.err}
	}
	return c.new(c.value, c.field).Int()
}

// Uint converts the input to uint type.
func (c *FloatConverter) Uint() *UnsignedIntegerConverter {
	if c.err != nil {
		return &UnsignedIntegerConverter{converter: c.converter, value: 0, err: c.err}
	}
	return c.new(c.value, c.field).Uint()
}

// Ptr converts the input to ptr type.
func (c *FloatConverter) Ptr() *FloatPtrConverter {
	if c.err != nil || c.isNil {
		return &FloatPtrConverter{converter: c.converter, value: nil, err: c.err}
	}
	return &FloatPtrConverter{converter: c.converter, value: &c.value, err: c.err}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *FloatConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *FloatConverter) convert(outV reflect.Value) error {
	if c.err != nil {
		err := *(c.err.(*ConvertError))
		err.DstType = outV.Type()
		return &err
	}
	if c.isNil {
		return nil
	}

	elemOutV := toInitializedNonPtrValue(outV)

	switch elemOutV.Kind() {
	case reflect.Float32:
		if float64(float32(c.value)) != c.value {
			return &ConvertError{
				Field:   c.field,
				SrcType: reflect.ValueOf(c.value).Type(),
				DstType: outV.Type(),
				Value:   c.value,
				Err:     ErrOverflow,
			}
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Float64:
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	default:
		return c.new(c.value, c.field).convert(outV)
	}
	return nil
}

// Result returns the conversion result and error.
func (c *FloatConverter) Result() (float64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *FloatConverter) Value() float64 {
	return c.value
}

// Error returns an error if the conversion fails
func (c *FloatConverter) Error() error {
	return c.err
}

// FloatPtrConverter is a converter that converts a pointer of float type to another type.
type FloatPtrConverter struct {
	converter
	value *float64
	err   error
}

// Result returns the conversion result and error.
func (c *FloatPtrConverter) Result() (*float64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *FloatPtrConverter) Value() *float64 {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *FloatPtrConverter) Error() error {
	return c.err
}

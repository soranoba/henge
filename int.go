package henge

import (
	"math"
	"reflect"
	"strconv"
)

// Int converts the input to int type.
func (c *ValueConverter) Int() *IntegerConverter {
	var (
		value int64
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	if inV.IsValid() {
		inT := inV.Type()
		outT := reflect.TypeOf(value)
		switch inT.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = inV.Convert(outT).Interface().(int64)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var u uint64
			u = inV.Convert(reflect.ValueOf(u).Type()).Interface().(uint64)
			i := inV.Convert(outT).Interface().(int64)
			if i < 0 && u > 0 {
				err = ErrOverflow
			} else {
				value = i
			}
		case reflect.Float32, reflect.Float64:
			var f float64
			f = inV.Convert(reflect.ValueOf(f).Type()).Interface().(float64)
			f = c.opts.numOpts.roundingFunc(f)
			if f > math.MaxInt64 || f < math.MinInt64 || ((f > 0) != (int64(f) > 0)) {
				err = ErrOverflow
			} else {
				value = int64(f)
			}
		case reflect.Bool:
			if inV.Interface().(bool) == true {
				value = 1
			}
		case reflect.String:
			value, err = strconv.ParseInt(inV.Interface().(string), 10, 64)
		default:
			err = ErrUnsupportedType
		}
	} else {
		err = ErrInvalidValue
	}

	if err != nil {
		err = c.wrapConvertError(c.value, reflect.ValueOf((*int64)(nil)).Type().Elem(), err)
	}
	return &IntegerConverter{baseConverter: c.baseConverter, value: value, err: err}
}

// IntPtr converts the input to pointer of int type.
func (c *ValueConverter) IntPtr() *IntegerPtrConverter {
	return c.Int().Ptr()
}

// IntegerConverter is a converter that converts an integer type to another type.
type IntegerConverter struct {
	*baseConverter
	value int64
	err   error
}

// Ptr converts the input to ptr type.
func (c *IntegerConverter) Ptr() *IntegerPtrConverter {
	if c.err != nil || c.isNil {
		return &IntegerPtrConverter{baseConverter: c.baseConverter, value: nil, err: c.err}
	}
	return &IntegerPtrConverter{baseConverter: c.baseConverter, value: &c.value, err: nil}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *IntegerConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Type().Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *IntegerConverter) convert(outV reflect.Value) error {
	if c.err != nil {
		return c.wrapConvertError(c.value, outV.Type(), c.err)
	}
	if c.isNil {
		return nil
	}

	elemOutV := toInitializedNonPtrValue(outV)
	overflowErr := c.wrapConvertError(c.value, outV.Type(), ErrOverflow)

	switch elemOutV.Kind() {
	case reflect.Int:
		if int64(int(c.value)) != c.value {
			return overflowErr
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Int8:
		if int64(int8(c.value)) != c.value {
			return overflowErr
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Int16:
		if int64(int16(c.value)) != c.value {
			return overflowErr
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Int32:
		if int64(int32(c.value)) != c.value {
			return overflowErr
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Int64:
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	default:
		return c.new(c.value, c.field).convert(outV)
	}
	return nil
}

// Result returns the conversion result and error.
func (c *IntegerConverter) Result() (int64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *IntegerConverter) Value() int64 {
	return c.value
}

// Interface returns the conversion result of interface type.
func (c *IntegerConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *IntegerConverter) Error() error {
	return c.err
}

// IntegerPtrConverter is a converter that converts a pointer of integer type to another type.
type IntegerPtrConverter struct {
	*baseConverter
	value *int64
	err   error
}

// Result returns the conversion result and error.
func (c *IntegerPtrConverter) Result() (*int64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *IntegerPtrConverter) Value() *int64 {
	return c.value
}

// Interface returns the conversion result of interface type.
func (c *IntegerPtrConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *IntegerPtrConverter) Error() error {
	return c.err
}

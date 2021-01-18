package henge

import (
	"math"
	"reflect"
	"strconv"
)

// Uint converts the input to uint type.
func (c *ValueConverter) Uint() *UnsignedIntegerConverter {
	var (
		value uint64
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	if inV.IsValid() {
		inT := inV.Type()
		outT := reflect.TypeOf(value)
		switch inT.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i = inV.Convert(reflect.ValueOf(i).Type()).Interface().(int64)
			if i < 0 {
				err = ErrNegativeNumber
			} else {
				value = uint64(i)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = inV.Convert(outT).Interface().(uint64)
		case reflect.Float32, reflect.Float64:
			var f float64
			f = inV.Convert(reflect.ValueOf(f).Type()).Interface().(float64)
			f = math.Floor(f)
			if f < 0 {
				err = ErrNegativeNumber
			} else if f > math.MaxUint64 || ((f > 0) != (int64(f) > 0)) {
				err = ErrOverflow
			} else {
				value = uint64(f)
			}
		case reflect.Bool:
			if inV.Interface().(bool) == true {
				value = 1
			}
		case reflect.String:
			value, err = strconv.ParseUint(inV.Interface().(string), 10, 64)
		default:
			err = ErrUnsupportedType
		}
	} else {
		err = ErrInvalidValue
	}

	if err != nil {
		err = c.wrapConvertError(c.value, reflect.ValueOf((*uint64)(nil)).Type().Elem(), err)
	}
	return &UnsignedIntegerConverter{converter: c.converter, value: value, err: err}
}

// UintPtr converts the input to pointer of uint type.
func (c *ValueConverter) UintPtr() *UnsignedIntegerPtrConverter {
	return c.Uint().Ptr()
}

// UnsignedIntegerConverter is a converter that converts an unsigned integer type to another type.
type UnsignedIntegerConverter struct {
	converter
	value uint64
	err   error
}

// Ptr converts the input to ptr type.
func (c *UnsignedIntegerConverter) Ptr() *UnsignedIntegerPtrConverter {
	if c.err != nil || c.isNil {
		return &UnsignedIntegerPtrConverter{converter: c.converter, value: nil, err: c.err}
	}
	return &UnsignedIntegerPtrConverter{converter: c.converter, value: &c.value, err: nil}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *UnsignedIntegerConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *UnsignedIntegerConverter) convert(outV reflect.Value) error {
	if c.err != nil {
		return c.wrapConvertError(c.value, outV.Type(), c.err)
	}
	if c.isNil {
		return nil
	}

	elemOutV := toInitializedNonPtrValue(outV)
	overflowErr := c.wrapConvertError(c.value, outV.Type(), ErrOverflow)

	switch elemOutV.Kind() {
	case reflect.Uint:
		if uint64(uint(c.value)) != c.value {
			return overflowErr
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Uint8:
		if uint64(uint8(c.value)) != c.value {
			return overflowErr
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Uint16:
		if uint64(uint16(c.value)) != c.value {
			return overflowErr
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Uint32:
		if uint64(uint32(c.value)) != c.value {
			return overflowErr
		}
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	case reflect.Uint64:
		elemOutV.Set(reflect.ValueOf(c.value).Convert(elemOutV.Type()))
	default:
		return c.new(c.value, c.field).convert(outV)
	}
	return nil
}

// Result returns the conversion result and error.
func (c *UnsignedIntegerConverter) Result() (uint64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *UnsignedIntegerConverter) Value() uint64 {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *UnsignedIntegerConverter) Error() error {
	return c.err
}

// UnsignedIntegerPtrConverter is a converter that converts a pointer of uint type to another type.
type UnsignedIntegerPtrConverter struct {
	converter
	value *uint64
	err   error
}

// Result returns the conversion result and error.
func (c *UnsignedIntegerPtrConverter) Result() (*uint64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *UnsignedIntegerPtrConverter) Value() *uint64 {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *UnsignedIntegerPtrConverter) Error() error {
	return c.err
}

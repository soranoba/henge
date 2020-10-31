package henge

import (
	"errors"
	"math"
	"reflect"
	"strconv"
)

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
				err = errors.New("negative number")
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
				err = errors.New("negative number")
			} else if f > math.MaxUint64 || ((f > 0) != (int64(f) > 0)) {
				err = errors.New("overflows")
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
			err = unsupportedTypeErr
		}
	} else {
		err = invalidValueErr
	}

	if err != nil {
		err = &ConvertError{
			Field:   c.field,
			SrcType: reflect.ValueOf(c.value).Type(),
			DstType: reflect.ValueOf((*uint64)(nil)).Type().Elem(),
			Value:   c.value,
			Err:     err,
		}
	}
	return &UnsignedIntegerConverter{converter: c.converter, value: value, err: err}
}

type UnsignedIntegerConverter struct {
	converter
	value uint64
	err   error
}

func (c *UnsignedIntegerConverter) Ptr() *UnsignedIntegerPtrConverter {
	if c.err != nil {
		return &UnsignedIntegerPtrConverter{converter: c.converter, value: nil, err: c.err}
	}
	return &UnsignedIntegerPtrConverter{converter: c.converter, value: &c.value, err: nil}
}

func (c *UnsignedIntegerConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}

	if c.err != nil {
		return c.err
	}

	for outV.Kind() == reflect.Ptr {
		outV = outV.Elem()
	}

	switch outV.Kind() {
	case reflect.Uint:
		if uint64(uint(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Uint8:
		if uint64(uint8(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Uint16:
		if uint64(uint16(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Uint32:
		if uint64(uint32(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Uint64:
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	default:
		return c.new(c.value, c.field).Convert(out)
	}
	return nil
}

func (c *UnsignedIntegerConverter) Result() (uint64, error) {
	return c.value, c.err
}

func (c *UnsignedIntegerConverter) Value() uint64 {
	return c.value
}

func (c *UnsignedIntegerConverter) Error() error {
	return c.err
}

type UnsignedIntegerPtrConverter struct {
	converter
	value *uint64
	err   error
}

func (c *UnsignedIntegerPtrConverter) Result() (*uint64, error) {
	return c.value, c.err
}

func (c *UnsignedIntegerPtrConverter) Value() *uint64 {
	return c.value
}

func (c *UnsignedIntegerPtrConverter) Error() error {
	return c.err
}

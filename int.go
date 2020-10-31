package henge

import (
	"errors"
	"math"
	"reflect"
	"strconv"
)

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
				err = errors.New("overflows")
			} else {
				value = i
			}
		case reflect.Float32, reflect.Float64:
			var f float64
			f = inV.Convert(reflect.ValueOf(f).Type()).Interface().(float64)
			f = math.Floor(f)
			if f > math.MaxInt64 || f < math.MinInt64 || ((f > 0) != (int64(f) > 0)) {
				err = errors.New("overflows")
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
			err = unsupportedTypeErr
		}
	} else {
		err = invalidValueErr
	}

	if err != nil {
		err = &ConvertError{
			Field:   c.field,
			SrcType: reflect.ValueOf(c.value).Type(),
			DstType: reflect.ValueOf((*int64)(nil)).Type().Elem(),
			Value:   c.value,
			Err:     err,
		}
	}
	return &IntegerConverter{converter: c.converter, value: value, err: err}
}

type IntegerConverter struct {
	converter
	value int64
	err   error
}

func (c *IntegerConverter) Ptr() *IntegerPtrConverter {
	if c.err != nil {
		return &IntegerPtrConverter{converter: c.converter, value: nil, err: c.err}
	}
	return &IntegerPtrConverter{converter: c.converter, value: &c.value, err: nil}
}

func (c *IntegerConverter) Convert(out interface{}) error {
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
	case reflect.Int:
		if int64(int(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Int8:
		if int64(int8(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Int16:
		if int64(int16(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Int32:
		if int64(int32(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Int64:
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	default:
		return c.new(c.value, c.field).Convert(out)
	}
	return nil
}

func (c *IntegerConverter) Result() (int64, error) {
	return c.value, c.err
}

func (c *IntegerConverter) Value() int64 {
	return c.value
}

func (c *IntegerConverter) Error() error {
	return c.err
}

type IntegerPtrConverter struct {
	converter
	value *int64
	err   error
}

func (c *IntegerPtrConverter) Result() (*int64, error) {
	return c.value, c.err
}

func (c *IntegerPtrConverter) Value() *int64 {
	return c.value
}

func (c *IntegerPtrConverter) Error() error {
	return c.err
}

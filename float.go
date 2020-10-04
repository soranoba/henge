package henge

import (
	"errors"
	"reflect"
	"strconv"
)

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
			err = unsupportedTypeErr
		}
	} else {
		err = invalidValueErr
	}

	if err != nil {
		err = &ConvertError{
			Field:   "",
			SrcType: reflect.ValueOf(c.value).Type(),
			DstType: reflect.ValueOf((*float64)(nil)).Type().Elem(),
			Err:     err,
		}
	}
	return &FloatConverter{value: value, err: err}
}

type FloatConverter struct {
	value float64
	err   error
}

func (c *FloatConverter) Int() *IntegerConverter {
	if c.err != nil {
		return &IntegerConverter{value: 0, err: c.err}
	}
	return New(c.value).Int()
}

func (c *FloatConverter) Uint() *UnsignedIntegerConverter {
	if c.err != nil {
		return &UnsignedIntegerConverter{value: 0, err: c.err}
	}
	return New(c.value).Uint()
}

func (c *FloatConverter) Ptr() *FloatPtrConverter {
	if c.err != nil {
		return &FloatPtrConverter{value: nil, err: c.err}
	}
	return &FloatPtrConverter{value: &c.value, err: c.err}
}

func (c *FloatConverter) Convert(out interface{}) error {
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
	case reflect.Float32:
		if float64(float32(c.value)) != c.value {
			return errors.New("overflows")
		}
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	case reflect.Float64:
		outV.Set(reflect.ValueOf(c.value).Convert(outV.Type()))
	default:
		return New(c.value).Convert(out)
	}
	return nil
}

func (c *FloatConverter) Result() (float64, error) {
	return c.value, c.err
}

func (c *FloatConverter) Value() float64 {
	return c.value
}

func (c *FloatConverter) Error() error {
	return c.err
}

type FloatPtrConverter struct {
	value *float64
	err   error
}

func (c *FloatPtrConverter) Result() (*float64, error) {
	return c.value, c.err
}

func (c *FloatPtrConverter) Value() *float64 {
	return c.value
}

func (c *FloatPtrConverter) Error() error {
	return c.err
}

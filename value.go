package henge

import (
	"reflect"
)

type ValueConverter struct {
	converter
	value interface{}
	err   error
}

func New(i interface{}, fs ...func(*ConverterOpts)) *ValueConverter {
	return &ValueConverter{converter: newConverter(fs...), value: i, err: nil}
}

func (c *ValueConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}

	for outV.Kind() == reflect.Ptr {
		outV = outV.Elem()
	}

	switch outV.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return c.Int().Convert(out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return c.Uint().Convert(out)
	case reflect.Float32, reflect.Float64:
		return c.Float().Convert(out)
	case reflect.String:
		return c.String().Convert(out)
	case reflect.Array, reflect.Slice:
		return c.Slice().Convert(out)
	case reflect.Struct:
		return c.Struct().Convert(out)
	default:
		return unsupportedTypeErr
	}
}

func (c *ValueConverter) Result() (interface{}, error) {
	return c.value, c.err
}

func (c *ValueConverter) Value() interface{} {
	return c.value
}

func (c *ValueConverter) Error() error {
	return c.err
}

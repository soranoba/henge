package henge

import (
	"reflect"
)

// Converter is an interface that has common functions of all Converters.
type Converter interface {
	InstanceGet(key string) interface{}
	InstanceSet(key string, value interface{})
}

type converter struct {
	isNil   bool
	field   string
	opts    ConverterOpts
	storage map[string]interface{}
}

func (c *converter) new(i interface{}, fieldName string) *ValueConverter {
	newConverter := New(i)
	newConverter.converter = *c
	newConverter.converter.field = fieldName
	return newConverter
}

// InstanceGet returns the value saved using InstanceSet.
func (c *converter) InstanceGet(key string) interface{} {
	return c.storage[key]
}

// InstanceSet saves the value by specifying the key.
func (c *converter) InstanceSet(key string, value interface{}) {
	c.storage[key] = value
}

// ValueConverter is a converter that converts an interface type to another type.
type ValueConverter struct {
	converter
	value interface{}
	err   error
}

// New returns a new ValueConverter
func New(i interface{}, fs ...func(*ConverterOpts)) *ValueConverter {
	opts := defaultConverterOpts()
	for _, f := range fs {
		f(&opts)
	}

	inV := reflect.ValueOf(i)
	isNil := false
	switch inV.Kind() {
	case reflect.Ptr:
		if isNil = inV.IsNil(); isNil {
			i = reflect.New(inV.Type().Elem()).Interface()
		}
	case reflect.Slice, reflect.Map:
		if isNil = inV.IsNil(); isNil {
			i = reflect.New(inV.Type()).Interface()
		}
	}

	return &ValueConverter{
		converter: converter{
			isNil:   isNil,
			opts:    opts,
			storage: map[string]interface{}{},
		},
		value: i,
		err:   nil,
	}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
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
		return ErrUnsupportedType
	}
}

// Result returns the conversion result and error.
func (c *ValueConverter) Result() (interface{}, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *ValueConverter) Value() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *ValueConverter) Error() error {
	return c.err
}

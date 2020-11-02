package henge

import (
	"reflect"
)

// Converter is an interface that has common functions of all Converters.
type Converter interface {
	InstanceGet(key string) interface{}
	InstanceSet(key string, value interface{})
	InstanceSetValues(m map[string]interface{})
}

type converter struct {
	isNil   bool
	field   string
	opts    ConverterOpts
	storage map[string]interface{}
}

func (c *converter) new(i interface{}, fieldName string) *ValueConverter {
	newConverter := New(i)
	newConverter.converter.field = fieldName
	newConverter.converter.opts = c.opts
	newConverter.converter.storage = c.storage
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

// InstanceSetValues saves multiple key-value pairs.
func (c *converter) InstanceSetValues(m map[string]interface{}) {
	for k, v := range m {
		c.storage[k] = v
	}
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
	outT := reflect.ValueOf(out).Type()
	if outT.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}

	for outT.Kind() == reflect.Ptr {
		outT = outT.Elem()
	}

	switch outT.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return c.Int().Convert(out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return c.Uint().Convert(out)
	case reflect.Float32, reflect.Float64:
		return c.Float().Convert(out)
	case reflect.Bool:
		return c.Bool().Convert(out)
	case reflect.String:
		return c.String().Convert(out)
	case reflect.Array, reflect.Slice:
		return c.Slice().Convert(out)
	case reflect.Map, reflect.Struct:
		t := reflect.ValueOf(c.value).Type()
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Map {
			return c.Map().Convert(out)
		}
		return c.Struct().Convert(out)
	case reflect.Interface:
		outV := reflect.ValueOf(out)
		for outV.Kind() == reflect.Ptr {
			outV = outV.Elem()
		}
		outV.Set(reflect.ValueOf(c.value))
	default:
		var srcType reflect.Type
		if reflect.ValueOf(c.value).IsValid() {
			srcType = reflect.ValueOf(c.value).Type()
		}
		return &ConvertError{
			Field:   c.field,
			SrcType: srcType,
			DstType: reflect.ValueOf(out).Type(),
			Value:   c.value,
			Err:     ErrUnsupportedType,
		}
	}
	return nil
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

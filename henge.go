// Package henge is a type conversion library for Golang.
//
// 変化 (Henge) means "Appearing with a different figure." in Japanese.
// Henge as the name implies can easily convert to different types.
//
// Links
//
// Source code: https://github.com/soranoba/henge
//
package henge

import (
	"reflect"
)

type (
	// InstanceStore is an interface for Converter holds some key-value pairs.
	InstanceStore interface {
		// Get returns the value saved using Set.
		Get(key string) interface{}
		// Set saves the value on the key.
		Set(key string, value interface{})
		// SetValues saves multiple key-value pairs.
		SetValues(m map[string]interface{})
	}

	// Converter is an interface that has common functions of all Converters.
	Converter interface {
		InstanceStore
		Interface() interface{}
		Error() error
	}

	// baseConverter is a struct inherited by each Converter and has common functions.
	baseConverter struct {
		isNil   bool
		field   string
		opts    ConverterOpts
		storage map[string]interface{}
	}
)

func (c *baseConverter) new(i interface{}, fieldName string) *ValueConverter {
	newConverter := New(i)
	newConverter.baseConverter.field = fieldName
	newConverter.baseConverter.opts = c.opts
	newConverter.baseConverter.storage = c.storage
	return newConverter
}

// Get returns the value saved using Set.
func (c *baseConverter) Get(key string) interface{} {
	return c.storage[key]
}

// Set saves the value on the key.
func (c *baseConverter) Set(key string, value interface{}) {
	c.storage[key] = value
}

// SetValues saves multiple key-value pairs.
func (c *baseConverter) SetValues(m map[string]interface{}) {
	for k, v := range m {
		c.storage[k] = v
	}
}

func (c *baseConverter) wrapConvertError(src interface{}, dstType reflect.Type, err error) error {
	if convertErr, ok := err.(*ConvertError); ok {
		err := *convertErr
		err.DstType = dstType
		return &err
	}
	srcValue := src
	var srcType reflect.Type
	if reflect.ValueOf(srcValue).IsValid() {
		srcType = reflect.ValueOf(srcValue).Type()
		if c.isNil {
			srcValue = reflect.New(srcType).Elem().Interface()
		}
	}
	return &ConvertError{
		Field:   c.field,
		SrcType: srcType,
		DstType: dstType,
		Value:   srcValue,
		Err:     err,
	}
}

// ValueConverter is a converter that converts an interface type to another type.
type ValueConverter struct {
	*baseConverter
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
		baseConverter: &baseConverter{
			isNil:   isNil,
			opts:    opts,
			storage: map[string]interface{}{},
		},
		value: i,
		err:   nil,
	}
}

// Model converts the input to the specified type.
func (c *ValueConverter) Model(t interface{}) *ValueConverter {
	var (
		value interface{}
		err   error
	)

	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Ptr && v.Elem().CanSet() {
		value = t
		err = c.Convert(t)
	} else {
		v = reflect.New(v.Type())
		err = c.Convert(v.Interface())
		value = v.Elem().Interface()
	}
	return &ValueConverter{baseConverter: c.baseConverter, value: value, err: err}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *ValueConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Type().Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *ValueConverter) convert(outV reflect.Value) error {
	outT := outV.Type()
	for outT.Kind() == reflect.Ptr {
		outT = outT.Elem()
	}

	switch outT.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return c.Int().convert(outV)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return c.Uint().convert(outV)
	case reflect.Float32, reflect.Float64:
		return c.Float().convert(outV)
	case reflect.Bool:
		return c.Bool().convert(outV)
	case reflect.String:
		return c.String().convert(outV)
	case reflect.Array, reflect.Slice:
		return c.Slice().convert(outV)
	case reflect.Map, reflect.Struct:
		t := reflect.ValueOf(c.value).Type()
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Map {
			return c.Map().convert(outV)
		}
		return c.Struct().convert(outV)
	case reflect.Interface:
		for outV.Kind() == reflect.Ptr {
			outV = outV.Elem()
		}
		outV.Set(reflect.ValueOf(c.value))
		return nil
	default:
		return c.wrapConvertError(c.value, outV.Type(), ErrUnsupportedType)
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

// Interface returns the conversion result of interface type.
func (c *ValueConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *ValueConverter) Error() error {
	return c.err
}

func toInitializedNonPtrValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

package henge

import (
	"reflect"
	"strconv"
)

// String converts the input to string type.
func (c *ValueConverter) String() *StringConverter {
	var (
		value string
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	if inV.IsValid() {
		inT := inV.Type()
		outT := reflect.TypeOf(value)
		switch inV.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i = inV.Convert(reflect.TypeOf(i)).Interface().(int64)
			value = strconv.FormatInt(i, 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var u uint64
			u = inV.Convert(reflect.TypeOf(u)).Interface().(uint64)
			value = strconv.FormatUint(u, 10)
		case reflect.Float32, reflect.Float64:
			var f float64
			f = inV.Convert(reflect.TypeOf(f)).Interface().(float64)
			value = strconv.FormatFloat(f, c.opts.stringOpts.fmt, c.opts.stringOpts.prec, 64)
		case reflect.Bool:
			if inV.Interface().(bool) == true {
				value = "true"
			} else {
				value = "false"
			}
		default:
			if inT.ConvertibleTo(outT) {
				value = inV.Convert(outT).Interface().(string)
			} else {
				err = ErrUnsupportedType
			}
		}
	} else {
		err = ErrInvalidValue
	}

	if err != nil {
		err = &ConvertError{
			Field:   c.field,
			SrcType: reflect.ValueOf(c.value).Type(),
			DstType: reflect.ValueOf((*string)(nil)).Type().Elem(),
			Value:   c.value,
			Err:     err,
		}
	}
	return &StringConverter{converter: c.converter, value: value, err: err}
}

// StringConverter is a converter that converts a string type to another type.
type StringConverter struct {
	converter
	value string
	err   error
}

// Ptr converts the input to ptr type.
func (c *StringConverter) Ptr() *StringPtrConverter {
	if c.err != nil || c.isNil {
		return &StringPtrConverter{converter: c.converter, value: nil, err: c.err}
	}
	return &StringPtrConverter{converter: c.converter, value: &c.value, err: nil}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *StringConverter) Convert(out interface{}) error {
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
	case reflect.String:
		outV.Set(reflect.ValueOf(c.value))
	default:
		return c.new(c.value, c.field).Convert(out)
	}
	return nil
}

// Result returns the conversion result and error.
func (c *StringConverter) Result() (string, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *StringConverter) Value() string {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *StringConverter) Error() error {
	return c.err
}

// StringPtrConverter is a converter that converts a pointer of string type to another type.
type StringPtrConverter struct {
	converter
	value *string
	err   error
}

// Result returns the conversion result and error.
func (c *StringPtrConverter) Result() (*string, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *StringPtrConverter) Value() *string {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *StringPtrConverter) Error() error {
	return c.err
}

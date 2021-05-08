package henge

import (
	"reflect"
	"strconv"
)

type (
	// StringConverter is a converter that converts a string type to another type.
	StringConverter struct {
		*baseConverter
		value string
		err   error
	}

	// StringPtrConverter is a converter that converts a pointer of string type to another type.
	StringPtrConverter struct {
		*baseConverter
		value *string
		err   error
	}
)

// --------------------------------------------------------------------- //
// ValueConverter
// --------------------------------------------------------------------- //

// String converts the input to string type.
func (c *ValueConverter) String() *StringConverter {
	var (
		value string
		err   error
	)

	inV := reflect.Indirect(c.reflectValue)
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
		err = c.wrapConvertError(c.value, reflect.ValueOf((*string)(nil)).Type().Elem(), err)
	}
	return &StringConverter{baseConverter: c.baseConverter, value: value, err: err}
}

// StringPtr converts the input to pointer of string type.
func (c *ValueConverter) StringPtr() *StringPtrConverter {
	return c.String().Ptr()
}

// --------------------------------------------------------------------- //
// StringConverter
// --------------------------------------------------------------------- //

// Ptr converts the input to ptr type.
func (c *StringConverter) Ptr() *StringPtrConverter {
	if c.err != nil || c.isNil {
		return &StringPtrConverter{baseConverter: c.baseConverter, value: nil, err: c.err}
	}
	return &StringPtrConverter{baseConverter: c.baseConverter, value: &c.value, err: nil}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *StringConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *StringConverter) convert(outV reflect.Value) error {
	if c.err != nil {
		return c.wrapConvertError(c.value, outV.Type(), c.err)
	}
	if c.isNil {
		return nil
	}

	elemOutV := toInitializedNonPtrValue(outV)

	switch elemOutV.Kind() {
	case reflect.String:
		elemOutV.Set(reflect.ValueOf(c.value))
	default:
		return c.new(c.value, c.field).convert(outV)
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

// Interface returns the conversion result of interface type.
func (c *StringConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *StringConverter) Error() error {
	return c.err
}

// --------------------------------------------------------------------- //
// StringPtrConverter
// --------------------------------------------------------------------- //

// Result returns the conversion result and error.
func (c *StringPtrConverter) Result() (*string, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *StringPtrConverter) Value() *string {
	return c.value
}

// Interface returns the conversion result of interface type.
func (c *StringPtrConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *StringPtrConverter) Error() error {
	return c.err
}

package henge

import (
	"reflect"
	"strconv"
)

func (c *ValueConverter) String(opts ...func(*StringConverterOpt) *StringConverterOpt) *StringConverter {
	var (
		value string
		err   error
	)

	opt := &StringConverterOpt{
		fmt:  'f',
		prec: -1,
	}
	for _, f := range opts {
		opt = f(opt)
	}

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
			value = strconv.FormatFloat(f, opt.fmt, opt.prec, 64)
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
				err = unsupportedTypeErr
			}
		}
	} else {
		err = invalidValueErr
	}

	if err != nil {
		err = &ConvertError{
			Field:   "",
			SrcType: reflect.ValueOf(c.value).Type(),
			DstType: reflect.ValueOf((*string)(nil)).Type().Elem(),
			Err:     err,
		}
	}
	return &StringConverter{value: value, err: err}
}

type StringConverterOpt struct {
	fmt  byte
	prec int
}

func WithFloatFormat(fmt byte, prec int) func(*StringConverterOpt) *StringConverterOpt {
	return func(opt *StringConverterOpt) *StringConverterOpt {
		opt.fmt = fmt
		opt.prec = prec
		return opt
	}
}

type StringConverter struct {
	value string
	err   error
}

func (c *StringConverter) Ptr() *StringPtrConverter {
	if c.err != nil {
		return &StringPtrConverter{value: nil, err: c.err}
	}
	return &StringPtrConverter{value: &c.value, err: nil}
}

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
		return New(c.value).Convert(out)
	}
	return nil
}

func (c *StringConverter) Result() (string, error) {
	return c.value, c.err
}

func (c *StringConverter) Value() string {
	return c.value
}

func (c *StringConverter) Error() error {
	return c.err
}

type StringPtrConverter struct {
	value *string
	err   error
}

func (c *StringPtrConverter) Result() (*string, error) {
	return c.value, c.err
}

func (c *StringPtrConverter) Value() *string {
	return c.value
}

func (c *StringPtrConverter) Error() error {
	return c.err
}

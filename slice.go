package henge

import "reflect"

// Slice converts the input to slice type.
func (c *ValueConverter) Slice() *SliceConverter {
	var (
		value []interface{}
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	switch inV.Kind() {
	case reflect.Array, reflect.Slice:
		value = make([]interface{}, inV.Len())
		for i := 0; i < inV.Len(); i++ {
			value[i] = inV.Index(i).Interface()
		}
	default:
		err = ErrUnsupportedType
	}

	if c.isNil {
		return &SliceConverter{converter: c.converter, value: nil, err: err}
	}
	return &SliceConverter{converter: c.converter, value: value, err: err}
}

// StringSlice converts the input to slice of string type.
func (c *ValueConverter) StringSlice() *StringSliceConverter {
	return c.Slice().StringSlice()
}

// IntSlice converts the input to slice of int type.
func (c *ValueConverter) IntSlice() *IntegerSliceConverter {
	return c.Slice().IntSlice()
}

// UintSlice converts the input to slice of uint type.
func (c *ValueConverter) UintSlice() *UnsignedIntegerSliceConverter {
	return c.Slice().UintSlice()
}

// FloatSlice converts the input to slice of fload type.
func (c *ValueConverter) FloatSlice() *FloatSliceConverter {
	return c.Slice().FloatSlice()
}

// SliceConverter is a converter that converts a slice type to another type.
type SliceConverter struct {
	converter
	value []interface{}
	err   error
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *SliceConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}

	for outV.Kind() == reflect.Ptr {
		if outV.IsNil() {
			outV.Set(reflect.New(outV.Type().Elem()))
		}
		outV = outV.Elem()
	}

	unsupportedTypeErr := &ConvertError{
		Field:   c.field,
		SrcType: reflect.ValueOf(c.value).Type(),
		DstType: outV.Type(),
		Value:   c.value,
		Err:     ErrUnsupportedType,
	}

	switch outV.Kind() {
	case reflect.Array:
		inV := reflect.Indirect(reflect.ValueOf(c.value))
		if inV.Kind() != reflect.Array && inV.Kind() != reflect.Slice {
			return unsupportedTypeErr
		}

		v := reflect.New(reflect.ArrayOf(outV.Len(), outV.Type().Elem())).Elem()
		for i := 0; i < inV.Len() && i < outV.Len(); i++ {
			elem := reflect.New(outV.Type().Elem())
			fieldName := c.field + "[" + New(i).String().Value() + "]"
			if err := c.new(inV.Index(i).Interface(), fieldName).Convert(elem.Interface()); err != nil {
				return err
			}
			v.Index(i).Set(elem.Elem())
		}
		outV.Set(v)
	case reflect.Slice:
		inV := reflect.Indirect(reflect.ValueOf(c.value))
		if inV.Kind() != reflect.Array && inV.Kind() != reflect.Slice {
			return unsupportedTypeErr
		}

		v := reflect.MakeSlice(reflect.SliceOf(outV.Type().Elem()), inV.Len(), inV.Len())
		for i := 0; i < inV.Len(); i++ {
			elem := reflect.New(outV.Type().Elem())
			fieldName := c.field + "[" + New(i).String().Value() + "]"
			if err := c.new(inV.Index(i).Interface(), fieldName).Convert(elem.Interface()); err != nil {
				return err
			}
			v.Index(i).Set(elem.Elem())
		}
		outV.Set(v)
	default:
		return unsupportedTypeErr
	}
	return nil
}

// IntSlice converts the input to slice of int type.
func (c *SliceConverter) IntSlice() *IntegerSliceConverter {
	var (
		value []int64 = make([]int64, len(c.value))
		err   error
	)

	for i, v := range c.value {
		fieldName := c.field + "[" + New(i).String().Value() + "]"
		value[i], err = c.new(v, fieldName).Int().Result()
		if err != nil {
			return &IntegerSliceConverter{converter: c.converter, value: nil, err: err}
		}
	}
	return &IntegerSliceConverter{converter: c.converter, value: value, err: nil}
}

// UintSlice converts the input to slice of uint type.
func (c *SliceConverter) UintSlice() *UnsignedIntegerSliceConverter {
	var (
		value []uint64 = make([]uint64, len(c.value))
		err   error
	)

	for i, v := range c.value {
		fieldName := c.field + "[" + New(i).String().Value() + "]"
		value[i], err = c.new(v, fieldName).Uint().Result()
		if err != nil {
			return &UnsignedIntegerSliceConverter{converter: c.converter, value: nil, err: err}
		}
	}
	return &UnsignedIntegerSliceConverter{converter: c.converter, value: value, err: nil}
}

// FloatSlice converts the input to slice of float type.
func (c *SliceConverter) FloatSlice() *FloatSliceConverter {
	var (
		value []float64 = make([]float64, len(c.value))
		err   error
	)

	for i, v := range c.value {
		fieldName := c.field + "[" + New(i).String().Value() + "]"
		value[i], err = c.new(v, fieldName).Float().Result()
		if err != nil {
			return &FloatSliceConverter{converter: c.converter, value: nil, err: err}
		}
	}
	return &FloatSliceConverter{converter: c.converter, value: value, err: nil}
}

// StringSlice converts the input to slice of string type.
func (c *SliceConverter) StringSlice() *StringSliceConverter {
	var (
		value []string = make([]string, len(c.value))
		err   error
	)

	for i, v := range c.value {
		fieldName := c.field + "[" + New(i).String().Value() + "]"
		value[i], err = c.new(v, fieldName).String().Result()
		if err != nil {
			return &StringSliceConverter{converter: c.converter, value: nil, err: err}
		}
	}
	return &StringSliceConverter{converter: c.converter, value: value, err: nil}
}

// Result returns the conversion result and error
func (c *SliceConverter) Result() ([]interface{}, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *SliceConverter) Value() []interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *SliceConverter) Error() error {
	return c.err
}

// IntegerSliceConverter is a converter that converts a slice of integer type to another type.
type IntegerSliceConverter struct {
	converter
	value []int64
	err   error
}

// Result returns the conversion result and error
func (c *IntegerSliceConverter) Result() ([]int64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *IntegerSliceConverter) Value() []int64 {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *IntegerSliceConverter) Error() error {
	return c.err
}

// UnsignedIntegerSliceConverter is a converter that converts a slice of uint type to another type.
type UnsignedIntegerSliceConverter struct {
	converter
	value []uint64
	err   error
}

// Result returns the conversion result and error.
func (c *UnsignedIntegerSliceConverter) Result() ([]uint64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *UnsignedIntegerSliceConverter) Value() []uint64 {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *UnsignedIntegerSliceConverter) Error() error {
	return c.err
}

// FloatSliceConverter is a converter that converts a slice of float type to another type.
type FloatSliceConverter struct {
	converter
	value []float64
	err   error
}

// Result returns the conversion result and error.
func (c *FloatSliceConverter) Result() ([]float64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *FloatSliceConverter) Value() []float64 {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *FloatSliceConverter) Error() error {
	return c.err
}

// StringSliceConverter is a converter that converts a slice of string type to another type.
type StringSliceConverter struct {
	converter
	value []string
	err   error
}

// Result returns the conversion result and error.
func (c *StringSliceConverter) Result() ([]string, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *StringSliceConverter) Value() []string {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *StringSliceConverter) Error() error {
	return c.err
}

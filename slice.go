package henge

import "reflect"

type (
	// SliceConverter is a converter that converts a slice type to another type.
	SliceConverter struct {
		*baseConverter
		value []interface{}
		err   error
	}

	// IntegerSliceConverter is a converter that converts a slice of integer type to another type.
	IntegerSliceConverter struct {
		*baseConverter
		value []int64
		err   error
	}

	// UnsignedIntegerSliceConverter is a converter that converts a slice of uint type to another type.
	UnsignedIntegerSliceConverter struct {
		*baseConverter
		value []uint64
		err   error
	}

	// FloatSliceConverter is a converter that converts a slice of float type to another type.
	FloatSliceConverter struct {
		*baseConverter
		value []float64
		err   error
	}

	// StringSliceConverter is a converter that converts a slice of string type to another type.
	StringSliceConverter struct {
		*baseConverter
		value []string
		err   error
	}
)

// --------------------------------------------------------------------- //
// ValueConverter
// --------------------------------------------------------------------- //

// Slice converts the input to slice type.
func (c *ValueConverter) Slice() *SliceConverter {
	var (
		value []interface{}
		err   error
	)

	inV := reflect.Indirect(c.reflectValue)
	switch inV.Kind() {
	case reflect.Array, reflect.Slice:
		value = make([]interface{}, inV.Len())
		for i := 0; i < inV.Len(); i++ {
			vConv := c.opts.sliceOpts.valueConversionFunc(c.new(inV.Index(i).Interface(), c.field+"["+New(i).String().Value()+"]"))
			if err = vConv.Error(); err != nil {
				break
			}
			value[i] = vConv.Interface()
		}
	default:
		err = ErrUnsupportedType
	}

	if err != nil {
		err = c.wrapConvertError(c.value, reflect.ValueOf(value).Type(), err)
	}

	if c.isNil {
		return &SliceConverter{baseConverter: c.baseConverter, value: nil, err: err}
	}
	return &SliceConverter{baseConverter: c.baseConverter, value: value, err: err}
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

// --------------------------------------------------------------------- //
// SliceConverter
// --------------------------------------------------------------------- //

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *SliceConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *SliceConverter) convert(outV reflect.Value) error {
	if c.err != nil {
		return c.wrapConvertError(c.value, outV.Type(), c.err)
	}
	if c.isNil {
		return nil
	}

	elemOutV := toInitializedNonPtrValue(outV)
	unsupportedTypeErr := c.wrapConvertError(c.value, outV.Type(), ErrUnsupportedType)

	switch elemOutV.Kind() {
	case reflect.Array:
		inV := reflect.Indirect(reflect.ValueOf(c.value))
		if inV.Kind() != reflect.Array && inV.Kind() != reflect.Slice {
			return unsupportedTypeErr
		}

		v := reflect.New(reflect.ArrayOf(elemOutV.Len(), elemOutV.Type().Elem())).Elem()
		for i := 0; i < inV.Len() && i < elemOutV.Len(); i++ {
			elem := reflect.New(elemOutV.Type().Elem()).Elem()
			fieldName := c.field + "[" + New(i).String().Value() + "]"
			if err := c.new(inV.Index(i).Interface(), fieldName).convert(elem); err != nil {
				return err
			}
			v.Index(i).Set(elem)
		}
		elemOutV.Set(v)
	case reflect.Slice:
		inV := reflect.Indirect(reflect.ValueOf(c.value))
		if inV.Kind() != reflect.Array && inV.Kind() != reflect.Slice {
			return unsupportedTypeErr
		}

		v := reflect.MakeSlice(reflect.SliceOf(elemOutV.Type().Elem()), inV.Len(), inV.Len())
		for i := 0; i < inV.Len(); i++ {
			elem := reflect.New(elemOutV.Type().Elem()).Elem()
			fieldName := c.field + "[" + New(i).String().Value() + "]"
			if err := c.new(inV.Index(i).Interface(), fieldName).convert(elem); err != nil {
				return err
			}
			v.Index(i).Set(elem)
		}
		elemOutV.Set(v)
	default:
		return unsupportedTypeErr
	}
	return nil
}

// IntSlice converts the input to slice of int type.
func (c *SliceConverter) IntSlice() *IntegerSliceConverter {
	var value []int64
	if err := c.Convert(&value); err != nil {
		return &IntegerSliceConverter{baseConverter: c.baseConverter, value: value, err: err}
	}
	return &IntegerSliceConverter{baseConverter: c.baseConverter, value: value, err: nil}
}

// UintSlice converts the input to slice of uint type.
func (c *SliceConverter) UintSlice() *UnsignedIntegerSliceConverter {
	var value []uint64
	if err := c.Convert(&value); err != nil {
		return &UnsignedIntegerSliceConverter{baseConverter: c.baseConverter, value: value, err: err}
	}
	return &UnsignedIntegerSliceConverter{baseConverter: c.baseConverter, value: value, err: nil}
}

// FloatSlice converts the input to slice of float type.
func (c *SliceConverter) FloatSlice() *FloatSliceConverter {
	var value []float64
	if err := c.Convert(&value); err != nil {
		return &FloatSliceConverter{baseConverter: c.baseConverter, value: value, err: err}
	}
	return &FloatSliceConverter{baseConverter: c.baseConverter, value: value, err: nil}
}

// StringSlice converts the input to slice of string type.
func (c *SliceConverter) StringSlice() *StringSliceConverter {
	var value []string
	if err := c.Convert(&value); err != nil {
		return &StringSliceConverter{baseConverter: c.baseConverter, value: value, err: err}
	}
	return &StringSliceConverter{baseConverter: c.baseConverter, value: value, err: nil}
}

// Result returns the conversion result and error
func (c *SliceConverter) Result() ([]interface{}, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *SliceConverter) Value() []interface{} {
	return c.value
}

// Interface returns the conversion result of interface type.
func (c *SliceConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *SliceConverter) Error() error {
	return c.err
}

// --------------------------------------------------------------------- //
// IntegerSliceConverter
// --------------------------------------------------------------------- //

// Result returns the conversion result and error
func (c *IntegerSliceConverter) Result() ([]int64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *IntegerSliceConverter) Value() []int64 {
	return c.value
}

// Interface returns the conversion result of interface type.
func (c *IntegerSliceConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *IntegerSliceConverter) Error() error {
	return c.err
}

// --------------------------------------------------------------------- //
// UnsignedIntegerSliceConverter
// --------------------------------------------------------------------- //

// Result returns the conversion result and error.
func (c *UnsignedIntegerSliceConverter) Result() ([]uint64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *UnsignedIntegerSliceConverter) Value() []uint64 {
	return c.value
}

// Interface returns the conversion result of interface type.
func (c *UnsignedIntegerSliceConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *UnsignedIntegerSliceConverter) Error() error {
	return c.err
}

// --------------------------------------------------------------------- //
// FloatSliceConverter
// --------------------------------------------------------------------- //

// Result returns the conversion result and error.
func (c *FloatSliceConverter) Result() ([]float64, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *FloatSliceConverter) Value() []float64 {
	return c.value
}

// Interface returns the conversion result of interface type.
func (c *FloatSliceConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *FloatSliceConverter) Error() error {
	return c.err
}

// --------------------------------------------------------------------- //
// StringSliceConverter
// --------------------------------------------------------------------- //

// Result returns the conversion result and error.
func (c *StringSliceConverter) Result() ([]string, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *StringSliceConverter) Value() []string {
	return c.value
}

// Interface returns the conversion result of interface type.
func (c *StringSliceConverter) Interface() interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *StringSliceConverter) Error() error {
	return c.err
}

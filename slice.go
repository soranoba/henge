package henge

import "reflect"

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
		err = unsupportedTypeErr
	}
	return &SliceConverter{value: value, err: err}
}

type SliceConverter struct {
	value []interface{}
	err   error
}

func (c *SliceConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}

	for outV.Kind() == reflect.Ptr {
		outV = outV.Elem()
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
			if err := New(inV.Index(i).Interface()).Convert(elem.Interface()); err != nil {
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
			if err := New(inV.Index(i).Interface()).Convert(elem.Interface()); err != nil {
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

func (c *SliceConverter) IntSlice() *IntegerSliceConverter {
	var (
		value []int64 = make([]int64, len(c.value))
		err   error
	)

	for i, v := range c.value {
		value[i], err = New(v).Int().Result()
		if err != nil {
			return &IntegerSliceConverter{value: nil, err: err}
		}
	}
	return &IntegerSliceConverter{value: value, err: nil}
}

func (c *SliceConverter) UintSlice() *UnsignedIntegerSliceConverter {
	var (
		value []uint64 = make([]uint64, len(c.value))
		err   error
	)

	for i, v := range c.value {
		value[i], err = New(v).Uint().Result()
		if err != nil {
			return &UnsignedIntegerSliceConverter{value: nil, err: err}
		}
	}
	return &UnsignedIntegerSliceConverter{value: value, err: nil}
}

func (c *SliceConverter) FloatSlice() *FloatSliceConverter {
	var (
		value []float64 = make([]float64, len(c.value))
		err   error
	)

	for i, v := range c.value {
		value[i], err = New(v).Float().Result()
		if err != nil {
			return &FloatSliceConverter{value: nil, err: err}
		}
	}
	return &FloatSliceConverter{value: value, err: nil}
}

func (c *SliceConverter) StringSlice() *StringSliceConverter {
	var (
		value []string = make([]string, len(c.value))
		err   error
	)

	for i, v := range c.value {
		value[i], err = New(v).String().Result()
		if err != nil {
			return &StringSliceConverter{value: nil, err: err}
		}
	}
	return &StringSliceConverter{value: value, err: nil}
}

func (c *SliceConverter) Result() ([]interface{}, error) {
	return c.value, c.err
}

func (c *SliceConverter) Value() []interface{} {
	return c.value
}

func (c *SliceConverter) Error() error {
	return c.err
}

type IntegerSliceConverter struct {
	value []int64
	err   error
}

func (c *IntegerSliceConverter) Result() ([]int64, error) {
	return c.value, c.err
}

func (c *IntegerSliceConverter) Value() []int64 {
	return c.value
}

func (c *IntegerSliceConverter) Error() error {
	return c.err
}

type UnsignedIntegerSliceConverter struct {
	value []uint64
	err   error
}

func (c *UnsignedIntegerSliceConverter) Result() ([]uint64, error) {
	return c.value, c.err
}

func (c *UnsignedIntegerSliceConverter) Value() []uint64 {
	return c.value
}

func (c *UnsignedIntegerSliceConverter) Error() error {
	return c.err
}

type FloatSliceConverter struct {
	value []float64
	err   error
}

func (c *FloatSliceConverter) Result() ([]float64, error) {
	return c.value, c.err
}

func (c *FloatSliceConverter) Value() []float64 {
	return c.value
}

func (c *FloatSliceConverter) Error() error {
	return c.err
}

type StringSliceConverter struct {
	value []string
	err   error
}

func (c *StringSliceConverter) Result() ([]string, error) {
	return c.value, c.err
}

func (c *StringSliceConverter) Value() []string {
	return c.value
}

func (c *StringSliceConverter) Error() error {
	return c.err
}

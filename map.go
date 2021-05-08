package henge

import "reflect"

type (
	// MapConverter is a converter that converts a map type to another type.
	MapConverter struct {
		*baseConverter
		value reflect.Value
		err   error
	}
)

// Map converts the input to map type.
func (c *ValueConverter) Map() *MapConverter {
	return c.mapWithDepth(0)
}

func (c *ValueConverter) makeOutputMapVar() reflect.Value {
	return reflect.New(reflect.MapOf(c.opts.mapOpts.keyType, interfaceType)).Elem()
}

func (c *ValueConverter) mapWithDepth(depth uint) *MapConverter {
	var (
		value = c.makeOutputMapVar()
		err   error
	)

	convAndSet := func(kVal reflect.Value, vVal reflect.Value) {
		if !kVal.CanInterface() || !vVal.CanInterface() {
			return
		}
		if !c.opts.mapOpts.filterFuns.All(kVal.Interface(), vVal.Interface()) {
			return
		}
		strKey := New(kVal.Interface()).String().Value()
		kConv := c.opts.mapOpts.keyConversionFunc(c.new(kVal.Interface(), c.field+"[]"+strKey))
		if err = kConv.Error(); err != nil {
			return
		}

		// NOTE: `reflect.ValueOf(nil).Convert(type)` occurs panic.
		//       Therefore, it uses ptr and Elem after store it on a variable once.
		k := kConv.Interface()
		convertedKeyVal := (func() reflect.Value {
			if k == nil {
				return reflect.ValueOf(&k).Elem().Convert(c.opts.keyType)
			} else {
				return reflect.ValueOf(k).Convert(c.opts.keyType)
			}
		})()

		switch reflect.Indirect(reflect.ValueOf(vVal.Interface())).Kind() {
		case reflect.Struct, reflect.Map:
			if depth < c.opts.mapOpts.maxDepth {
				conv := c.new(vVal.Interface(), c.field+"."+strKey).mapWithDepth(depth + 1)
				if err = conv.err; err != nil {
					return
				}
				value.SetMapIndex(convertedKeyVal, conv.value)
				break
			}
			fallthrough
		default:
			vConv := c.opts.mapOpts.valueConversionFunc(c.new(vVal.Interface(), c.field+"."+strKey))
			if err = vConv.Error(); err != nil {
				return
			}
			v := vConv.Interface()
			value.SetMapIndex(convertedKeyVal, reflect.ValueOf(&v).Elem())
		}
	}

	inV := reflect.Indirect(c.reflectValue)
	switch inV.Kind() {
	case reflect.Map:
		value = reflect.MakeMap(value.Type())
		iter := inV.MapRange()
		for iter.Next() {
			convAndSet(iter.Key(), iter.Value())
			if err != nil {
				break
			}
		}
	case reflect.Struct:
		value = reflect.MakeMap(value.Type())
		for i := 0; i < inV.NumField(); i++ {
			convAndSet(reflect.ValueOf(inV.Type().Field(i).Name), inV.Field(i))
			if err != nil {
				break
			}
		}
	default:
		err = ErrUnsupportedType
	}

	if err != nil {
		err = c.wrapConvertError(c.value, value.Type(), err)
	}

	return &MapConverter{baseConverter: c.baseConverter, value: value, err: err}
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *MapConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Type().Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *MapConverter) convert(outV reflect.Value) error {
	if c.err != nil {
		return c.wrapConvertError(c.value.Interface(), outV.Type(), c.err)
	}
	if c.isNil {
		return nil
	}

	for outV.Kind() == reflect.Ptr {
		if outV.IsNil() {
			outV.Set(reflect.New(outV.Type().Elem()))
		}
		outV = outV.Elem()
	}

	switch outV.Kind() {
	case reflect.Map:
		if outV.IsNil() {
			outV.Set(reflect.MakeMap(outV.Type()))
		}
		iter := c.value.MapRange()
		for iter.Next() {
			keyV := reflect.New(outV.Type().Key())
			valueV := reflect.New(outV.Type().Elem())
			strKey := New(keyV).String().Value()
			if err := c.new(iter.Key().Interface(), c.field+"[]"+strKey).convert(keyV); err != nil {
				return err
			}
			if err := c.new(iter.Value().Interface(), c.field+"["+strKey+"]").convert(valueV); err != nil {
				return err
			}
			outV.SetMapIndex(keyV.Elem(), valueV.Elem())
		}
	case reflect.Struct:
		m := map[string]interface{}{}
		iter := c.value.MapRange()
		for iter.Next() {
			strKey, err := c.new(iter.Key().Interface(), c.field+"[]").String().Result()
			if err != nil {
				return err
			}
			m[strKey] = iter.Value().Interface()
		}

		for _, outField := range getStructFields(outV.Type()) {
			if outField.isIgnore() {
				continue
			}

			if value, ok := m[outField.name]; ok {
				// NOTE: initialized embedded field.
				anchor := outV
				for _, index := range outField.index {
					v := anchor.Field(index)
					if v.Kind() == reflect.Ptr {
						if v.IsNil() {
							v.Set(reflect.New(v.Type().Elem()))
						}
						anchor = v.Elem()
					}
				}

				target := outV.FieldByIndex(outField.index)
				if err := c.new(value, c.field+"."+outField.name).convert(target); err != nil {
					return err
				}
			}
		}
	default:
		return c.wrapConvertError(c.value.Interface(), outV.Type(), ErrUnsupportedType)
	}
	return nil
}

// Result returns the conversion result and error.
func (c *MapConverter) Result() (map[interface{}]interface{}, error) {
	if c.isNil {
		return nil, c.err
	}
	return c.value.Interface().(map[interface{}]interface{}), c.err
}

// Value returns the conversion result.
func (c *MapConverter) Value() map[interface{}]interface{} {
	if c.isNil {
		return nil
	}
	return c.value.Interface().(map[interface{}]interface{})
}

// Interface returns the conversion result of interface type.
func (c *MapConverter) Interface() interface{} {
	if c.isNil {
		return nil
	}
	return c.value
}

// Error returns an error if the conversion fails.
func (c *MapConverter) Error() error {
	return c.err
}

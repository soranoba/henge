package henge

import "reflect"

// Map converts the input to map type.
func (c *ValueConverter) Map() *MapConverter {
	return c.mapWithDepth(0)
}

// JsonMap converts the input to json map type.
func (c *ValueConverter) JsonMap() *JsonMapConverter {
	return c.jsonMapWithDepth(0)
}

func (c *ValueConverter) mapWithDepth(depth uint) *MapConverter {
	var (
		value map[interface{}]interface{}
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	switch inV.Kind() {
	case reflect.Map:
		value = map[interface{}]interface{}{}
		iter := inV.MapRange()
		for iter.Next() {
			iterK := iter.Key()
			iterV := iter.Value()
			if c.opts.mapOpts.filterFunc != nil && !c.opts.mapOpts.filterFunc(iterK.Interface(), iterV.Interface()) {
				continue
			}
			if reflect.Indirect(iterV).Kind() == reflect.Struct && depth < c.opts.mapOpts.maxDepth {
				strKey := New(iterK.Interface()).String().Value()
				var v interface{}
				if v, err = c.new(iterV.Interface(), c.field+"."+strKey).mapWithDepth(depth + 1).Result(); err != nil {
					break
				}
				value[iterK.Interface()] = v
			} else {
				value[iterK.Interface()] = iterV.Interface()
			}
		}
	case reflect.Struct:
		value = map[interface{}]interface{}{}
		for i := 0; i < inV.NumField(); i++ {
			field := inV.Field(i)
			if !field.CanInterface() {
				continue
			}

			key := inV.Type().Field(i).Name
			if c.opts.mapOpts.filterFunc != nil && !c.opts.mapOpts.filterFunc(key, field.Interface()) {
				continue
			}
			if reflect.Indirect(field).Kind() == reflect.Struct && depth < c.opts.mapOpts.maxDepth {
				var v interface{}
				if v, err = c.new(field.Interface(), c.field+"."+key).mapWithDepth(depth + 1).Result(); err != nil {
					break
				}
				value[key] = v
			} else {
				value[key] = field.Interface()
			}
		}
	default:
		err = ErrUnsupportedType
	}

	if err != nil {
		var srcType reflect.Type
		if reflect.ValueOf(c.value).IsValid() {
			srcType = reflect.ValueOf(c.value).Type()
		}
		err = &ConvertError{
			Field:   c.field,
			SrcType: srcType,
			DstType: reflect.ValueOf(value).Type(),
			Value:   c.value,
			Err:     err,
		}
	}

	if c.isNil {
		return &MapConverter{converter: c.converter, value: nil, err: err}
	}
	return &MapConverter{converter: c.converter, value: value, err: err}
}

func (c *ValueConverter) jsonMapWithDepth(depth uint) *JsonMapConverter {
	var (
		value map[string]interface{}
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	switch inV.Kind() {
	case reflect.Map:
		value = map[string]interface{}{}
		iter := inV.MapRange()
		for iter.Next() {
			iterK := iter.Key()
			iterV := iter.Value()
			if c.opts.mapOpts.filterFunc != nil && !c.opts.mapOpts.filterFunc(iterK.Interface(), iterV.Interface()) {
				continue
			}
			strKey := New(iter.Key().Interface()).String().Value()
			if reflect.Indirect(iterV).Kind() == reflect.Struct && depth < c.opts.mapOpts.maxDepth {
				var v interface{}
				if v, err = c.new(iterV.Interface(), c.field+"."+strKey).jsonMapWithDepth(depth + 1).Result(); err != nil {
					break
				}
				value[strKey] = v
			} else {
				value[strKey] = iterV.Interface()
			}
		}
	case reflect.Struct:
		value = map[string]interface{}{}
		for i := 0; i < inV.NumField(); i++ {
			field := inV.Field(i)
			if !field.CanInterface() {
				continue
			}

			key := inV.Type().Field(i).Name
			if c.opts.mapOpts.filterFunc != nil && !c.opts.mapOpts.filterFunc(key, field.Interface()) {
				continue
			}
			if reflect.Indirect(field).Kind() == reflect.Struct && depth < c.opts.mapOpts.maxDepth {
				var v interface{}
				if v, err = c.new(field.Interface(), c.field+"."+key).jsonMapWithDepth(depth + 1).Result(); err != nil {
					break
				}
				value[key] = v
			} else {
				value[key] = field.Interface()
			}
		}
	default:
		err = ErrUnsupportedType
	}

	if err != nil {
		var srcType reflect.Type
		if reflect.ValueOf(c.value).IsValid() {
			srcType = reflect.ValueOf(c.value).Type()
		}
		err = &ConvertError{
			Field:   c.field,
			SrcType: srcType,
			DstType: reflect.ValueOf(value).Type(),
			Value:   c.value,
			Err:     err,
		}
	}

	if c.isNil {
		return &JsonMapConverter{converter: c.converter, value: nil, err: err}
	}
	return &JsonMapConverter{converter: c.converter, value: value, err: err}
}

// MapConverter is a converter that converts a map type to another type.
type MapConverter struct {
	converter
	value map[interface{}]interface{}
	err   error
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *MapConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}

	if c.err != nil {
		return c.err
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
		for k, v := range c.value {
			keyV := reflect.New(outV.Type().Key())
			valueV := reflect.New(outV.Type().Elem())
			strKey := New(keyV).String().Value()
			if err := c.new(k, c.field+"[]"+strKey).Convert(keyV.Interface()); err != nil {
				return err
			}
			if err := c.new(v, c.field+"["+strKey+"]").Convert(valueV.Interface()); err != nil {
				return err
			}
			outV.SetMapIndex(keyV.Elem(), valueV.Elem())
		}
	case reflect.Struct:
		m := map[string]interface{}{}
		for k, v := range c.value {
			strKey, err := c.new(k, c.field+"[]").String().Result()
			if err != nil {
				return err
			}
			m[strKey] = v
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
				if target.Kind() != reflect.Ptr {
					target = target.Addr()
				}
				if err := c.new(value, c.field+"."+outField.name).Convert(target.Interface()); err != nil {
					return err
				}
			}
		}
	default:
		return ErrUnsupportedType
	}
	return nil
}

// Result returns the conversion result and error.
func (c *MapConverter) Result() (map[interface{}]interface{}, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *MapConverter) Value() map[interface{}]interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *MapConverter) Error() error {
	return c.err
}

// JsonMapConverter is a converter that converts a json map type to another type.
type JsonMapConverter struct {
	converter
	value map[string]interface{}
	err   error
}

// Result returns the conversion result and error.
func (c *JsonMapConverter) Result() (map[string]interface{}, error) {
	return c.value, c.err
}

// Value returns the conversion result.
func (c *JsonMapConverter) Value() map[string]interface{} {
	return c.value
}

// Error returns an error if the conversion fails.
func (c *JsonMapConverter) Error() error {
	return c.err
}

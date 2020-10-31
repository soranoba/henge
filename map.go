package henge

import "reflect"

func (c *ValueConverter) Map() *MapConverter {
	return c.mapWithDepth(0)
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
			iterV := iter.Value()
			if reflect.Indirect(iterV).Kind() == reflect.Struct && depth < c.opts.mapOpts.maxDepth {
				strKey := New(iter.Key().Interface()).String().Value()
				var v interface{}
				if v, err = c.new(iterV.Interface(), c.field+"."+strKey).mapWithDepth(depth + 1).Result(); err != nil {
					break
				}
				value[iter.Key().Interface()] = v
			} else {
				value[iter.Key().Interface()] = iterV.Interface()
			}
		}
	case reflect.Struct:
		value = map[interface{}]interface{}{}
		for i := 0; i < inV.NumField(); i++ {
			field := inV.Field(i)
			if reflect.Indirect(field).Kind() == reflect.Struct && depth < c.opts.mapOpts.maxDepth {
				key := inV.Type().Field(i).Name
				var v interface{}
				if v, err = c.new(field.Interface(), c.field+"."+key).mapWithDepth(depth + 1).Result(); err != nil {
					break
				}
				value[key] = v
			} else {
				value[inV.Type().Field(i).Name] = inV.Field(i).Interface()
			}
		}
	default:
		err = unsupportedTypeErr
	}
	return &MapConverter{converter: c.converter, value: value, err: err}
}

type MapConverter struct {
	converter
	value map[interface{}]interface{}
	err   error
}

func (c *MapConverter) Convert(out interface{}) error {
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}

	for outV.Kind() == reflect.Ptr {
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
				// NOTE: initialized embeded field.
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
		return unsupportedTypeErr
	}
	return nil
}

func (c *MapConverter) Result() (map[interface{}]interface{}, error) {
	return c.value, c.err
}

func (c *MapConverter) Value() map[interface{}]interface{} {
	return c.value
}

func (c *MapConverter) Error() error {
	return c.err
}
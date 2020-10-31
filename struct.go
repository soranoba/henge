package henge

import (
	"reflect"
)

func (c *ValueConverter) Struct() *StructConverter {
	var (
		value interface{}
		err   error
	)

	inV := reflect.Indirect(reflect.ValueOf(c.value))
	switch inV.Kind() {
	case reflect.Struct:
		value = c.value
	default:
		err = unsupportedTypeErr
	}
	return &StructConverter{converter: c.converter, value: value, err: err}
}

type StructConverter struct {
	converter
	value interface{}
	err   error
}

func (c *StructConverter) Convert(out interface{}) error {
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
		outV = outV.Elem()
	}

	switch outV.Kind() {
	case reflect.Struct:
		inV := reflect.Indirect(reflect.ValueOf(c.value))

		// NOTE: Types that are simply converted (it also copies private fields)
		if inV.Type().ConvertibleTo(outV.Type()) {
			outV.Set(inV.Convert(outV.Type()))
			return nil
		}

		inFields := getStructFields(inV.Type())
		for _, outField := range getStructFields(outV.Type()) {
			if outField.isIgnore() {
				continue
			}

			if f, ok := inV.Type().FieldByName(outField.name); ok {
				inField := func() structField {
					for _, field := range inFields {
						if field.isMatch(f) {
							return field
						}
					}
					panic("field not found")
				}()
				if inField.isIgnore() {
					continue
				}

				// NOTE: initialized embeded field.
				anchor := outV
				for _, index := range outField.index {
					v := anchor.Field(index)
					if v.Kind() == reflect.Ptr {
						if !v.CanSet() {
							break
						}
						if v.IsNil() {
							v.Set(reflect.New(v.Type().Elem()))
						}
						anchor = v.Elem()
					}
				}

				v := inV.FieldByIndex(inField.index)
				target := outV.FieldByIndex(outField.index)

				// NOTE: private field
				if !v.CanInterface() {
					continue
				}

				if target.Kind() != reflect.Ptr {
					target = target.Addr()
				}
				if err := c.new(v.Interface(), c.field+"."+outField.name).Convert(target.Interface()); err != nil {
					return err
				}
			}
		}
	default:
		return c.new(c.value, c.field).Convert(out)
	}
	return nil
}

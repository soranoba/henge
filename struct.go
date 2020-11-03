package henge

import (
	"errors"
	"reflect"
)

// Struct converts the input to struct type.
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
		err = ErrUnsupportedType
	}
	return &StructConverter{converter: c.converter, value: value, err: err}
}

// StructConverter is a converter that converts a struct type to another type.
type StructConverter struct {
	converter
	value interface{}
	err   error
}

// Convert converts the input to the out type and assigns it.
// If the conversion fails, the method returns an error.
func (c *StructConverter) Convert(out interface{}) error {
	var (
		err error
	)

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

	if beforeCallback, ok := outV.Addr().Interface().(BeforeCallback); ok {
		if err = beforeCallback.BeforeConvert(c.value, &c.converter); err != nil {
			goto failed
		}
	}

	switch outV.Kind() {
	case reflect.Struct:
		inV := reflect.Indirect(reflect.ValueOf(c.value))

		// NOTE: Types that are simply converted (it also copies private fields)
		if inV.Type().ConvertibleTo(outV.Type()) {
			outV.Set(inV.Convert(outV.Type()))
			break
		}

		inFields := getStructFields(inV.Type())
	Loop:
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

				v := inV.FieldByIndex(inField.index)
				// NOTE: private field
				if !v.CanInterface() {
					continue
				}
				conv := c.new(v.Interface(), c.field+"."+outField.name)

				// NOTE: initialized embedded field.
				anchor := outV
				for _, index := range outField.index {
					v := anchor.Field(index)
					if v.Kind() == reflect.Ptr {
						if !v.CanSet() {
							continue Loop
						}
						if conv.isNil {
							// NOTE: set nil.
							v.Set(reflect.New(v.Type()).Elem())
							continue Loop
						}
						if v.IsNil() {
							v.Set(reflect.New(v.Type().Elem()))
						}
						anchor = v.Elem()
					} else {
						anchor = v
					}
				}

				target := outV.FieldByIndex(outField.index)
				if target.Kind() != reflect.Ptr {
					target = target.Addr()
				}
				if err = conv.Convert(target.Interface()); err != nil {
					goto failed
				}
			}
		}
	default:
		err = c.new(c.value, c.field).Convert(out)
	}

	if afterCallback, ok := outV.Addr().Interface().(AfterCallback); ok {
		err = afterCallback.AfterConvert(c.value, &c.converter)
	}

failed:
	var convertError *ConvertError
	if err != nil && !errors.As(err, &convertError) {
		var srcType reflect.Type
		if reflect.ValueOf(c.value).IsValid() {
			srcType = reflect.ValueOf(c.value).Type()
		}
		err = &ConvertError{
			Field:   c.field,
			SrcType: srcType,
			DstType: outV.Type(),
			Value:   c.value,
			Err:     err,
		}
	}
	return err
}

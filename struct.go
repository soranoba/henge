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

	if err != nil {
		err = c.wrapConvertError(c.value, reflect.ValueOf((*uint64)(nil)).Type().Elem(), err)
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
	outV := reflect.ValueOf(out)
	if outV.Kind() != reflect.Ptr {
		panic("out must be ptr")
	}
	return c.convert(outV.Elem())
}

func (c *StructConverter) convert(outV reflect.Value) error {
	if c.err != nil {
		return c.wrapConvertError(c.value, outV.Type(), c.err)
	}
	if c.isNil {
		return nil
	}

	var err error
	elemOutV := toInitializedNonPtrValue(outV)

	if beforeCallback, ok := elemOutV.Addr().Interface().(BeforeCallback); ok {
		if err = beforeCallback.BeforeConvert(c.value, &c.converter); err != nil {
			goto failed
		}
	}

	switch elemOutV.Kind() {
	case reflect.Struct:
		inV := reflect.Indirect(reflect.ValueOf(c.value))

		// NOTE: Types that are simply converted (it also copies private fields)
		if inV.Type().ConvertibleTo(elemOutV.Type()) {
			elemOutV.Set(inV.Convert(elemOutV.Type()))
			break
		}

		inFields := getStructFields(inV.Type())
	Loop:
		for _, outField := range getStructFields(elemOutV.Type()) {
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
				anchor := elemOutV
				for i, index := range outField.index {
					v := anchor.Field(index)
					if v.Kind() == reflect.Ptr {
						if !v.CanSet() {
							continue Loop
						}
						if conv.isNil {
							if i == len(outField.index)-1 { // last index only.
								// NOTE: set nil.
								v.Set(reflect.New(v.Type()).Elem())
								continue Loop
							} else if v.IsNil() {
								continue Loop
							}
						}
						if v.IsNil() {
							v.Set(reflect.New(v.Type().Elem()))
						}
						anchor = v.Elem()
					} else {
						anchor = v
					}
				}

				target := elemOutV.FieldByIndex(outField.index)
				if err = conv.convert(target); err != nil {
					goto failed
				}
			}
		}
	default:
		err = c.new(c.value, c.field).convert(outV)
	}

	if afterCallback, ok := elemOutV.Addr().Interface().(AfterCallback); ok {
		err = afterCallback.AfterConvert(c.value, &c.converter)
	}

failed:
	var convertError *ConvertError
	if err != nil && !errors.As(err, &convertError) {
		err = c.wrapConvertError(c.value, outV.Type(), err)
	}
	return err
}

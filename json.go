package henge

import "reflect"

type (
	// JSONValueConverter is a converter that converts a JSON value type to another type.
	JSONValueConverter struct {
		Converter
	}
	// JSONArrayConverter is a converter that converts a JSON array type to another type.
	JSONArrayConverter struct {
		*SliceConverter
	}
	// JSONObjectConverter is a converter that converts a JSON object type to another type.
	JSONObjectConverter struct {
		*baseConverter
		value map[string]interface{}
		err   error
	}
)

// --------------------------------------------------------------------- //
// ValueConverter
// --------------------------------------------------------------------- //

// JSONValue converts the input to JSON value (boolean or string or numeric or array or map)
func (c *ValueConverter) JSONValue() *JSONValueConverter {
	if c.isNil {
		return &JSONValueConverter{Converter: c}
	}

	switch reflect.Indirect(c.reflectValue).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &JSONValueConverter{Converter: c.Int()}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &JSONValueConverter{Converter: c.Uint()}
	case reflect.Float32, reflect.Float64:
		return &JSONValueConverter{Converter: c.Float()}
	case reflect.Bool:
		return &JSONValueConverter{Converter: c.Bool()}
	case reflect.Array, reflect.Slice:
		return &JSONValueConverter{Converter: c.JSONArray()}
	case reflect.Map, reflect.Struct:
		return &JSONValueConverter{Converter: c.JSONObject()}
	default:
		return &JSONValueConverter{Converter: c.String()}
	}
}

func (c *ValueConverter) JSONArray() *JSONArrayConverter {
	newConv := c.new(c.value, c.field)
	newConv.opts.sliceOpts.valueConversionFunc = func(converter *ValueConverter) Converter {
		return converter.JSONValue()
	}
	return &JSONArrayConverter{SliceConverter: newConv.Slice()}
}

func (c *ValueConverter) JSONObject() *JSONObjectConverter {
	if c.isNil {
		return &JSONObjectConverter{baseConverter: c.baseConverter, value: nil, err: nil}
	}

	var out map[string]interface{}
	newConv := c.new(c.value, c.field)
	newConv.opts.mapOpts.keyType = reflect.TypeOf((*string)(nil)).Elem()
	newConv.opts.mapOpts.keyConversionFunc = func(converter *ValueConverter) Converter {
		return converter.String()
	}
	newConv.opts.mapOpts.valueConversionFunc = func(converter *ValueConverter) Converter {
		return converter.JSONValue()
	}
	err := newConv.Convert(&out)
	return &JSONObjectConverter{baseConverter: newConv.baseConverter, value: out, err: err}
}

// --------------------------------------------------------------------- //
// JSONValueConverter
// --------------------------------------------------------------------- //

func (c *JSONValueConverter) Result() (interface{}, error) {
	return c.Interface(), c.Error()
}

func (c *JSONValueConverter) Value() interface{} {
	return c.Interface()
}

// --------------------------------------------------------------------- //
// JSONObjectConverter
// --------------------------------------------------------------------- //

func (c *JSONObjectConverter) Result() (map[string]interface{}, error) {
	return c.value, c.err
}

func (c *JSONObjectConverter) Value() map[string]interface{} {
	return c.value
}

func (c *JSONObjectConverter) Interface() interface{} {
	return c.value
}

func (c *JSONObjectConverter) Error() error {
	return c.err
}

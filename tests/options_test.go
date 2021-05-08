package tests

import (
	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithSliceValueConverter_conversionFailed(t *testing.T) {
	in := []string{"1.5", "2", "2.5", "a"}

	assert.Error(t, henge.New(in, henge.WithSliceValueConverter(func(converter *henge.ValueConverter) henge.Converter {
		return converter.Float().Int()
	})).Slice().Error())
}

func TestWithMapKeyConverter_conversionFailed(t *testing.T) {
	in := map[interface{}]interface{}{
		"1.0": map[float64]interface{}{1.5: "a"},
		"b": map[uint64]interface{}{2: "b"},
	}

	assert.Error(t, henge.New(in, henge.WithMapKeyConverter(func(converter *henge.ValueConverter) henge.Converter {
			return converter.Float().Int()
	})).Map().Error())
}

func TestWithValueKeyConverter_conversionFailed(t *testing.T) {
	in := map[interface{}]interface{}{
		"a": map[interface{}]interface{}{"a.1": "a", "a.2": "b"},
		"b": map[interface{}]interface{}{"b.1": 2.5, "b.2": 2},
	}
	assert.Error(t, henge.New(in, henge.WithMapValueConverter(func(converter *henge.ValueConverter) henge.Converter {
		return converter.Float().Int()
	})).Map().Error())

	in = map[interface{}]interface{}{
		"b": map[interface{}]interface{}{"b.1": 2.5, "b.2": 2},
		"c": struct { X string }{X: "a"},
	}
	assert.Error(t, henge.New(in, henge.WithMapValueConverter(func(converter *henge.ValueConverter) henge.Converter {
		return converter.Float().Int()
	})).Map().Error())
}

package tests

import (
	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithMapKeyConverter_conversionFailed(t *testing.T) {
	in := map[interface{}]interface{}{
		"1.0": map[float64]interface{}{1.5: "a"},
		"b": map[uint64]interface{}{2: "b"},
	}

	assert.Error(t, henge.New(in, henge.WithMapKeyConverter(func(keyConverter *henge.ValueConverter) henge.Converter {
			return keyConverter.Float().Int()
	})).Map().Error())
}

func TestWithValueKeyConverter_conversionFailed(t *testing.T) {
	in := map[interface{}]interface{}{
		"a": map[interface{}]interface{}{"a.1": "a", "a.2": "b"},
		"b": map[interface{}]interface{}{"b.1": 2.5, "b.2": 2},
	}
	assert.Error(t, henge.New(in, henge.WithMapValueConverter(func(key interface{}, valueConverter *henge.ValueConverter) henge.Converter {
		return valueConverter.Float().Int()
	})).Map().Error())

	in = map[interface{}]interface{}{
		"b": map[interface{}]interface{}{"b.1": 2.5, "b.2": 2},
		"c": struct { X string }{X: "a"},
	}
	assert.Error(t, henge.New(in, henge.WithMapValueConverter(func(key interface{}, valueConverter *henge.ValueConverter) henge.Converter {
		return valueConverter.Float().Int()
	})).Map().Error())
}

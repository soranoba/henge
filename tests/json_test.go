package tests

import (
	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONValueConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil).JSONValue()
}

func TestJSONArrayConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil).JSONArray()
}

func TestJSONObjectConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil).JSONObject()
}

func TestJSONValueConverter_nil(t *testing.T) {
	val, err := henge.New(nil).JSONValue().Result()
	assert.Nil(t, val)
	assert.NoError(t, err)
}

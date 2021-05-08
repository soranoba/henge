package tests

import (
	"github.com/soranoba/henge"
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

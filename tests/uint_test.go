package tests

import (
	"testing"

	"github.com/soranoba/henge/v2"
	"github.com/stretchr/testify/assert"
)

func TestUnsignedIntegerConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil).Uint()
}

func TestUnsignedIntegerConverter_Convert(t *testing.T) {
	var s string
	if assert.NoError(t, henge.New(24.0).Uint().Convert(&s)) {
		assert.Equal(t, "24", s)
	}

	var i int
	if assert.NoError(t, henge.New(24.0).Uint().Convert(&i)) {
		assert.Equal(t, 24, i)
	}
}

func TestUnsignedIntegerConverter_Ptr(t *testing.T) {
	ptr, err := henge.New(struct{}{}).Uint().Ptr().Result()
	assert.Nil(t, ptr)
	assert.EqualError(t, err, "Failed to convert from struct {} to uint64: fields=, value=struct {}{}, error=unsupported type")

	ptr, err = henge.New(124).Uint().Ptr().Result()
	if assert.NotNil(t, ptr) {
		assert.Equal(t, uint64(124), *ptr)
	}
	assert.NoError(t, err)

	// NOTE: nil treats as a zero value, but Ptr keeps nil
	ptr, err = henge.New((*int)(nil)).Uint().Ptr().Result()
	assert.Nil(t, ptr)
	assert.NoError(t, err)

	ptr, err = henge.New((*struct{})(nil)).Uint().Ptr().Result()
	assert.Nil(t, ptr)
	assert.EqualError(t, err, "Failed to convert from *struct {} to uint64: fields=, value=(*struct {})(nil), error=unsupported type")
}

func TestUnsignedIntegerPtrConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil).UintPtr()
}

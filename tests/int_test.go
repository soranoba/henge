package tests

import (
	"testing"

	"github.com/soranoba/henge/v2"
	"github.com/stretchr/testify/assert"
)

func TestIntegerConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil).Int()
}

func TestIntegerConverter_Ptr(t *testing.T) {
	ptr, err := henge.New(struct{}{}).Int().Ptr().Result()
	assert.Nil(t, ptr)
	assert.EqualError(t, err, "Failed to convert from struct {} to int64: fields=, value=struct {}{}, error=unsupported type")

	ptr, err = henge.New("24").Int().Ptr().Result()
	if assert.NotNil(t, ptr) {
		assert.Equal(t, int64(24), *ptr)
	}
	assert.NoError(t, err)

	// NOTE: nil treats as a zero value, but Ptr keeps nil
	ptr, err = henge.New((*uint)(nil)).Int().Ptr().Result()
	assert.Nil(t, ptr)
	assert.NoError(t, err)

	ptr, err = henge.New((*struct{})(nil)).Int().Ptr().Result()
	assert.Nil(t, ptr)
	assert.EqualError(t, err, "Failed to convert from *struct {} to int64: fields=, value=(*struct {})(nil), error=unsupported type")
}

func TestIntegerPtrConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil).IntPtr()
}

package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestFloatConverter(t *testing.T) {
	var _ henge.Converter = henge.New(nil).Float()
}

func TestFloatConverter_Ptr(t *testing.T) {
	ptr, err := henge.New(struct{}{}).Float().Ptr().Result()
	assert.Nil(t, ptr)
	assert.EqualError(t, err, "Failed to convert from struct {} to float64: fields=, value=struct {}{}, error=unsupported type")

	ptr, err = henge.New("24.5").Float().Ptr().Result()
	if assert.NotNil(t, ptr) {
		assert.Equal(t, float64(24.5), *ptr)
	}
	assert.NoError(t, err)

	// NOTE: nil treats as a zero value, but Ptr keeps nil
	ptr, err = henge.New((*int)(nil)).Float().Ptr().Result()
	assert.Nil(t, ptr)
	assert.NoError(t, err)

	ptr, err = henge.New((*struct{})(nil)).Float().Ptr().Result()
	assert.Nil(t, ptr)
	assert.EqualError(t, err, "Failed to convert from *struct {} to float64: fields=, value=(*struct {})(nil), error=unsupported type")
}

package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

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
	assert.Error(t, err)

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
	assert.Error(t, err)
}

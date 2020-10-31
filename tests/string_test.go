package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestStringConverter_Convert(t *testing.T) {
	var s string
	if assert.NoError(t, henge.New(24.0).String().Convert(&s)) {
		assert.Equal(t, "24", s)
	}

	var i int
	if assert.NoError(t, henge.New(24.0).String().Convert(&i)) {
		assert.Equal(t, 24, i)
	}
}

func TestStringConverter_Ptr(t *testing.T) {
	ptr, err := henge.New(struct{}{}).String().Ptr().Result()
	assert.Nil(t, ptr)
	assert.Error(t, err)

	ptr, err = henge.New(1).String().Ptr().Result()
	if assert.NotNil(t, ptr) {
		assert.Equal(t, "1", *ptr)
	}
	assert.NoError(t, err)

	// NOTE: nil treats as a zero value, but Ptr keeps nil
	ptr, err = henge.New((*int)(nil)).String().Ptr().Result()
	assert.Nil(t, ptr)
	assert.NoError(t, err)

	ptr, err = henge.New((*struct{})(nil)).String().Ptr().Result()
	assert.Nil(t, ptr)
	assert.Error(t, err)
}

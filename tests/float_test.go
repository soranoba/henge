package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestFloatConverter_Ptr(t *testing.T) {
	ptr, err := henge.New(struct{}{}).Float().Ptr().Result()
	assert.Nil(t, ptr)
	assert.Error(t, err)

	ptr, err = henge.New("24.5").Float().Ptr().Result()
	if assert.NotNil(t, ptr) {
		assert.Equal(t, float64(24.5), *ptr)
	}
	assert.NoError(t, err)
}

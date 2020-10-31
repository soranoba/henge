package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestIntegerConverterPtr(t *testing.T) {
	ptr, err := henge.New(struct{}{}).Int().Ptr().Result()
	assert.Nil(t, ptr)
	assert.Error(t, err)

	ptr, err = henge.New("24").Int().Ptr().Result()
	if assert.NotNil(t, ptr) {
		assert.Equal(t, int64(24), *ptr)
	}
	assert.NoError(t, err)
}

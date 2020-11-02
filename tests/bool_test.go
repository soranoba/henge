package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestBoolConverter_Ptr(t *testing.T) {
	ptr, err := henge.New(nil).Bool().Ptr().Result()
	assert.Nil(t, ptr)
	assert.Error(t, err)

	ptr, err = henge.New("aaaa").Bool().Ptr().Result()
	if assert.NotNil(t, ptr) {
		assert.Equal(t, true, *ptr)
	}
	assert.NoError(t, err)

	// NOTE: nil treats as a zero value, but Ptr keeps nil
	ptr, err = henge.New((*int)(nil)).Bool().Ptr().Result()
	assert.Nil(t, ptr)
	assert.NoError(t, err)

	ptr, err = henge.New((*struct{})(nil)).Bool().Ptr().Result()
	assert.Nil(t, ptr)
	assert.NoError(t, err)
}

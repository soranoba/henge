package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestSliceConverter_Nil(t *testing.T) {
	s, err := henge.New(([]int)(nil)).Slice().Result()
	assert.NoError(t, err)
	assert.Nil(t, s)

	s, err = henge.New((*int)(nil)).Slice().Result()
	assert.Error(t, err)
	assert.Nil(t, s)
}

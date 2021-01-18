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
	assert.EqualError(t, err, "Failed to convert from *int to []interface {}: fields=, value=(*int)(nil), error=unsupported type")
	assert.Nil(t, s)
}

func TestSliceConverter_PtrSlice(t *testing.T) {
	s := make([]*int, 0)
	assert.NoError(t, henge.New([]int{1, 2, 3}).Slice().Convert(&s))
	if assert.Equal(t, 3, len(s)) {
		assert.Equal(t, 1, *s[0])
		assert.Equal(t, 2, *s[1])
		assert.Equal(t, 3, *s[2])
	}

	a := [2]*int{}
	assert.NoError(t, henge.New([]int{1, 2, 3}).Slice().Convert(&a))
	assert.Equal(t, 1, *a[0])
	assert.Equal(t, 2, *a[1])
}

func TestSliceConverter_NilValue(t *testing.T) {
	s := make([]*int, 0)
	assert.NoError(t, henge.New(make([]*uint, 3)).Slice().Convert(&s))
	if assert.Equal(t, 3, len(s)) {
		assert.Nil(t, s[0])
		assert.Nil(t, s[1])
		assert.Nil(t, s[2])
	}

	a := [2]*int{}
	assert.NoError(t, henge.New(make([]*uint, 3)).Slice().Convert(&a))
	if assert.Equal(t, 2, len(a)) {
		assert.Nil(t, a[0])
		assert.Nil(t, a[1])
	}
}

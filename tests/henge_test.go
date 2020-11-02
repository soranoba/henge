package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestConverter_InstanceGet(t *testing.T) {
	c := henge.New("some value")

	c.InstanceSet("key", 1)
	if v, ok := c.InstanceGet("key").(int); ok {
		assert.Equal(t, 1, v)
	}

	_, ok := c.InstanceGet("no").(bool)
	assert.False(t, ok)
}

func TestConverter_Interface(t *testing.T) {
	var i interface{}
	assert.NoError(t, henge.New("a").Convert(&i))
	assert.Equal(t, "a", i)
}

func TestConverter_MapToStruct(t *testing.T) {
	type Out struct {
		A string
		B string
	}
	var out Out
	assert.NoError(t, henge.New(map[string]string{"A": "a", "B": "b"}).Convert(&out))
	assert.Equal(t, "a", out.A)
	assert.Equal(t, "b", out.B)
}

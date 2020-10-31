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

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

func TestValueConverter_Convert_Interface(t *testing.T) {
	var i interface{}
	assert.NoError(t, henge.New("a").Convert(&i))
	assert.Equal(t, "a", i)
}

func TestValueConverter_Convert_MapToStruct(t *testing.T) {
	type Out struct {
		A string
		B string
	}
	var out Out
	assert.NoError(t, henge.New(map[string]string{"A": "a", "B": "b"}).Convert(&out))
	assert.Equal(t, "a", out.A)
	assert.Equal(t, "b", out.B)
}
func TestValueConverter_Model(t *testing.T) {
	type In struct {
		X string
	}
	type Out struct {
		X int
	}

	v, err := henge.New(In{X: "125"}).Model(Out{}).Result()
	assert.NoError(t, err)
	assert.Equal(t, Out{X: 125}, v)

	out := Out{}
	v, err = henge.New(In{X: "125"}).Model(&out).Result()
	assert.NoError(t, err)
	assert.Equal(t, &out, v)
	assert.Equal(t, Out{X: 125}, out)

	// Case. Non assignable
	type T struct {
		out *Out
	}
	v, err = henge.New(In{X: "125"}).Model(*(T{out: &Out{}}.out)).Result()
	assert.NoError(t, err)
	assert.Equal(t, Out{X: 125}, v)

	v, err = henge.New(In{X: "125"}).Model(T{}.out).Result()
	assert.NoError(t, err)
	if assert.NotNil(t, v) {
		assert.Equal(t, Out{X: 125}, *(v.(*Out)))
	}

	// Case. Conversion fails
	v, err = henge.New(In{X: "125"}).Model("").Result()
	assert.Error(t, err)
	assert.Equal(t, "", v)

	v, err = henge.New(In{X: "125"}).Model((*string)(nil)).Result()
	assert.Error(t, err)
	assert.Nil(t, v)
}

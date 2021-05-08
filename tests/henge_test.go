package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestConverter_Get(t *testing.T) {
	c := henge.New("some value")

	c.Set("key", 1)
	if v, ok := c.Get("key").(int); ok {
		assert.Equal(t, 1, v)
	}

	assert.Nil(t, c.Get("no"))
}

func TestConverter_SetValues(t *testing.T) {
	c := henge.New("some value")

	c.SetValues(map[string]interface{}{"a": 1, "b": 2})
	if v, ok := c.Get("a").(int); ok {
		assert.Equal(t, 1, v)
	}
	if v, ok := c.Get("b").(int); ok {
		assert.Equal(t, 2, v)
	}

	assert.Nil(t, c.Get("c"))
}

func TestValueConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil)
}

func TestValueConverter_Value(t *testing.T) {
	assert.Equal(t, nil, henge.New(nil).Value())
	assert.Equal(t, (*string)(nil), henge.New((*string)(nil)).Value())
	assert.Equal(t, 1, henge.New(1).Value())

	var i int
	assert.Equal(t, &i, henge.New(&i).Value())
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
	assert.EqualError(t, err, "Failed to convert from tests.In to string: fields=, value=tests.In{X:\"125\"}, error=unsupported type")
	assert.Equal(t, "", v)

	v, err = henge.New(In{X: "125"}).Model((*string)(nil)).Result()
	assert.EqualError(t, err, "Failed to convert from tests.In to *string: fields=, value=tests.In{X:\"125\"}, error=unsupported type")
	assert.Nil(t, v)
}

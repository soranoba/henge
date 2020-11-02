package tests

import (
	"testing"
	"time"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestMapConverter_PrivateField(t *testing.T) {
	// NOTE: private fields cannot be copied
	m, err := henge.New(time.Now()).Map().Result()
	assert.NoError(t, err)
	assert.Equal(t, map[interface{}]interface{}{}, m)
}

func TestMapConverter_Nil(t *testing.T) {
	m, err := henge.New((*struct{})(nil)).Map().Result()
	assert.NoError(t, err)
	assert.Nil(t, m)

	m, err = henge.New((map[string]string)(nil)).Map().Result()
	assert.NoError(t, err)
	assert.Nil(t, m)

	m, err = henge.New((*int)(nil)).Map().Result()
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestMapConverter_ConvertPtr(t *testing.T) {
	src := map[string]int{
		"a": 1,
		"b": 2,
	}
	var dst1 map[string]*int
	assert.NoError(t, henge.New(src).Map().Convert(&dst1))
	assert.Equal(t, "1", henge.New(dst1["a"]).String().Value())
	assert.Equal(t, "2", henge.New(dst1["b"]).String().Value())

	var dst2 map[string]*uint
	assert.NoError(t, henge.New(src).Map().Convert(&dst2))
	assert.Equal(t, "1", henge.New(dst2["a"]).String().Value())
	assert.Equal(t, "2", henge.New(dst2["b"]).String().Value())

	var dst3 map[string]*float64
	assert.NoError(t, henge.New(src).Map().Convert(&dst3))
	assert.Equal(t, "1", henge.New(dst3["a"]).String().Value())
	assert.Equal(t, "2", henge.New(dst3["b"]).String().Value())

	var dst4 map[string]*bool
	assert.NoError(t, henge.New(src).Map().Convert(&dst4))
	assert.Equal(t, true, henge.New(dst4["a"]).Bool().Value())
	assert.Equal(t, true, henge.New(dst4["b"]).Bool().Value())

	var dst5 map[string]*string
	assert.NoError(t, henge.New(src).Map().Convert(&dst5))
	assert.Equal(t, "1", henge.New(dst5["a"]).String().Value())
	assert.Equal(t, "2", henge.New(dst5["b"]).String().Value())
}

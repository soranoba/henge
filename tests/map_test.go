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

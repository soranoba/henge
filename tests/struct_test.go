package tests

import (
	"testing"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

func TestStructConverter_EmbededField(t *testing.T) {
	type In struct {
		A string
		B string
	}
	type Embeded2 struct {
		A string
	}
	type Embeded1 struct {
		*Embeded2
		B string
	}
	type Out struct {
		*Embeded1
		A string
	}

	in := In{A: "a", B: "b"}
	out := Out{}
	if err := henge.New(in).Struct().Convert(&out); err != nil {
		assert.NoError(t, err)
	}
	if assert.NotNil(t, out.Embeded1) && assert.NotNil(t, out.Embeded2) {
		// NOTE: If the names conflict, it will assign to everything possible.
		assert.Equal(t, "a", out.A)
		assert.Equal(t, "a", out.Embeded1.Embeded2.A)
		assert.Equal(t, "b", out.B)
	}

	out = Out{A: "a", Embeded1: &Embeded1{Embeded2: &Embeded2{A: "Embeded2.a"}, B: "b"}}
	in = In{}
	if err := henge.New(out).Struct().Convert(&in); err != nil {
		assert.NoError(t, err)
	}
	// NOTE: If the input has the same name properties, the higher-level property takes precedence.
	assert.Equal(t, "a", in.A)
	assert.Equal(t, "b", in.B)
}

func TestStructConverter_IgnoreField(t *testing.T) {
	type Embeded1 struct {
		X string `henge:"-"`
		Y string
	}
	type Embeded2 struct {
		Z string
	}

	type In struct {
		Embeded1 `henge:"-"`
		Embeded2
		A string `henge:"-"`
		B string
	}
	type Out struct {
		*Embeded1
		*Embeded2 `henge:"-"`
		A         string
		B         string `henge:"-"`
		X         string
		Y         string
		Z         string
	}

	in := In{A: "a", B: "b", Embeded1: Embeded1{X: "x", Y: "y"}, Embeded2: Embeded2{Z: "z"}}
	out := Out{}
	if err := henge.New(in).Struct().Convert(&out); err != nil {
		assert.NoError(t, err)
	}
	// NOTE: If ignore is specified somewhere in the path, it will not be copied.
	assert.Equal(t, "", out.A)
	assert.Equal(t, "", out.B)
	assert.Equal(t, "", out.X)
	assert.Equal(t, "", out.Y)
	assert.Equal(t, "z", out.Z)
	assert.Nil(t, out.Embeded1)
	assert.Nil(t, out.Embeded2)
}

package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/soranoba/henge"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
	Age  int
}

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

func TestStructConverter_InternalField(t *testing.T) {
	var out struct {
		wall uint64
	}

	// NOTE: private fields cannot be copied
	assert.NoError(t, henge.New(time.Now()).Struct().Convert(&out))
	assert.Equal(t, uint64(0), out.wall)
}

func TestStructConverter_SameStruct(t *testing.T) {
	now := time.Now()

	// NOTE: For the same struct, the private fields will also be copied
	var time time.Time
	assert.NoError(t, henge.New(now).Struct().Convert(&time))
	assert.Equal(t, now, time)
}

func TestStructConverter_SameStructPtr(t *testing.T) {
	now := time.Now()

	var time time.Time
	assert.NoError(t, henge.New(&now).Struct().Convert(&time))
	assert.Equal(t, now, time)
}

func TestStructConverter_Nil(t *testing.T) {
	var out *struct{}
	assert.NoError(t, henge.New((*struct{})(nil)).Struct().Convert(&out))
	assert.Nil(t, out)

	assert.Error(t, henge.New((*int)(nil)).Struct().Convert(&out))
	assert.Nil(t, out)
}

type BeforeCallbackT struct {
	Name string
	Age  int
}

func (t *BeforeCallbackT) BeforeConvert(src interface{}, converter henge.Converter) error {
	if _, ok := src.(User); ok {
		return nil
	}
	return errors.New("failed")
}

type AfterCallbackT struct {
	Name string
	Age  int
}

func (t *AfterCallbackT) AfterConvert(src interface{}, converter henge.Converter) error {
	if u, ok := src.(User); ok {
		diff, _ := converter.InstanceGet("diff").(int)
		t.Age = u.Age + diff
		return nil
	}
	return errors.New("failed")
}

func TestStructConverter_Callbacks(t *testing.T) {
	user := User{
		Name: "Alice",
		Age:  25,
	}

	out1 := BeforeCallbackT{}
	assert.NoError(t, henge.New(user).Struct().Convert(&out1))
	assert.Equal(t, user.Name, out1.Name)
	assert.Equal(t, user.Age, out1.Age)

	out1 = BeforeCallbackT{}
	assert.Error(t, henge.New(struct{ Name string }{"Bob"}).Convert(&out1))

	out2 := AfterCallbackT{}
	conv := henge.New(user)
	conv.InstanceSet("diff", 23)
	assert.NoError(t, conv.Struct().Convert(&out2))
	assert.Equal(t, user.Name, out2.Name)
	assert.Equal(t, 48, out2.Age)

	out2 = AfterCallbackT{}
	assert.Error(t, henge.New(struct{ Name string }{"Carol"}).Convert(&out2))
}

func TestStructConverter_NilField(t *testing.T) {
	type In struct {
		A *string
		B *uint
		C *int
		D *bool
		E *float64
	}
	type Out struct {
		A *string
		B *uint
		C *int
		D *bool
		E *float64
	}
	var out Out
	assert.NoError(t, henge.New(&In{}).Convert(&out))
	assert.Nil(t, out.A)
	assert.Nil(t, out.B)
	assert.Nil(t, out.C)
	assert.Nil(t, out.D)
	assert.Nil(t, out.E)
}

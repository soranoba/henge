package tests

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/soranoba/henge/v2"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
	Age  int
}

func TestStructConverter_interface(t *testing.T) {
	var _ henge.Converter = henge.New(nil).Struct()
}

func TestStructConverter_EmbeddedField(t *testing.T) {
	type In struct {
		A string
		B string
	}
	type Embedded2 struct {
		A string
	}
	type Embedded1 struct {
		*Embedded2
		B string
	}
	type Out struct {
		*Embedded1
		A string
	}

	in := In{A: "a", B: "b"}
	out := Out{}
	if err := henge.New(in).Struct().Convert(&out); err != nil {
		assert.NoError(t, err)
	}
	if assert.NotNil(t, out.Embedded1) && assert.NotNil(t, out.Embedded2) {
		// NOTE: If the names conflict, it will assign to everything possible.
		assert.Equal(t, "a", out.A)
		assert.Equal(t, "a", out.Embedded1.Embedded2.A)
		assert.Equal(t, "b", out.B)
	}

	out = Out{A: "a", Embedded1: &Embedded1{Embedded2: &Embedded2{A: "Embedded2.a"}, B: "b"}}
	in = In{}
	if err := henge.New(out).Struct().Convert(&in); err != nil {
		assert.NoError(t, err)
	}
	// NOTE: If the input has the same name properties, the higher-level property takes precedence.
	assert.Equal(t, "a", in.A)
	assert.Equal(t, "b", in.B)
}

func TestStructConverter_EmbeddedField2(t *testing.T) {
	type In struct {
		A *string
		B *string
		C *string
	}
	type Embedded2 struct {
		A string
		C string
	}
	type Embedded1 struct {
		*Embedded2
		B string
	}
	type Out struct {
		*Embedded1
		A *string
	}

	in := In{A: henge.ToStringPtr("a"), B: nil, C: henge.ToStringPtr("c")}
	out := Out{}
	if err := henge.New(in).Struct().Convert(&out); err != nil {
		assert.NoError(t, err)
	}
	if assert.NotNil(t, out.Embedded1) && assert.NotNil(t, out.Embedded2) {
		// NOTE: If the names conflict, it will assign to everything possible.
		assert.Equal(t, "a", henge.ToString(out.A))
		assert.Equal(t, "a", out.Embedded1.Embedded2.A)
		assert.Equal(t, "", out.B)
		assert.Equal(t, "c", out.C)
	}
}

func TestStructConverter_EmbeddedPtrField(t *testing.T) {
	type In struct {
		A string
	}
	type Embedded3 struct {
		A string
	}
	type Embedded2 struct {
		*Embedded3
	}
	type Embedded1 struct {
		Embedded2
	}
	type Out struct {
		*Embedded1
	}
	in := In{A: "a"}
	out := Out{}
	assert.NoError(t, henge.New(in).Struct().Convert(&out))
	assert.Equal(t, "a", out.A)
}

func TestStructConverter_IgnoreField(t *testing.T) {
	type Embedded1 struct {
		X string `henge:"-"`
		Y string
	}
	type Embedded2 struct {
		Z string
	}

	type In struct {
		Embedded1 `henge:"-"`
		Embedded2
		A string `henge:"-"`
		B string
	}
	type Out struct {
		*Embedded1
		*Embedded2 `henge:"-"`
		A          string
		B          string `henge:"-"`
		X          string
		Y          string
		Z          string
	}

	in := In{A: "a", B: "b", Embedded1: Embedded1{X: "x", Y: "y"}, Embedded2: Embedded2{Z: "z"}}
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
	assert.Nil(t, out.Embedded1)
	assert.Nil(t, out.Embedded2)
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

	assert.EqualError(
		t,
		henge.New((*int)(nil)).Struct().Convert(&out),
		"Failed to convert from *int to *struct {}: fields=, value=(*int)(nil), error=unsupported type",
	)
	assert.Nil(t, out)
}

type BeforeCallbackT struct {
	Name string
	Age  int
}

func (t *BeforeCallbackT) BeforeConvert(src interface{}, store henge.InstanceStore) error {
	if _, ok := src.(User); ok {
		return nil
	}
	return errors.New("failed")
}

func TestBeforeCallbackT(t *testing.T) {
	var _ henge.BeforeCallback = &BeforeCallbackT{}
}

type AfterCallbackT struct {
	Name string
	Age  int
}

func (t *AfterCallbackT) AfterConvert(src interface{}, store henge.InstanceStore) error {
	if u, ok := src.(User); ok {
		diff, _ := store.InstanceGet("diff").(int)
		t.Age = u.Age + diff
		return nil
	}
	return errors.New("failed")
}

func TestAfterCallbackT(t *testing.T) {
	var _ henge.AfterCallback = &AfterCallbackT{}
}

func TestStructConverter_Callbacks(t *testing.T) {
	user := User{
		Name: "Alice",
		Age:  25,
	}

	// NOTE: BeforeCallbackT converts only from User{} and returns an error otherwise.
	out1 := BeforeCallbackT{}
	assert.NoError(t, henge.New(user).Struct().Convert(&out1))
	assert.Equal(t, user.Name, out1.Name)
	assert.Equal(t, user.Age, out1.Age)

	out1 = BeforeCallbackT{}
	assert.EqualError(
		t,
		henge.New(&user).Struct().Convert(&out1),
		"Failed to convert from *tests.User to tests.BeforeCallbackT: fields=, value=&tests.User{Name:\"Alice\", Age:25}, error=failed",
	)
	out1 = BeforeCallbackT{}
	assert.EqualError(
		t,
		henge.New(struct{ Name string }{"Bob"}).Convert(&out1),
		"Failed to convert from struct { Name string } to tests.BeforeCallbackT: fields=, value=struct { Name string }{Name:\"Bob\"}, error=failed",
	)

	// NOTE: AfterCallbackT converts only from User{} and returns an error otherwise.
	out2 := AfterCallbackT{}
	conv := henge.New(user)
	conv.InstanceSet("diff", 23)
	assert.NoError(t, conv.Struct().Convert(&out2))
	assert.Equal(t, user.Name, out2.Name)
	assert.Equal(t, 48, out2.Age)

	out2 = AfterCallbackT{}
	assert.EqualError(
		t,
		henge.New(&user).Struct().Convert(&out2),
		"Failed to convert from *tests.User to tests.AfterCallbackT: fields=, value=&tests.User{Name:\"Alice\", Age:25}, error=failed",
	)
	out2 = AfterCallbackT{}
	assert.EqualError(
		t,
		henge.New(struct{ Name string }{"Carol"}).Convert(&out2),
		"Failed to convert from struct { Name string } to tests.AfterCallbackT: fields=, value=struct { Name string }{Name:\"Carol\"}, error=failed",
	)
}

func TestStructConverter_Callbacks2(t *testing.T) {
	user := User{
		Name: "Alice",
		Age:  25,
	}

	type InP struct {
		User *User
	}
	type InV struct {
		User User
	}
	type OutP struct {
		User *BeforeCallbackT
	}
	type OutV struct {
		User BeforeCallbackT
	}

	// NOTE: In the case of User{}, no error occurs, but in the case of *User{}, an error occurs.
	outP := OutP{}
	assert.EqualError(
		t,
		henge.New(InP{User: &user}).Convert(&outP),
		"Failed to convert from *tests.User to *tests.BeforeCallbackT: fields=.User, value=&tests.User{Name:\"Alice\", Age:25}, error=failed",
	)
	assert.NoError(t, henge.New(InV{User: user}).Convert(&outP))

	outV := OutV{}
	assert.EqualError(
		t,
		henge.New(InP{User: &user}).Convert(&outV),
		"Failed to convert from *tests.User to tests.BeforeCallbackT: fields=.User, value=&tests.User{Name:\"Alice\", Age:25}, error=failed",
	)
	assert.NoError(t, henge.New(InV{User: user}).Convert(&outV))
}

func TestStructConverter_NilField(t *testing.T) {
	type Embedded struct {
		S *string
	}
	type Embeded2 struct {
		S *string
	}
	type In struct {
		A *string
		B *uint
		C *int
		D *bool
		E *float64
		S *string
		x uint
	}
	type Out struct {
		*Embedded
		*Embeded2
		A *string
		B *uint
		C *int
		D *bool
		E *float64
	}
	out := Out{A: henge.ToStringPtr("a"), Embedded: &Embedded{S: henge.ToStringPtr("s")}}
	assert.NoError(t, henge.New(&In{}).Convert(&out))
	assert.Nil(t, out.A) // overwrite to nil
	assert.Nil(t, out.B)
	assert.Nil(t, out.C)
	assert.Nil(t, out.D)
	assert.Nil(t, out.E)
	assert.NotNil(t, out.Embedded) // not overwrite at the middle path
	assert.Nil(t, out.Embedded.S)  // overwrite to nil
	assert.Nil(t, out.Embeded2)    // keep nil
}

func TestStructConverter_MapField(t *testing.T) {
	type In struct {
		X map[string]interface{}
		Y map[string]int
	}
	type Out struct {
		X map[int]string
		Y map[string]interface{}
	}
	var out Out
	assert.NoError(t, henge.New(In{X: map[string]interface{}{"1": "a", "2": "b"}}).Convert(&out))
	if assert.NotNil(t, out.X) {
		assert.Equal(t, "a", out.X[1])
		assert.Equal(t, "b", out.X[2])
	}
	assert.Nil(t, out.Y)
}

func TestStructConverter_Error(t *testing.T) {
	type InV struct {
		X struct {
			Y string
			Z int
		}
	}
	type OutV struct {
		X struct {
			Y int
		}
	}
	outV := OutV{}
	err := henge.New(InV{X: struct {
		Y string
		Z int
	}{Y: "aa"}}).Convert(&outV)
	var convertError *henge.ConvertError
	if assert.True(t, errors.As(err, &convertError)) {
		assert.Equal(t, ".X.Y", convertError.Field)
		assert.Equal(t, "aa", convertError.Value)
		assert.Equal(t, reflect.TypeOf(string("")), convertError.SrcType)
		assert.Equal(t, reflect.TypeOf(int(1)), convertError.DstType)
	}

	type InP struct {
		X struct {
			Y *string
			Z *int
		}
	}
	type OutP struct {
		X struct {
			Y int
		}
	}
	outP := OutP{}
	err = henge.New(InP{X: struct {
		Y *string
		Z *int
	}{Y: henge.New("aa").StringPtr().Value()}}).Convert(&outP)
	if assert.True(t, errors.As(err, &convertError)) {
		assert.Equal(t, ".X.Y", convertError.Field)
		assert.Equal(t, "aa", *(convertError.Value.(*string)))
		assert.Equal(t, reflect.TypeOf((*string)(nil)), convertError.SrcType)
		assert.Equal(t, reflect.TypeOf((int)(1)), convertError.DstType)
	}
}

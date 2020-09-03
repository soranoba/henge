package henge

import (
	"reflect"
	"runtime"
	"testing"
)

type TestStruct1 struct {
	A string
	B int
	C *string
	D struct {
		D2 string
	}
	E *struct {
		E2 *string
	}
}

func TestInt(t *testing.T) {
	// int64
	i := int64(1)
	assertEqual(t, Int(i), int64(1))
	assertEqual(t, Int(&i), int64(1))

	// uint64
	ui := uint64(1)
	assertEqual(t, Int(ui), int64(1))
	assertEqual(t, Int(&ui), int64(1))

	// float64
	f := 1.2
	assertEqual(t, Int(f), int64(1))
	assertEqual(t, Int(&f), int64(1))

	// bool
	b := true
	assertEqual(t, Int(b), int64(1))
	assertEqual(t, Int(&b), int64(1))

	// string
	s := "1"
	assertEqual(t, Int(s), int64(1))
	assertEqual(t, Int(&s), int64(1))
	assertEqual(t, Int("abcd"), int64(0))

	// nil
	var null interface{}
	var ptr *int64
	assertEqual(t, Int(nil), int64(0))
	assertEqual(t, Int(null), int64(0))
	assertEqual(t, Int(ptr), int64(0))
}

func TestUint(t *testing.T) {
	// int64
	i := int64(1)
	assertEqual(t, Uint(i), uint64(1))
	assertEqual(t, Uint(&i), uint64(1))

	// uint64
	ui := uint64(1)
	assertEqual(t, Uint(ui), uint64(1))
	assertEqual(t, Uint(&ui), uint64(1))

	// float64
	f := 1.2
	assertEqual(t, Uint(f), uint64(1))
	assertEqual(t, Uint(&f), uint64(1))

	// bool
	b := true
	assertEqual(t, Uint(b), uint64(1))
	assertEqual(t, Uint(&b), uint64(1))

	// string
	s := "1"
	assertEqual(t, Uint(s), uint64(1))
	assertEqual(t, Uint(&s), uint64(1))
	assertEqual(t, Uint("abcd"), uint64(0))

	// nil
	var null interface{}
	var ptr *uint64
	assertEqual(t, Uint(nil), uint64(0))
	assertEqual(t, Uint(null), uint64(0))
	assertEqual(t, Uint(ptr), uint64(0))
}

func TestFloat(t *testing.T) {
	// int64
	i := 1
	assertEqual(t, Float(i), float64(1))
	assertEqual(t, Float(&i), float64(1))

	// uint64
	ui := uint64(1)
	assertEqual(t, Float(ui), float64(1))
	assertEqual(t, Float(&ui), float64(1))

	// float64
	f := 1.2
	assertEqual(t, Float(f), float64(1.2))
	assertEqual(t, Float(&f), float64(1.2))

	// bool
	b := true
	assertEqual(t, Float(b), float64(1))
	assertEqual(t, Float(&b), float64(1))

	// string
	s := "1.2"
	assertEqual(t, Float(s), float64(1.2))
	assertEqual(t, Float(&s), float64(1.2))
	assertEqual(t, Float("abcd"), float64(0))

	// nil
	var null interface{}
	var ptr *float64
	assertEqual(t, Float(nil), float64(0))
	assertEqual(t, Float(null), float64(0))
	assertEqual(t, Float(ptr), float64(0))
}

func TestString(t *testing.T) {
	// int64
	i := 1
	assertEqual(t, String(i), "1")
	assertEqual(t, String(&i), "1")

	// uint64
	ui := uint64(1)
	assertEqual(t, String(ui), "1")
	assertEqual(t, String(&ui), "1")

	// float64
	f := 1.2
	assertEqual(t, String(f), "1.2")
	assertEqual(t, String(&f), "1.2")
	f = 1.125
	assertEqual(t, String(f), "1.125")
	assertEqual(t, String(&f), "1.125")

	// bool
	b := true
	assertEqual(t, String(b), "true")
	assertEqual(t, String(&b), "true")
	b = false
	assertEqual(t, String(b), "false")
	assertEqual(t, String(&b), "false")

	// string
	s := "hoge"
	assertEqual(t, String(s), "hoge")
	assertEqual(t, String(&s), "hoge")

	// bytes
	by := []byte("hoge")
	assertEqual(t, String(by), "hoge")
	assertEqual(t, String(&by), "hoge")

	// nil
	var null interface{}
	var ptr *string
	assertEqual(t, String(nil), "")
	assertEqual(t, String(null), "")
	assertEqual(t, String(ptr), "")
}

func TestIntPtr(t *testing.T) {
	assertNil(t, IntPtr(nil))
	assertNil(t, IntPtr((*string)(nil)))
	assertEqual(t, *IntPtr(123), int64(123))
	assertEqual(t, *IntPtr("123"), int64(123))
}

func TestUintPtr(t *testing.T) {
	assertNil(t, UintPtr(nil))
	assertNil(t, UintPtr((*string)(nil)))
	assertEqual(t, *UintPtr(123), uint64(123))
	assertEqual(t, *UintPtr("123"), uint64(123))
}

func TestFloatPtr(t *testing.T) {
	assertNil(t, FloatPtr(nil))
	assertNil(t, FloatPtr((*string)(nil)))
	assertEqual(t, *FloatPtr(123.5), float64(123.5))
	assertEqual(t, *FloatPtr("123.5"), float64(123.5))
}

func TestStringPtr(t *testing.T) {
	assertNil(t, StringPtr(nil))
	assertNil(t, StringPtr((*int)(nil)))
	assertEqual(t, *StringPtr(123), "123")
	assertEqual(t, *StringPtr("hoge"), "hoge")
}

func TestIntSlice(t *testing.T) {
	assertNil(t, IntSlice(([]interface{})(nil)))
	assertEqual(t, IntSlice([]interface{}{"1", 2, "x"}), []int64{1, 2, 0})
	assertEqual(t, IntSlice([]string{"1", "2", "3"}), []int64{1, 2, 3})
}

func TestUintSlice(t *testing.T) {
	assertNil(t, UintSlice(([]interface{})(nil)))
	assertEqual(t, UintSlice([]interface{}{"1", 2, "x"}), []uint64{1, 2, 0})
	assertEqual(t, UintSlice([]string{"1", "2", "3"}), []uint64{1, 2, 3})
}

func TestFloatSlice(t *testing.T) {
	assertNil(t, FloatSlice(([]interface{})(nil)))
	assertEqual(t, FloatSlice([]interface{}{"1.5", 2.5, "x"}), []float64{1.5, 2.5, 0})
	assertEqual(t, FloatSlice([]string{"1.5", "2.5", "3.5"}), []float64{1.5, 2.5, 3.5})
}

func TestStringSlice(t *testing.T) {
	assertNil(t, StringSlice(([]interface{})(nil)))
	assertEqual(t, StringSlice([]interface{}{1, "2", 3.5}), []string{"1", "2", "3.5"})
	assertEqual(t, StringSlice([]int{1, 2, 3}), []string{"1", "2", "3"})
}

func TestMap(t *testing.T) {
	var s = TestStruct1{
		A: "a",
		B: 1,
		C: StringPtr("c"),
		D: struct{ D2 string }{"d2"},
		E: &struct{ E2 *string }{StringPtr("e2")},
	}

	m := MapWithDepth(s, UnlimitedDepth)
	expected := map[string]interface{}{
		"A": "a",
		"B": 1,
		"C": "c",
		"D": map[string]interface{}{
			"D2": "d2",
		},
		"E": map[string]interface{}{
			"E2": "e2",
		},
	}
	assertEqual(t, m, expected)

	// nil
	var null interface{}
	var ptr *map[string]interface{}
	assertEqual(t, Map(nil), map[string]interface{}{})
	assertEqual(t, Map(null), map[string]interface{}{})
	assertEqual(t, Map(ptr), map[string]interface{}{})
	assertEqual(t, Map(struct{ A *string }{}), map[string]interface{}{})
}

func TestMapWithDepth(t *testing.T) {
	in := struct {
		A string
		B struct {
			C string
			D struct {
				E string
			}
		}
	}{
		A: "aaa",
		B: struct {
			C string
			D struct{ E string }
		}{
			C: "ccc",
			D: struct{ E string }{
				E: "eee",
			},
		},
	}
	assertEqual(t, MapWithDepth(in, 0), map[string]interface{}{
		"A": "aaa",
		"B": struct {
			C string
			D struct{ E string }
		}{
			C: "ccc",
			D: struct{ E string }{
				E: "eee",
			},
		},
	})
	assertEqual(t, MapWithDepth(in, 1), map[string]interface{}{
		"A": "aaa",
		"B": map[string]interface{}{
			"C": "ccc",
			"D": struct{ E string }{
				E: "eee",
			},
		},
	})
	assertEqual(t, MapWithDepth(in, 2), map[string]interface{}{
		"A": "aaa",
		"B": map[string]interface{}{
			"C": "ccc",
			"D": map[string]interface{}{
				"E": "eee",
			},
		},
	})
}

func assertEqual(t *testing.T, got, expected interface{}) bool {
	if !reflect.DeepEqual(got, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Not equals:\n  file    : %s:%d\n  got     : %#v\n  expected: %#v\n", file, line, got, expected)
		return false
	}
	return true
}

func assertNotEqual(t *testing.T, got, expected interface{}) bool {
	if reflect.DeepEqual(got, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Equals:\n  file    : %s:%d\n  got     : %#v\n", file, line, got)
		return false
	}
	return true
}

func assertNil(t *testing.T, got interface{}) (ok bool) {
	if !reflect.ValueOf(got).IsNil() {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Not nil:\n  file    : %s:%d\n  got     : %#v\n", file, line, got)
		return false
	}
	return true
}

func assertNotNil(t *testing.T, got interface{}) (ok bool) {
	if reflect.ValueOf(got).IsNil() {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Nil:\n  file    : %s:%d\n  got     : %#v\n", file, line, got)
		return false
	}
	return true
}

func assertPanic(t *testing.T, fun func()) (ok bool) {
	ok = true

	defer func() {
		recover()
	}()

	fun()

	ok = false

	_, file, line, _ := runtime.Caller(1)
	t.Errorf("No panic:\n  file    : %s:%d", file, line)
	return
}

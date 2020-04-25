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
}

func TestUint(t *testing.T) {
	// int64
	i := int64(1)
	assertEqual(t, UInt(i), uint64(1))
	assertEqual(t, UInt(&i), uint64(1))

	// uint64
	ui := uint64(1)
	assertEqual(t, UInt(ui), uint64(1))
	assertEqual(t, UInt(&ui), uint64(1))

	// float64
	f := 1.2
	assertEqual(t, UInt(f), uint64(1))
	assertEqual(t, UInt(&f), uint64(1))

	// bool
	b := true
	assertEqual(t, UInt(b), uint64(1))
	assertEqual(t, UInt(&b), uint64(1))

	// string
	s := "1"
	assertEqual(t, UInt(s), uint64(1))
	assertEqual(t, UInt(&s), uint64(1))
	assertEqual(t, UInt("abcd"), uint64(0))
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
	assertEqual(t, String(f), "1.200000")
	assertEqual(t, String(&f), "1.200000")

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
}

func TestStringPtr(t *testing.T) {
	assertEqual(t, *StringPtr("hoge"), "hoge")
}

func TestMap(t *testing.T) {
	var s = TestStruct1{
		A: "a",
		B: 1,
		C: StringPtr("c"),
		D: struct{ D2 string }{"d2"},
		E: &struct{ E2 *string }{StringPtr("e2")},
	}

	m := Map(s)
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

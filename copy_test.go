package henge

import (
	"testing"
	"time"
)

func TestCopy_Primitive(t *testing.T) {
	assertPanic(t, func() {
		Copy("hoge", "fuga")
	})

	var i int64
	Copy(1, &i)
	assertEqual(t, i, int64(1))

	var f32 float32
	Copy(1, &f32)
	assertEqual(t, f32, float32(1))

	var f64 float64
	Copy(1, &f64)
	assertEqual(t, f64, float64(1))
}

func TestCopy_Nil(t *testing.T) {
	var src *string
	var dst string

	Copy(nil, &dst)
	assertEqual(t, dst, "")
	Copy(src, &dst)
	assertEqual(t, dst, "")
}

func TestCopy_Array(t *testing.T) {
	arrInt := []int{1, 2, 3}
	arrIntDst1 := make([]int, 3)
	Copy(arrInt, &arrIntDst1)
	assertEqual(t, arrIntDst1, arrInt)

	var arrIntDst2 []int
	Copy(arrInt, &arrIntDst2)
	assertEqual(t, arrIntDst2, arrInt)
}

func TestCopy_Map(t *testing.T) {
	assertPanic(t, func() {
		var m map[string]string
		Copy(m, m)
	})

	m1 := map[string]string{
		"A": "a",
		"B": "b",
	}
	var m2 map[string]string
	Copy(m1, &m2)

	assertEqual(t, m1, m2)
}

func TestCopy_Struct(t *testing.T) {
	assertPanic(t, func() {
		var s struct{}
		Copy(struct{}{}, s)
	})

	s1 := TestStruct1{
		A: "hoge",
		B: 100,
		D: struct{ D2 string }{"fuga"},
	}
	var s2 struct {
		A string
		B *int
		D struct {
			D2 *string
			D3 string
		}
		E *struct {
		}
		F *struct {
		}
	}

	Copy(s1, &s2)
	assertEqual(t, s2.A, "hoge")
	assertEqual(t, *(s2.B), 100)
	assertEqual(t, *(s2.D.D2), "fuga")
	assertEqual(t, s2.D.D3, "")
	assertEqual(t, *(s2.E), struct{}{})
	assertNil(t, s2.F)
}

func TestCopy_InternalField(t *testing.T) {
	src := time.Now()
	var dst struct {
		wall uint64
	}

	Copy(src, &dst)
	assertEqual(t, dst, struct{ wall uint64 }{})
}

func TestCopy_SameType(t *testing.T) {
	src := time.Now()
	var t1, t2 time.Time

	Copy(src, &t1)
	assertEqual(t, t1, src)

	Copy(&src, &t2)
	assertEqual(t, t2, src)
}

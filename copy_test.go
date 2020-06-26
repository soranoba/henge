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
	arrInt := [...]int{1, 2, 3}
	arrIntDst1 := make([]int, 3)
	Copy(arrInt, &arrIntDst1)
	assertEqual(t, arrIntDst1, []int{1, 2, 3})

	var arrIntDst2 []int
	Copy(arrInt, &arrIntDst2)
	assertEqual(t, arrIntDst2, []int{1, 2, 3})

	var arrIntDst3 [2]int
	Copy(arrInt, &arrIntDst3)
	assertEqual(t, arrIntDst3, [2]int{1, 2})
}

func TestCopy_Slice(t *testing.T) {
	sliceInt := []int{1, 2, 3}
	sliceIntDst1 := make([]int, 3)
	Copy(sliceInt, &sliceIntDst1)
	assertEqual(t, sliceIntDst1, []int{1, 2, 3})

	var sliceIntDst2 []int
	Copy(sliceInt, &sliceIntDst2)
	assertEqual(t, sliceIntDst2, []int{1, 2, 3})

	var sliceIntDst3 [2]int
	Copy(sliceInt, &sliceIntDst3)
	assertEqual(t, sliceIntDst3, [2]int{1, 2})
}

func TestCopy_ArrayStruct(t *testing.T) {
	type T struct {
		A string
	}

	arrSrc1 := [...]T{{A: "a"}, {A: "b"}}
	arrSrc2 := [...]*T{{A: "a"}, {A: "b"}}

	arrDst1 := make([]T, 2)
	Copy(arrSrc1, &arrDst1)
	assertEqual(t, arrDst1, []T{{A: "a"}, {A: "b"}})

	arrDst2 := make([]T, 2)
	Copy(arrSrc2, &arrDst2)
	assertEqual(t, arrDst2, []T{{A: "a"}, {A: "b"}})

	arrDst3 := make([]*T, 2)
	Copy(arrSrc1, &arrDst3)
	assertEqual(t, *arrDst3[0], T{A: "a"})
	assertEqual(t, *arrDst3[1], T{A: "b"})

	arrDst4 := make([]*T, 2)
	Copy(arrSrc2, &arrDst4)
	assertEqual(t, *arrDst4[0], T{A: "a"})
	assertEqual(t, *arrDst4[1], T{A: "b"})

	var arrDst5 []T
	Copy(arrSrc1, &arrDst5)
	assertEqual(t, arrDst5, []T{{A: "a"}, {A: "b"}})

	var arrDst6 []T
	Copy(arrSrc2, &arrDst6)
	assertEqual(t, arrDst6, []T{{A: "a"}, {A: "b"}})

	var arrDst7 []*T
	Copy(arrSrc1, &arrDst7)
	assertEqual(t, *arrDst7[0], T{A: "a"})
	assertEqual(t, *arrDst7[1], T{A: "b"})

	var arrDst8 []*T
	Copy(arrSrc2, &arrDst8)
	assertEqual(t, *arrDst8[0], T{A: "a"})
	assertEqual(t, *arrDst8[1], T{A: "b"})
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
		E: &struct{ E2 *string }{StringPtr("piyo")},
	}
	var s2 struct {
		A string
		B *int
		D *struct {
			D2 *string
			D3 string
		}
		E struct {
			E2 string
		}
		F *struct {
		}
	}

	Copy(s1, &s2)
	assertEqual(t, s2.A, "hoge")
	assertEqual(t, *(s2.B), 100)
	assertEqual(t, *(s2.D.D2), "fuga")
	assertEqual(t, s2.D.D3, "")
	assertEqual(t, s2.E.E2, "piyo")
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

func TestCopy_NilPointerInStruct(t *testing.T) {
	var t1 struct {
		A *string
		T *time.Time
	}
	var t2 struct {
		A *string
		B *int
		T *time.Time
	}

	Copy(t1, &t2)
	assertNil(t, t2.A)
	assertNil(t, t2.B)
	assertNil(t, t2.T)
}

func TestCopy_EmbededField(t *testing.T) {
	type T struct {
		A string
	}

	t0 := struct {
		A string
	}{
		A: "aaaa",
	}

	var t1 struct {
		T
	}
	var t2 struct {
		*T
	}

	Copy(t0, &t1)
	Copy(t0, &t2)
	assertEqual(t, t1.A, t0.A)
	assertEqual(t, String(t2.A), t0.A)
}

type AfterCallbackT struct {
	A string
	B string
}

func (t *AfterCallbackT) afterHenge(src interface{}) {
	t.A = "a_overwrite"
	if _, ok := src.(AfterCallbackT); ok {
		t.B = "b_overwrite"
	}
}

func TestCopy_Callbacks(t *testing.T) {
	src1 := struct {
		A string
	}{
		A: "aaaa",
	}

	var a1 AfterCallbackT
	Copy(src1, &a1)
	assertEqual(t, a1.A, "a_overwrite")
	assertEqual(t, a1.B, "")

	var a2 AfterCallbackT
	Copy(a1, &a2)
	assertEqual(t, a2.A, "a_overwrite")
	assertEqual(t, a2.B, "b_overwrite")

	var a3 AfterCallbackT
	Copy(&a1, &a3)
	assertEqual(t, a2.A, "a_overwrite")
	assertEqual(t, a2.B, "b_overwrite")
}

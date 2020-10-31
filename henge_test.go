package henge

import (
	"fmt"
	"math"
)

func ExampleNew() {
	i, err := New("-12").Int().Result()
	fmt.Printf("New(\"-12\").Int().Result() = (%#v, %#v)\n", i, err)

	err = New("abc").Int().Error()
	fmt.Printf("New(\"abc\").Error() = %#v\n", err.Error())

	ip := New("-12").Int().Ptr().Value()
	fmt.Printf("*New(\"-12\").Int().Ptr().Value() = %#v", *ip)

	// Output:
	// New("-12").Int().Result() = (-12, <nil>)
	// New("abc").Error() = "Failed to convert from string to int64: fields=, value=\"abc\", error=strconv.ParseInt: parsing \"abc\": invalid syntax"
	// *New("-12").Int().Ptr().Value() = -12
}

func ExampleValueConverter_Convert() {
	fmt.Println("int")

	var i8 int8
	fmt.Println(New(math.MaxInt16).Convert(&i8), i8)
	fmt.Println(New(math.MinInt16).Convert(&i8), i8)
	fmt.Println(New("24").Convert(&i8), i8)

	var i16 int16
	fmt.Println(New(math.MaxInt32).Convert(&i16), i16)
	fmt.Println(New(math.MinInt32).Convert(&i16), i16)
	fmt.Println(New("24").Convert(&i16), i16)

	var i32 int32
	fmt.Println(New(math.MaxInt64).Convert(&i32), i32)
	fmt.Println(New(math.MinInt64).Convert(&i32), i32)
	fmt.Println(New("24").Convert(&i32), i32)

	var i64 int64
	fmt.Println(New(math.MaxInt64).Convert(&i64), i64)
	fmt.Println(New(math.MinInt64).Convert(&i64), i64)
	fmt.Println(New("24").Convert(&i64), i64)
	fmt.Println()

	fmt.Println("uint")

	var u8 uint8
	fmt.Println(New(math.MaxUint16).Convert(&u8), u8)
	fmt.Println(New("24").Convert(&u8), u8)

	var u16 uint16
	fmt.Println(New(math.MaxUint32).Convert(&u16), u16)
	fmt.Println(New("24").Convert(&u16), u16)

	var u32 uint32
	fmt.Println(New(uint64(math.MaxUint64)).Convert(&u32), u32)
	fmt.Println(New("24").Convert(&u32), u32)

	var u64 uint64
	fmt.Println(New(uint64(math.MaxUint64)).Convert(&u64), u64)
	fmt.Println(New("24").Convert(&u64), u64)
	fmt.Println()

	fmt.Println("float")

	var f32 float32
	fmt.Println(New(math.MaxFloat64).Convert(&f32), f32)
	fmt.Println(New("24").Convert(&f32), f32)

	var f64 float64
	fmt.Println(New(math.MaxFloat64).Convert(&f64), f64)
	fmt.Println(New("24").Convert(&f64), f64)
	fmt.Println()

	fmt.Println("string")

	var s string
	fmt.Println(New(123).Convert(&s), s)
	fmt.Println()

	// Output:
	// int
	// overflows 0
	// overflows 0
	// <nil> 24
	// overflows 0
	// overflows 0
	// <nil> 24
	// overflows 0
	// overflows 0
	// <nil> 24
	// <nil> 9223372036854775807
	// <nil> -9223372036854775808
	// <nil> 24
	//
	// uint
	// overflows 0
	// <nil> 24
	// overflows 0
	// <nil> 24
	// overflows 0
	// <nil> 24
	// <nil> 18446744073709551615
	// <nil> 24
	//
	// float
	// overflows 0
	// <nil> 24
	// <nil> 1.7976931348623157e+308
	// <nil> 24
	//
	// string
	// <nil> 123
}

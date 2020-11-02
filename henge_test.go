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
	if err := New(math.MaxInt16).Convert(&i8); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i8)
	}
	if err := New(math.MinInt16).Convert(&i8); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i8)
	}
	if err := New("24").Convert(&i8); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i8)
	}

	var i16 int16
	if err := New(math.MaxInt32).Convert(&i16); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i16)
	}
	if err := New(math.MinInt32).Convert(&i16); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i16)
	}
	if err := New("24").Convert(&i16); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i16)
	}

	var i32 int32
	if err := New(math.MaxInt64).Convert(&i32); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i32)
	}
	if err := New(math.MinInt64).Convert(&i32); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i32)
	}
	if err := New("24").Convert(&i32); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i32)
	}

	var i64 int64
	if err := New(math.MaxInt64).Convert(&i64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i64)
	}
	if err := New(math.MinInt64).Convert(&i64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i64)
	}
	if err := New("24").Convert(&i64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i64)
	}

	fmt.Println()
	fmt.Println("uint")

	var u8 uint8
	if err := New(math.MaxUint16).Convert(&u8); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u8)
	}
	if err := New("24").Convert(&u8); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u8)
	}

	var u16 uint16
	if err := New(math.MaxUint32).Convert(&u16); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u16)
	}
	if err := New("24").Convert(&u16); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u16)
	}

	var u32 uint32
	if err := New(uint64(math.MaxUint64)).Convert(&u32); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u32)
	}
	if err := New("24").Convert(&u32); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u32)
	}

	var u64 uint64
	if err := New(uint64(math.MaxUint64)).Convert(&u64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u64)
	}
	if err := New("24").Convert(&u64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u64)
	}

	fmt.Println()
	fmt.Println("float")

	var f32 float32
	if err := New(math.MaxFloat64).Convert(&f32); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(f32)
	}
	if err := New("24").Convert(&f32); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(f32)
	}

	var f64 float64
	if err := New(math.MaxFloat64).Convert(&f64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(f64)
	}
	if err := New("24").Convert(&f64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(f64)
	}

	fmt.Println()
	fmt.Println("string")

	var s string
	if err := New(123).Convert(&s); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(s)
	}

	// Output:
	// int
	// Failed to convert from int64 to int8: fields=, value=32767, error=overflows
	// Failed to convert from int64 to int8: fields=, value=-32768, error=overflows
	// 24
	// Failed to convert from int64 to int16: fields=, value=2147483647, error=overflows
	// Failed to convert from int64 to int16: fields=, value=-2147483648, error=overflows
	// 24
	// Failed to convert from int64 to int32: fields=, value=9223372036854775807, error=overflows
	// Failed to convert from int64 to int32: fields=, value=-9223372036854775808, error=overflows
	// 24
	// 9223372036854775807
	// -9223372036854775808
	// 24
	//
	// uint
	// Failed to convert from uint64 to uint8: fields=, value=0xffff, error=overflows
	// 24
	// Failed to convert from uint64 to uint16: fields=, value=0xffffffff, error=overflows
	// 24
	// Failed to convert from uint64 to uint32: fields=, value=0xffffffffffffffff, error=overflows
	// 24
	// 18446744073709551615
	// 24
	//
	// float
	// Failed to convert from float64 to float32: fields=, value=1.7976931348623157e+308, error=overflows
	// 24
	// 1.7976931348623157e+308
	// 24
	//
	// string
	// 123
}

package henge

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"
)

func ExampleValueConverter_Int() {
	fmt.Println("int64 to int64")
	fmt.Printf("%v -> %v\n", math.MaxInt64, New(math.MaxInt64).Int().Value())
	fmt.Printf("%v -> %v\n", math.MinInt64, New(math.MinInt64).Int().Value())
	fmt.Println()

	fmt.Println("uint64 to int64")
	fmt.Printf("%v -> %v\n", uint64(123), New(uint64(123)).Int().Value())
	fmt.Printf("%#v\n", New(uint64(math.MaxUint64)).Int().Error().Error())
	fmt.Println()

	fmt.Println("float64 to int64")
	fmt.Printf("%v -> %v\n", 1.25, New(1.25).Int().Value())
	fmt.Printf("%v -> %v\n", -1.25, New(-1.25).Int().Value())
	fmt.Printf("%v -> %v\n", math.MaxInt32, New(float64(math.MaxInt32)).Int().Value())
	fmt.Printf("%#v\n", New(float64(math.MaxInt64)).Int().Error().Error())
	fmt.Printf("%#v\n", New(math.MaxFloat64).Int().Error().Error())
	fmt.Println()

	fmt.Println("bool to int64")
	fmt.Printf("%v -> %v\n", true, New(true).Int().Value())
	fmt.Printf("%v -> %v\n", false, New(false).Int().Value())
	fmt.Println()

	fmt.Println("string to int64")
	fmt.Printf("\"%v\" -> %v\n", math.MaxInt64, New(strconv.FormatInt(math.MaxInt64, 10)).Int().Value())
	fmt.Printf("%#v\n", New("1.5").Int().Error().Error())
	fmt.Printf("%#v\n", New("-1.5").Int().Error().Error())

	// Output:
	// int64 to int64
	// 9223372036854775807 -> 9223372036854775807
	// -9223372036854775808 -> -9223372036854775808
	//
	// uint64 to int64
	// 123 -> 123
	// "Failed to convert from uint64 to int64: fields=, error=overflows"
	//
	// float64 to int64
	// 1.25 -> 1
	// -1.25 -> -2
	// 2147483647 -> 2147483647
	// "Failed to convert from float64 to int64: fields=, error=overflows"
	// "Failed to convert from float64 to int64: fields=, error=overflows"
	//
	// bool to int64
	// true -> 1
	// false -> 0
	//
	// string to int64
	// "9223372036854775807" -> 9223372036854775807
	// "Failed to convert from string to int64: fields=, error=strconv.ParseInt: parsing \"1.5\": invalid syntax"
	// "Failed to convert from string to int64: fields=, error=strconv.ParseInt: parsing \"-1.5\": invalid syntax"
}

func TestIntegerConverterPtr(t *testing.T) {
	ptr, err := (&IntegerConverter{value: 0, err: errors.New("error")}).Ptr().Result()
	if ptr != nil || err == nil {
		t.Error(ptr)
	}
	ptr, err = (&IntegerConverter{value: 0, err: nil}).Ptr().Result()
	if ptr == nil || err != nil {
		t.Error(err)
	} else if *ptr != 0 {
		t.Error(*ptr)
	}
	ptr, err = (&IntegerConverter{value: 1, err: nil}).Ptr().Result()
	if ptr == nil || err != nil {
		t.Error(err)
	} else if *ptr != 1 {
		t.Error(*ptr)
	}
}

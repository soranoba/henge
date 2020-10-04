package henge

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"
)

func ExampleValueConverter_Uint() {
	fmt.Println("int64 to uint64")
	fmt.Printf("%v -> %v\n", math.MaxInt64, New(math.MaxInt64).Uint().Value())
	fmt.Printf("%#v\n", New(math.MinInt64).Uint().Error().Error())
	fmt.Println()

	fmt.Println("uint64 to uint64")
	fmt.Printf("%v -> %v\n", uint64(math.MaxUint64), New(uint64(math.MaxUint64)).Uint().Value())
	fmt.Println()

	fmt.Println("float64 to uint64")
	fmt.Printf("%v -> %v\n", 1.25, New(1.25).Uint().Value())
	fmt.Printf("%#v\n", New(-1.25).Uint().Error().Error())
	fmt.Printf("%v -> %v\n", math.MaxUint32, New(float64(math.MaxUint32)).Uint().Value())
	fmt.Printf("%#v\n", New(float64(math.MaxUint64)).Uint().Error().Error())
	fmt.Printf("%#v\n", New(math.MaxFloat64).Uint().Error().Error())
	fmt.Println()

	fmt.Println("bool to uint64")
	fmt.Printf("%v -> %v\n", true, New(true).Uint().Value())
	fmt.Printf("%v -> %v\n", false, New(false).Uint().Value())
	fmt.Println()

	fmt.Println("string to uint64")
	fmt.Printf("\"%v\" -> %v\n", uint64(math.MaxUint64), New(strconv.FormatUint(math.MaxUint64, 10)).Uint().Value())
	fmt.Printf("%#v\n", New("1.5").Uint().Error().Error())
	fmt.Printf("%#v\n", New("-2").Uint().Error().Error())

	// Output:
	// int64 to uint64
	// 9223372036854775807 -> 9223372036854775807
	// "Failed to convert from int to uint64: fields=, error=negative number"
	//
	// uint64 to uint64
	// 18446744073709551615 -> 18446744073709551615
	//
	// float64 to uint64
	// 1.25 -> 1
	// "Failed to convert from float64 to uint64: fields=, error=negative number"
	// 4294967295 -> 4294967295
	// "Failed to convert from float64 to uint64: fields=, error=overflows"
	// "Failed to convert from float64 to uint64: fields=, error=overflows"
	//
	// bool to uint64
	// true -> 1
	// false -> 0
	//
	// string to uint64
	// "18446744073709551615" -> 18446744073709551615
	// "Failed to convert from string to uint64: fields=, error=strconv.ParseUint: parsing \"1.5\": invalid syntax"
	// "Failed to convert from string to uint64: fields=, error=strconv.ParseUint: parsing \"-2\": invalid syntax"
}

func TestUnsignedIntegerConverterPtr(t *testing.T) {
	ptr, err := (&UnsignedIntegerConverter{value: 0, err: errors.New("error")}).Ptr().Result()
	if ptr != nil || err == nil {
		t.Error(ptr)
	}
	ptr, err = (&UnsignedIntegerConverter{value: 0, err: nil}).Ptr().Result()
	if ptr == nil || err != nil {
		t.Error(err)
	} else if *ptr != 0 {
		t.Error(*ptr)
	}
	ptr, err = (&UnsignedIntegerConverter{value: 1, err: nil}).Ptr().Result()
	if ptr == nil || err != nil {
		t.Error(err)
	} else if *ptr != 1 {
		t.Error(*ptr)
	}
}

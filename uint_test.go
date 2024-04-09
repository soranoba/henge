package henge

import (
	"fmt"
	"math"
	"strconv"
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
	// This behavior varies depending on the architecture.
	if val, err := New(float64(math.MaxUint64)).Uint().Result(); err == nil {
		if val != math.MaxUint64 {
			fmt.Printf("%v -> %v\n", uint64(math.MaxUint64), val)
		}
	}
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
	// "Failed to convert from int to uint64: fields=, value=-9223372036854775808, error=negative number"
	//
	// uint64 to uint64
	// 18446744073709551615 -> 18446744073709551615
	//
	// float64 to uint64
	// 1.25 -> 1
	// "Failed to convert from float64 to uint64: fields=, value=-1.25, error=negative number"
	// 4294967295 -> 4294967295
	// "Failed to convert from float64 to uint64: fields=, value=1.7976931348623157e+308, error=overflows"
	//
	// bool to uint64
	// true -> 1
	// false -> 0
	//
	// string to uint64
	// "18446744073709551615" -> 18446744073709551615
	// "Failed to convert from string to uint64: fields=, value=\"1.5\", error=strconv.ParseUint: parsing \"1.5\": invalid syntax"
	// "Failed to convert from string to uint64: fields=, value=\"-2\", error=strconv.ParseUint: parsing \"-2\": invalid syntax"
}

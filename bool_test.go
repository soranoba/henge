package henge

import (
	"fmt"
	"math"
)

func ExampleValueConverter_Bool() {
	fmt.Println("int64 to bool")
	fmt.Printf("%v -> %v\n", math.MinInt64, New(math.MinInt64).Bool().Value())
	fmt.Printf("%v -> %v\n", 0, New(0).Bool().Value())
	fmt.Println()

	fmt.Println("uint64 to bool")
	fmt.Printf("%v -> %v\n", uint64(math.MaxUint64), New(uint64(math.MaxUint64)).Bool().Value())
	fmt.Printf("%v -> %v\n", 0, New(0).Bool().Value())
	fmt.Println()

	fmt.Println("float64 to bool")
	fmt.Printf("%v -> %v\n", -1*math.MaxFloat64, New(-1*math.MaxFloat64).Bool().Value())
	fmt.Printf("%v -> %v\n", 0.0, New(0.0).Bool().Value())
	fmt.Println()

	fmt.Println("bool to bool")
	fmt.Printf("%v -> %v\n", true, New(true).Bool().Value())
	fmt.Printf("%v -> %v\n", false, New(false).Bool().Value())
	fmt.Println()

	fmt.Println("string to bool")
	fmt.Printf("%#v -> %v\n", "aaaa", New("aaaa").Bool().Value())
	fmt.Printf("%#v -> %v\n", "0", New("0").Bool().Value())
	fmt.Printf("%#v -> %v\n", "", New("").Bool().Value())

	// Output:
	// int64 to bool
	// -9223372036854775808 -> true
	// 0 -> false
	//
	// uint64 to bool
	// 18446744073709551615 -> true
	// 0 -> false
	//
	// float64 to bool
	// -1.7976931348623157e+308 -> true
	// 0 -> false
	//
	// bool to bool
	// true -> true
	// false -> false
	//
	// string to bool
	// "aaaa" -> true
	// "0" -> true
	// "" -> false
}

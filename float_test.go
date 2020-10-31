package henge

import (
	"fmt"
	"math"
	"strconv"
)

func ExampleValueConverter_Float() {
	fmt.Println("int64 to float64")
	fmt.Printf("%v -> %v\n", math.MaxInt64, New(math.MaxInt64).Float().Value())
	fmt.Printf("%v -> %v\n", math.MinInt64, New(math.MinInt64).Float().Value())
	fmt.Println()

	fmt.Println("uint64 to float64")
	fmt.Printf("%v -> %v\n", uint64(math.MaxUint64), New(uint64(math.MaxUint64)).Float().Value())
	fmt.Println()

	fmt.Println("float64 to float64")
	fmt.Printf("%v -> %v\n", math.MaxFloat64, New(math.MaxFloat64).Float().Value())
	fmt.Printf("%v -> %v\n", math.SmallestNonzeroFloat64, New(math.SmallestNonzeroFloat64).Float().Value())
	fmt.Printf("%v -> %v\n", -1*math.MaxFloat64, New(-1*math.MaxFloat64).Float().Value())
	fmt.Println()

	fmt.Println("bool to float64")
	fmt.Printf("%v -> %v\n", true, New(true).Float().Value())
	fmt.Printf("%v -> %v\n", false, New(false).Float().Value())
	fmt.Println()

	fmt.Println("string to float64")
	fmt.Printf("%#v -> %v\n", strconv.FormatFloat(math.MaxFloat64, 'f', 10, 64), New(strconv.FormatFloat(math.MaxFloat64, 'f', 10, 64)).Float().Value())
	fmt.Printf("%#v -> %v\n", strconv.FormatFloat(1.79769e+308, 'e', 10, 64), New(strconv.FormatFloat(1.79769e+308, 'e', 10, 64)).Float().Value())
	fmt.Printf("%#v -> %v\n", strconv.FormatFloat(math.MaxFloat64, 'e', 10, 64), New(strconv.FormatFloat(math.MaxFloat64, 'e', 10, 64)).Float().Value())
	fmt.Printf("%#v\n", New("1.1.1").Float().Error().Error())

	// Output:
	// int64 to float64
	// 9223372036854775807 -> 9.223372036854776e+18
	// -9223372036854775808 -> -9.223372036854776e+18
	//
	// uint64 to float64
	// 18446744073709551615 -> 1.8446744073709552e+19
	//
	// float64 to float64
	// 1.7976931348623157e+308 -> 1.7976931348623157e+308
	// 5e-324 -> 5e-324
	// -1.7976931348623157e+308 -> -1.7976931348623157e+308
	//
	// bool to float64
	// true -> 1
	// false -> 0
	//
	// string to float64
	// "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.0000000000" -> 1.7976931348623157e+308
	// "1.7976900000e+308" -> 1.79769e+308
	// "1.7976931349e+308" -> +Inf
	// "Failed to convert from string to float64: fields=, value=\"1.1.1\", error=strconv.ParseFloat: parsing \"1.1.1\": invalid syntax"
}

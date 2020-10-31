package henge

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func ExampleValueConverter_String() {
	fmt.Println("int64 to string")
	fmt.Printf("%v -> %#v\n", math.MaxInt64, New(math.MaxInt64).String().Value())
	fmt.Printf("%v -> %#v\n", math.MinInt64, New(math.MinInt64).String().Value())
	fmt.Println()

	fmt.Println("uint64 to string")
	fmt.Printf("%v -> %#v\n", uint64(math.MaxUint64), New(uint64(math.MaxUint64)).String().Value())
	fmt.Println()

	fmt.Println("float64 to string")
	fmt.Printf("%v -> %#v\n", math.MaxFloat64, New(math.MaxFloat64).String().Value())
	fmt.Printf("%v -> %#v\n", math.MaxFloat64, New(math.MaxFloat64, WithFloatFormat('e', 2)).String().Value())
	fmt.Printf("%v -> %#v\n", math.SmallestNonzeroFloat64, New(math.SmallestNonzeroFloat64).String().Value())
	fmt.Printf("%v -> %#v\n", -1*math.MaxFloat64, New(-1*math.MaxFloat64).String().Value())
	fmt.Println()

	fmt.Println("bool to string")
	fmt.Printf("%v -> %#v\n", true, New(true).String().Value())
	fmt.Printf("%v -> %#v\n", false, New(false).String().Value())
	fmt.Println()

	fmt.Println("string to string")
	fmt.Printf("%v -> %v\n", "aaa", New("aaa").String().Value())

	// Output:
	// int64 to string
	// 9223372036854775807 -> "9223372036854775807"
	// -9223372036854775808 -> "-9223372036854775808"
	//
	// uint64 to string
	// 18446744073709551615 -> "18446744073709551615"
	//
	// float64 to string
	// 1.7976931348623157e+308 -> "179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	// 1.7976931348623157e+308 -> "1.80e+308"
	// 5e-324 -> "0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005"
	// -1.7976931348623157e+308 -> "-179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	//
	// bool to string
	// true -> "true"
	// false -> "false"
	//
	// string to string
	// aaa -> aaa
}

func TestStringConverterPtr(t *testing.T) {
	ptr, err := (&StringConverter{value: "", err: errors.New("error")}).Ptr().Result()
	if ptr != nil || err == nil {
		t.Error(ptr)
	}
	ptr, err = (&StringConverter{value: "", err: nil}).Ptr().Result()
	if ptr == nil || err != nil {
		t.Error(err)
	} else if *ptr != "" {
		t.Error(*ptr)
	}
	ptr, err = (&StringConverter{value: "aaa", err: nil}).Ptr().Result()
	if ptr == nil || err != nil {
		t.Error(err)
	} else if *ptr != "aaa" {
		t.Error(*ptr)
	}
}

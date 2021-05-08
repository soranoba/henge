package henge

import (
	"fmt"
	"math"
)

func ExampleValueConverter_JSONValue() {
	fmt.Println("int64")
	fmt.Printf("%v -> %v\n", math.MinInt64, New(math.MinInt64).JSONValue().Value())
	fmt.Printf("%v -> %v\n", 0, New(0).JSONValue().Value())
	fmt.Println()

	fmt.Println("uint64")
	fmt.Printf("%v -> %v\n", uint64(math.MaxUint64), New(uint64(math.MaxUint64)).JSONValue().Value())
	fmt.Printf("%v -> %v\n", 0, New(0).JSONValue().Value())
	fmt.Println()

	fmt.Println("float64")
	fmt.Printf("%v -> %v\n", -1*math.MaxFloat64, New(-1*math.MaxFloat64).JSONValue().Value())
	fmt.Printf("%v -> %v\n", 0.0, New(0.0).JSONValue().Value())
	fmt.Println()

	fmt.Println("bool")
	fmt.Printf("%v -> %v\n", true, New(true).JSONValue().Value())
	fmt.Printf("%v -> %v\n", false, New(false).JSONValue().Value())
	fmt.Println()

	fmt.Println("string")
	fmt.Printf("%v -> %v\n", "aaaa", New("aaaa").JSONValue().Value())
	fmt.Println()

	fmt.Println("null")
	fmt.Printf("%v -> %v\n", nil, New(nil).JSONValue().Value())
	fmt.Println()

	type Y struct {
		Z string
	}
	type X struct {
		Y
	}
	type Object struct {
		X
	}

	fmt.Println("object")
	fmt.Printf("%#v\n  -> %#v\n", Object{X: X{Y{Z: "z"}}}, New(Object{X: X{Y{Z: "z"}}}).JSONValue().Value())
	fmt.Printf(
		"%#v\n  -> %#v\n",
		map[interface{}]interface{}{"a": Object{X: X{Y{Z: "z"}}}},
		New(map[interface{}]interface{}{"a": Object{X: X{Y{Z: "z"}}}}).JSONValue().Value(),
	)
	fmt.Println()

	fmt.Println("array")
	fmt.Printf(
		"%#v\n  -> %#v\n",
		[]interface{}{Object{X: X{Y{Z: "z"}}}},
		New([]interface{}{Object{X: X{Y{Z: "z"}}}}).JSONValue().Value(),
	)

	// Output:
	// int64
	// -9223372036854775808 -> -9223372036854775808
	// 0 -> 0
	//
	// uint64
	// 18446744073709551615 -> 18446744073709551615
	// 0 -> 0
	//
	// float64
	// -1.7976931348623157e+308 -> -1.7976931348623157e+308
	// 0 -> 0
	//
	// bool
	// true -> true
	// false -> false
	//
	// string
	// aaaa -> aaaa
	//
	// null
	// <nil> -> <nil>
	//
	// object
	// henge.Object{X:henge.X{Y:henge.Y{Z:"z"}}}
	//   -> map[string]interface {}{"X":map[string]interface {}{"Y":map[string]interface {}{"Z":"z"}}}
	// map[interface {}]interface {}{"a":henge.Object{X:henge.X{Y:henge.Y{Z:"z"}}}}
	//   -> map[string]interface {}{"a":map[string]interface {}{"X":map[string]interface {}{"Y":map[string]interface {}{"Z":"z"}}}}
	//
	// array
	// []interface {}{henge.Object{X:henge.X{Y:henge.Y{Z:"z"}}}}
	//   -> []interface {}{map[string]interface {}{"X":map[string]interface {}{"Y":map[string]interface {}{"Z":"z"}}}}
}

func ExampleValueConverter_JSONObject() {
	type Y struct {
		Z string
	}
	type X struct {
		Y
	}
	type Object struct {
		X
	}

	value := map[interface{}]interface{}{"a": Object{X: X{Y{Z: "z"}}}}
	fmt.Printf("%#v\n  -> %#v\n", value, New(value).JSONValue().Value())

	// Output:
	// map[interface {}]interface {}{"a":henge.Object{X:henge.X{Y:henge.Y{Z:"z"}}}}
	//   -> map[string]interface {}{"a":map[string]interface {}{"X":map[string]interface {}{"Y":map[string]interface {}{"Z":"z"}}}}
}

func ExampleValueConverter_JSONArray() {
	type Y struct {
		Z string
	}
	type X struct {
		Y
	}
	type Object struct {
		X
	}

	value := []interface{}{
		1,
		true,
		"a",
		map[interface{}]interface{}{1: Y{Z: "z"}},
		Object{X: X{Y{Z: "z"}}},
	}
	for i, v := range New(value).JSONArray().Value() {
		fmt.Printf("[%v] %#v -> %#v\n", i, value[i], v)
	}

	// Output:
	// [0] 1 -> 1
	// [1] true -> true
	// [2] "a" -> "a"
	// [3] map[interface {}]interface {}{1:henge.Y{Z:"z"}} -> map[string]interface {}{"1":map[string]interface {}{"Z":"z"}}
	// [4] henge.Object{X:henge.X{Y:henge.Y{Z:"z"}}} -> map[string]interface {}{"X":map[string]interface {}{"Y":map[string]interface {}{"Z":"z"}}}
}

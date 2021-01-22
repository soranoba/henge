package henge

import (
	"fmt"
	"math"
	"reflect"
)

func ExampleWithFloatFormat() {
	fmt.Printf(
		"Default:                 %v\n",
		New(0.0125).String().Value(),
	)
	fmt.Printf(
		"WithFloatFormat('e', 2): %v\n",
		New(0.0125, WithFloatFormat('e', 2)).String().Value(),
	)

	// Output:
	// Default:                 0.0125
	// WithFloatFormat('e', 2): 1.25e-02
}

func ExampleWithRoundingFunc() {
	fmt.Println("int")
	fmt.Printf(
		"Default:                     %v\n",
		New(1.25).Int().Value(),
	)
	fmt.Printf(
		"WithRoundingFunc(math.Ceil): %v\n",
		New(1.25, WithRoundingFunc(math.Ceil)).Int().Value(),
	)
	fmt.Println()

	fmt.Println("uint")
	fmt.Printf(
		"Default:                     %v\n",
		New(1.25).Uint().Value(),
	)
	fmt.Printf(
		"WithRoundingFunc(math.Ceil): %v\n",
		New(1.25, WithRoundingFunc(math.Ceil)).Uint().Value(),
	)

	// Output:
	// int
	// Default:                     1
	// WithRoundingFunc(math.Ceil): 2
	//
	// uint
	// Default:                     1
	// WithRoundingFunc(math.Ceil): 2
}

func ExampleWithMapMaxDepth() {
	type Nested struct {
		Y string
	}
	type Value struct {
		Nested
		X string
	}
	in := map[string]Value{"a": {X: "a", Nested: Nested{Y: "y"}}}

	fmt.Printf(
		"Default:            %v\n",
		New(in).Map().Value(),
	)
	fmt.Printf(
		"WithMapMaxDepth(1): %v\n",
		New(in, WithMapMaxDepth(1)).Map().Value(),
	)

	// Output:
	// Default:            map[a:map[Nested:map[Y:y] X:a]]
	// WithMapMaxDepth(1): map[a:map[Nested:{y} X:a]]
}

func ExampleWithMapFilter() {
	type Value struct {
		X string
	}
	in := map[string]*Value{"a": {X: "a"}, "b": nil}

	fmt.Printf(
		"Default:                  %v\n",
		New(in).Map().Value(),
	)

	fmt.Printf(
		"Except when value is nil: %v\n",
		New(in, WithMapFilter(func(k interface{}, v interface{}) bool {
			r := reflect.ValueOf(v)
			return r.Kind() != reflect.Ptr || !r.IsNil()
		})).Map().Value(),
	)

	fmt.Printf(
		"Multiple Filters:         %v\n",
		New(in,
			WithMapFilter(func(k interface{}, v interface{}) bool {
				return true
			}),
			WithMapFilter(func(k interface{}, v interface{}) bool {
				return false
			}),
		).Map().Value(),
	)

	// Output:
	// Default:                  map[a:map[X:a] b:<nil>]
	// Except when value is nil: map[a:map[X:a]]
	// Multiple Filters:         map[]
}

func ExampleIsNil() {
	fmt.Printf("string: %v\n", isNil(""))
	fmt.Printf("*string: %v\n", isNil((*string)(nil)))
	fmt.Printf("nil: %v\n", isNil(nil))
	fmt.Printf("[]string: %v\n", isNil(([]string)(nil)))

	// Output:
	// string: false
	// *string: true
	// nil: true
	// []string: true
}

func ExampleIsZero() {
	fmt.Printf("string: %v\n", isZero(""))
	fmt.Printf("*string: %v\n", isZero((*string)(nil)))
	fmt.Printf("nil: %v\n", isZero(nil))
	fmt.Printf("[]string: %v\n", isZero(([]string)(nil)))

	// Output:
	// string: true
	// *string: true
	// nil: true
	// []string: true
}

func ExampleWithoutNilMapKey() {
	in := map[interface{}]interface{}{
		nil:            "a",
		(*string)(nil): "b",
		"":             "c",
	}

	fmt.Printf(
		"Default:                    %v\n",
		New(in).Map().Value(),
	)

	fmt.Printf(
		"Except when the key is nil: %v\n",
		New(in, WithoutNilMapKey()).Map().Value(),
	)

	// Output:
	// Default:                    map[<nil>:a <nil>:b :c]
	// Except when the key is nil: map[:c]
}

func ExampleWithoutNilMapValue() {
	in := map[interface{}]interface{}{
		"a": nil,
		"b": (*string)(nil),
		"c": "",
	}

	fmt.Printf(
		"Default:                      %v\n",
		New(in).Map().Value(),
	)

	fmt.Printf(
		"Except when the value is nil: %v\n",
		New(in, WithoutNilMapValue()).Map().Value(),
	)

	// Output:
	// Default:                      map[a:<nil> b:<nil> c:]
	// Except when the value is nil: map[c:]
}

func ExampleWithoutZeroMapKey() {
	in := map[interface{}]interface{}{
		nil:            "a",
		(*string)(nil): "b",
		"":             "c",
		"d":            "d",
	}

	fmt.Printf(
		"Default:                     %v\n",
		New(in).Map().Value(),
	)

	fmt.Printf(
		"Except when the key is zero: %v\n",
		New(in, WithoutZeroMapKey()).Map().Value(),
	)

	// Output:
	// Default:                     map[<nil>:a <nil>:b :c d:d]
	// Except when the key is zero: map[d:d]
}

func ExampleWithoutZeroMapValue() {
	in := map[interface{}]interface{}{
		"a": nil,
		"b": (*string)(nil),
		"c": "",
		"d": "d",
	}

	fmt.Printf(
		"Default:                       %v\n",
		New(in).Map().Value(),
	)

	fmt.Printf(
		"Except when the value is zero: %v\n",
		New(in, WithoutZeroMapValue()).Map().Value(),
	)

	// Output:
	// Default:                       map[a:<nil> b:<nil> c: d:d]
	// Except when the value is zero: map[d:d]
}

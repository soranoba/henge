package henge

import (
	"fmt"
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

func ExampleWithDepth() {
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

func ExampleWithFilter() {
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

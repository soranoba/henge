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

func ExampleWithMapKeyConverter() {
	in := map[interface{}]interface{}{
		"1.0": map[float64]interface{}{1.5: "a"},
		"2.0": map[uint64]interface{}{2: "b"},
	}

	fmt.Printf(
		"Default:            %v\n",
		New(in).Map().Value(),
	)

	fmt.Printf(
		"Convert key to int: %v\n",
		New(in, WithMapKeyConverter(func(keyConverter *ValueConverter) Converter {
			return keyConverter.Float().Int()
		})).Map().Value(),
	)

	// Output:
	// Default:            map[1.0:map[1.5:a] 2.0:map[2:b]]
	// Convert key to int: map[1:map[1:a] 2:map[2:b]]
}

func ExampleWithMapValueConverter() {
	in := map[interface{}]interface{}{
		"a": map[interface{}]interface{}{"a.1": 1.5, "a.2": 1},
		"b": map[interface{}]interface{}{"b.1": 2.5, "b.2": 2},
		"c": struct { X float64 }{X: 3.5},
	}

	fmt.Printf(
		"Default:                          %v\n",
		New(in).Map().Value(),
	)

	fmt.Printf(
		"Convert values of b and c to int: %v\n",
		New(in, WithMapValueConverter(func(key interface{}, keyConverter *ValueConverter) Converter {
			if key == "a.1" || key == "a.2" {
				return keyConverter
			}
			return keyConverter.Float().Int()
		})).Map().Value(),
	)

	// Output:
	// Default:                          map[a:map[a.1:1.5 a.2:1] b:map[b.1:2.5 b.2:2] c:map[X:3.5]]
	// Convert values of b and c to int: map[a:map[a.1:1.5 a.2:1] b:map[b.1:2 b.2:2] c:map[X:3]]
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

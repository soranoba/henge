package henge

import (
	"fmt"
	"reflect"
)

func ExampleGetStructFieldIndexes() {
	type Embedded2 struct {
		A string // [0 0 0]
	}
	type Embedded1 struct {
		*Embedded2        // [0 0]
		B          string // [0 1]
	}
	type Out struct {
		*Embedded1     // [0]
		A          int // [1]
	}

	fieldNames := getStructFieldIndexes(reflect.ValueOf(Out{}).Type())
	fmt.Println(fieldNames)

	// Output:
	// [[0] [0 0] [0 0 0] [0 1] [1]]
}

func ExampleGetStructFields() {
	type Embedded2 struct {
		A string `henge:"-"`
	}
	type Embedded1 struct {
		*Embedded2
		B string
	}
	type Out struct {
		*Embedded1 `henge:"-"`
		A          int
	}

	for _, field := range getStructFields(reflect.ValueOf(Out{}).Type()) {
		fmt.Println(field)
	}

	// Output:
	// {Embedded1 [0] [{true}]}
	// {Embedded2 [0 0] [{true} {false}]}
	// {A [0 0 0] [{true} {false} {true}]}
	// {B [0 1] [{true} {false}]}
	// {A [1] [{false}]}
}

package henge

import "fmt"

func ExampleValueConverter_Struct() {
	type In struct {
		A int
		B string
	}
	type Out struct {
		A string
		B *int
	}

	var in = In{
		A: 125,
		B: "25",
	}
	var out Out
	if err := New(in).Struct().Convert(&out); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("A = %#v\n", out.A)
		fmt.Printf("B = %#v\n", *out.B)
	}

	// Output:
	// A = "125"
	// B = 25
}

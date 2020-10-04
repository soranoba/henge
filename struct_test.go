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

	var s1 = In{
		A: 125,
		B: "25",
	}
	var s2 Out
	err := New(s1).Struct().Convert(&s2)
	if err != nil {
		return
	}

	fmt.Printf("A = %#v\n", s2.A)
	fmt.Printf("B = %#v\n", *s2.B)

	// Output:
	// A = "125"
	// B = 25
}

package henge

import "fmt"

func ExampleValueConverter_Map() {
	type In struct {
		A int
		B string
	}

	var in = In{
		A: 125,
		B: "25",
	}
	var s2 map[string]int
	err := New(in).Map().Convert(&s2)
	if err != nil {
		return
	}

	fmt.Printf("%#v\n", s2)

	// Output:
	// map[string]int{"A":125, "B":25}
}

func ExampleValueConverter_Map_Nested() {
	type Nested2 struct {
		Z int
	}
	type Nested1 struct {
		Nested2
		X string
		Y int
	}
	type In struct {
		A int
		B Nested1
	}

	in := In{A: 1, B: Nested1{X: "x", Y: 2, Nested2: Nested2{Z: 3}}}
	m, err := New(in).Map(WithMaxDepth(1)).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(m)
	}

	m, err = New(in).Map().Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(m)
	}

	// Output:
	// map[A:1 B:map[Nested2:{3} X:x Y:2]]
	// map[A:1 B:map[Nested2:map[Z:3] X:x Y:2]]
}

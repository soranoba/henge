package henge

import (
	"fmt"
)

func ExampleValueConverter_Slice() {
	var r1 []string
	fmt.Printf("%v | %#v -> %#v\n", New([...]int{1, 2, 3}).Slice().Convert(&r1), [...]int{1, 2, 3}, r1)
	var r2 [2]string
	fmt.Printf("%v | %#v -> %#v\n", New([]int{1, 2, 3}).Slice().Convert(&r2), []int{1, 2, 3}, r2)
	var r3 []int
	fmt.Printf("%v | %#v\n", New([]string{"1", "a"}).Slice().Convert(&r3), r3)

	type In struct {
		A string
	}
	type Out struct {
		A int
	}
	var out []Out
	if err := New([]In{{A: "123"}, {A: "234"}}).Convert(&out); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v -> %#v\n", []In{{A: "123"}, {A: "234"}}, out)
	}

	// Output:
	// <nil> | [3]int{1, 2, 3} -> []string{"1", "2", "3"}
	// <nil> | []int{1, 2, 3} -> [2]string{"1", "2"}
	// Failed to convert from string to int64: fields=, error=strconv.ParseInt: parsing "a": invalid syntax | []int(nil)
	// []henge.In{henge.In{A:"123"}, henge.In{A:"234"}} -> []henge.Out{henge.Out{A:123}, henge.Out{A:234}}
}

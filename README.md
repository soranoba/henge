# Henge
[![CircleCI](https://circleci.com/gh/soranoba/henge.svg?style=svg&circle-token=3c8c20a0a57a6333fb949dd6b901c610656e9da6)](https://circleci.com/gh/soranoba/henge)
[![Go Report Card](https://goreportcard.com/badge/github.com/soranoba/henge)](https://goreportcard.com/report/github.com/soranoba/henge)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/soranoba/henge)](https://pkg.go.dev/github.com/soranoba/henge)

Henge is a type conversion library for Golang.

変化 (Henge) means "Appearing with a different figure." in Japanese.<br>
Henge as the name implies can easily convert to different types.

## Overviews

- 💫　Easily converting to various types
  - int64, uint64, float64, bool, string, slice, map, and struct.
- ⚡　Simple and minimal code.
- 🔧　Support for custom conversions by callbacks before and after conversion.

## Motivations

### Easily converting pointer and non-pointer types.

In Golang world, there is a trade-off between pointer and non-pointer type, so it is used both as needed.<br>
For the reason, it is often necessary to convert between pointers and non-pointers.<br>
<br>
When using Henge, it can easily convert even if it's a struct field.

### Easily converting to a different struct.

There are many cases where the API server response is not the same as the DB record.<br>
Henge is very useful if you just want to copy, but want to ignore some fields.

### Easily create pointer-type values.

If we try to assign a non-zero constant value to a String or Int pointer type, we need to write codes of multiple lines.<br>
When using Henge, it easy to create pointer types while preserving the benefits of types.

## Installation

To install it, run:

```
go get -u github.com/soranoba/henge
```

## Usage

### Conversion to the specified type.

```go
import (
	"fmt"

	"github.com/soranoba/henge"
)

func main() {
	type UserRequest struct {
		Name *string
		Age  *int
	}
	type User struct {
		Name string // *string to string
		Age  string // *int to string
	}

	name := "Alice"
	age := 30
	var user User
	if err := henge.New(UserRequest{Name: &name, Age: &age}).Convert(&user); err != nil {
		return
	}
	fmt.Printf("%#v", user)
}
```

### Conversion by method chain.

```go
import (
	"fmt"

	"github.com/soranoba/henge"
)

func main() {
	i, err := henge.New("1.25").Float().Int().Result()
	if err != nil {
		return
  }
  // "1.25" -> 1.25 -> 1
	fmt.Println(i)
}
```

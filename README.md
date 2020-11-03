# Henge
[![CircleCI](https://circleci.com/gh/soranoba/henge.svg?style=svg&circle-token=3c8c20a0a57a6333fb949dd6b901c610656e9da6)](https://circleci.com/gh/soranoba/henge)
[![Go Report Card](https://goreportcard.com/badge/github.com/soranoba/henge)](https://goreportcard.com/report/github.com/soranoba/henge)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/soranoba/henge)](https://pkg.go.dev/github.com/soranoba/henge)

Henge is a type conversion library for Golang.

変化 (Henge) means "appear in a different appearance" in Japanese.<br>
Henge as the name implies can easily convert to different types.

## Overviews

- 💫　Easily converting to various types
  - int64, uint64, float64, bool, string, slice, map, and struct.
- ⚡　Simple and minimal code.
- 🔧　Support for custom conversions by callbacks before and after conversion.

## Motivations

### Easily converting pointer and non-pointer types.

There are several ways to handle null in Golang world, but there are tradeoffs in all cases.

Case 1. When it distinguish from Zero value by using structure like [sql.NullXX](https://golang.org/pkg/database/sql/), it will extra effort when using third-party library like [faker](https://github.com/bxcodec/faker) and [swag](https://github.com/swaggo/swag).<br>
Case 2. When using pointers, it will extra effort of conversion to non-pointer type.<br>

Henge aims to make the conversion easier for those who choose Case 2.

### Easily converting to a different struct.

There are many cases where the API server response is not the same as the DB record.<br>
Henge is very useful if you just want to copy, but want to ignore some fields.

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
	fmt.Println(i)
}
```

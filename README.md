# Henge
[![CircleCI](https://circleci.com/gh/soranoba/henge.svg?style=svg&circle-token=3c8c20a0a57a6333fb949dd6b901c610656e9da6)](https://circleci.com/gh/soranoba/henge)

Henge is a struct transrate library for Golang.  
变化 (Henge) means "appear in a different appearance" in Japanese.  

## Motivations

There are several ways to handle null in Golang world, but there are tradeoffs in all cases.  

Case 1. When it distinguish from Zero value by using structure like [sql.NullXX](https://golang.org/pkg/database/sql/),  it will extra effort when using third-party library like [faker](https://github.com/bxcodec/faker) and [swag](https://github.com/swaggo/swag).  

Case 2. When using pointers, it will extra effort of conversion to non-pointer type.  

`Henge` aims to eliminate the hassle of conversion in Case 2. 

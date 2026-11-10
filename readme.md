# Danom

Danom provides generic `Option`, `Result`, and `Choice` types, along with composable validation, for Go. It offers a type-safe way to represent optional values, expected errors, and values that may be one of two types, with functions for matching and transforming them.

## Key Features

- `option.Option[T]` for values that may be present or absent.
- `result.Result[T, E]` for success or error outcomes.
- `choice.Choice[T1, T2]` for values of either of two types.
- Functional helpers such as `Map`, `Bind`, and `Match`.
- Composable validation rules and built-in facts in the `check` package.

## Getting Started

Install the module:

```sh
go get github.com/eastcitysoftware/danom
```

Requires Go 1.25.0 or later.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/choice"
    "github.com/eastcitysoftware/danom/option"
    "github.com/eastcitysoftware/danom/result"
)

func divide(numerator, denominator int) result.Result[int, string] {
    if denominator == 0 {
        return result.Error[int, string]("Cannot divide by zero")
    }

    return result.Ok[int, string](numerator / denominator)
}

func main() {
    // Option
    doubled := option.Map(option.Some(5), func(x int) int {
        return x * 2
    })
    fmt.Println(option.Match(
        doubled,
        func(x int) string { return fmt.Sprintf("Value: %d", x) },
        func() string { return "No value" },
    ))

    // Result
    answer := result.Map(divide(10, 2), func(x int) int {
        return x + 1
    })
    fmt.Println(result.Match(
        answer,
        func(x int) string { return fmt.Sprintf("Result: %d", x) },
        func(err string) string { return fmt.Sprintf("Error: %s", err) },
    ))

    // Choice
    value := choice.OfT2[int, string]("hello")
    fmt.Println(choice.Match(
        value,
        func(x int) string { return fmt.Sprintf("Got int: %d", x) },
        func(s string) string { return fmt.Sprintf("Got string: %s", s) },
    ))
}
```

## Option

An `Option[T]` represents a value that may or may not be present. Create one with `option.Some(value)` or `option.None[T]()`. `Some` represents a present value independently of whether that value is nil.

Use `Value` when you want ordinary Go control flow, or `Match` to handle both cases:

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/option"
)

func main() {
    value, ok := option.Some(5).Value()
    if ok {
        fmt.Println(value)
    }

    description := option.Match(
        option.None[int](),
        func(value int) string { return fmt.Sprintf("Value: %d", value) },
        func() string { return "No value" },
    )
    fmt.Println(description)
}
```

`option.Map` transforms a present value. Use `option.Bind` when the transformation itself returns an `Option`. The package also provides `Apply`, `OrElse`, `OrElseWith`, `Tap`, `Iter`, `Filter`, and `Flatten`, as well as `DefaultValue` and `DefaultWith` methods.

## Result

A `Result[T, E]` represents either a value of type `T` or an error value of type `E`. Create one with `result.Ok[T, E](value)` or `result.Error[T, E](err)`. The error type is chosen by the caller and is not required to implement Go's `error` interface.

Use `Value` and `ErrorValue` for procedural handling, or `Match` to provide a handler for each case:

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/result"
)

func main() {
    outcome := result.Ok[int, string](42)
    if value, ok := outcome.Value(); ok {
        fmt.Printf("Value: %d\n", value)
    } else if err, ok := outcome.ErrorValue(); ok {
        fmt.Printf("Error: %s\n", err)
    }
}
```

`result.Map` transforms a successful value, while `result.Bind` chains an operation that returns another `Result`. Use `result.MapError` to transform the error value. Other helpers include `Apply`, `OrElse`, `OrElseWith`, `Tap`, and `Iter`; `Result` also provides `DefaultValue` and `DefaultWith` methods.

## Choice

A `Choice[T1, T2]` holds a value of either `T1` or `T2`. Use `choice.OfT1[T1, T2](value)` or `choice.OfT2[T1, T2](value)` to construct one:

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/choice"
)

func main() {
    value := choice.OfT1[int, string](42)
    message := choice.Match(
        value,
        func(n int) string { return fmt.Sprintf("Got int: %d", n) },
        func(s string) string { return fmt.Sprintf("Got string: %s", s) },
    )
    fmt.Println(message)
}
```

`IsT1` and `IsT2` inspect which type is held; `T1` and `T2` retrieve it. `MapT1` and `MapT2` transform one side, while `BindT1` and `BindT2` chain a function that returns another `Choice`. The package also provides `Do` and `Swap`.

## Validation

The `check` package builds validators from named rules. A rule selects a value from the input and applies one or more facts to it. A validator evaluates its rules and returns the original input on success, or the collected validation errors on failure.

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/check"
)

type User struct {
    Name string
    Age  int
}

func main() {
    validate := check.NewValidator(
        check.NewRule("name", func(user User) string { return user.Name },
            check.StrNotEmpty(),
            check.StrLengthBetween(2, 40)),
        check.NewRule("age", func(user User) int { return user.Age },
            check.GreaterThanOrEqual(18)),
    )

    outcome := validate.Validate(User{Name: "", Age: 16})
    if user, ok := outcome.Value(); ok {
        fmt.Printf("Valid user: %+v\n", user)
    } else if errs, ok := outcome.ErrorValue(); ok {
        fmt.Println(errs)
    }
}
```

Built-in facts cover equality, booleans, numeric and ordered values, strings, slices, maps, optional values, times, and durations. Facts can be combined with `When`, `AnyOf`, and `NoneOf`. Create a custom fact with `check.NewFact`, and customize its message with `check.NewMessage`. Validation errors are available as `check.Errors`, which implements Go's `error` interface.

## Find a bug?

There's an [issue](https://github.com/eastcitysoftware/danom/issues) for that.

## License

Licensed under [MIT](https://github.com/eastcitysoftware/danom/blob/master/LICENSE).
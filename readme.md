# danom

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/eastcitysoftware/danom/v4)
[![Latest Release](https://img.shields.io/github/v/release/eastcitysoftware/danom?style=flat-square)](https://github.com/eastcitysoftware/danom/releases/latest)
[![Build Status](https://github.com/eastcitysoftware/danom/actions/workflows/build.yml/badge.svg)](https://github.com/eastcitysoftware/danom/actions/workflows/build.yml)

danom provides generic `Option` and `Result` types, string conversion helpers, and composable validation for Go, inspired by F#. These discriminated unions offer a type-safe way to represent nullable values, expected errors and decisions, providing an API (e.g., Map, Bind etc.) for value transformations and safe data access.

## Key Features

- `option.Option[T]` makes the presence or absence of a value explicit.
- `result.Result[T, E]` represents a successful value or a typed error.
- `parse` returns standard-library parsing results as `Result` values.
- `convert` exposes the same conversions as `parse` but produces `Option` values.
- `check` provides reusable validation rules, composable facts, and collected errors.

## Design Goals

- Enhance functional programming in Go.
- Support robust use through opinionated monads.
- Provide accumulative input validation.

## Getting Started

Install the module with Go:

```sh
go get github.com/eastcitysoftware/danom
```

danom requires Go 1.21.0 or later. Import only the packages you need: `option`, `result`, `parse`, `convert`, and `check` are separate packages. See the [API reference](reference.md) for the complete list of types, functions, and built-in validation facts.

## Option

Use `Option[T]` when a value may legitimately be absent, for example, a lookup that might not find a record. This keeps absence distinct from a type's zero value. `option.Some(value)` represents a present value even when `value` is nil; `option.None[T]()` represents no value. The zero value of `Option[T]` is also `None`.

Use `Value` when ordinary Go control flow is clearest, or `Match` when both cases should produce the same result type:

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/option"
)

func doubler(n int) int {
    return n * 2
}

func main() {
    value := option.Some(5)
    doubled := option.Map(value, doubler)
    fmt.Println(doubled.DefaultValue(0))

    value = option.None[int]()
    doubled = option.Map(value, doubler)
    fmt.Println(doubled.DefaultValue(86))

    description := option.Match(
        option.None[int](),
        func(n int) string { return fmt.Sprintf("Value: %d", n) },
        func() string { return "No value" },
    )
    fmt.Println(description)

    number := option.Some(86).
        Filter(func(n int) bool { return n > 0 }).
        Tap(func(n int) { fmt.Println("Value:", n) }).
        DefaultValue(0)
    fmt.Println("Number:", number)
}
```

`option.Map` transforms a present value and leaves `None` absent. Use `option.Bind` when the transformation already returns an `Option`, so you don't end up with a nested `Option[Option[T]]`. `Filter` keeps a present value only when a predicate succeeds; `Flatten` removes one layer of nesting.

When you need a fallback, `DefaultValue` uses a value you already have, while `DefaultWith` computes one only if the option is `None`. `OrElse` and `OrElseWith` provide the corresponding fallback as another `Option`. `Tap` and `Iter` let you perform an action for a present value without changing the value.

These methods can be chained to express a small pipeline. `Filter` can turn a present value into `None`, `Tap` runs only while a value is present, and `DefaultValue` unwraps the result with a fallback:

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/option"
)

type User struct {
    Name   string
    Active bool
}

func isActive(user User) bool {
    return user.Active
}

func logUser(user User) {
    fmt.Printf("active user: %s\n", user.Name)
}

func main() {
    anonymousUser := User{Name: "anonymous"}
    user := option.Some(User{Name: "Ada", Active: true}).
        Filter(isActive).
        Tap(logUser).
        DefaultValue(anonymousUser)
    fmt.Printf("user: %+v\n", user)
}
```


`Option[T]` supports Go's `encoding/json`: `Some(value)` is encoded as the JSON value, and `None` as `null`. Decoding `null` produces `None`. Consequently, `Some(nil)` also encodes as `null` and decodes as `None`, so that distinction cannot be preserved through JSON. On Go 1.24 and later, use the `omitzero` tag to omit `None` instead of encoding it as `null`:

```go
import "github.com/eastcitysoftware/danom/option"

type User struct {
    Nickname option.Option[string] `json:"nickname,omitzero"`
}
```

`IsZero` reports true only for `None`, so `Some` containing a zero value is retained. Go 1.21 through 1.23 do not recognize the `omitzero` tag; on those versions, `None` is emitted as `null`.

## Result

Use `Result[T, E]` when an operation can fail in an expected way and callers should handle that outcome. `Ok` carries a value of type `T`; `Error` carries a value of type `E`. `E` can be a Go `error`, a string, or another type suited to the operation. The zero value of `Result[T, E]` is an `Error` containing the zero value of `E`.

For example, a parser can return a `Result` so its caller must choose how to handle success and failure:

```go
package main

import (
    "errors"
    "fmt"
    "strconv"

    "github.com/eastcitysoftware/danom/result"
)

func parsePort(raw string) result.Result[int, error] {
    port, err := strconv.Atoi(raw)
    if err != nil {
        return result.Error[int, error](fmt.Errorf("parse port: %w", err))
    }
    return result.Ok[int, error](port)
}

func main() {
    port := parsePort("8080").
        Filter(
            func(port int) bool { return port >= 1 && port <= 65535 },
            errors.New("port must be between 1 and 65535"),
        ).
        Tap(func(port int) {
            fmt.Printf("Validated port: %d\n", port)
        })

    message := result.Match(
        port,
        func(port int) string {
            return fmt.Sprintf("Listening on port %d", port)
        },
        func(err error) string {
            return fmt.Sprintf("Invalid port: %v", err)
        })
    fmt.Println(message)
}
```

`Value` and `ErrorValue` return a value and a boolean for conventional branching; `Match` and `Do` handle both cases through callbacks. `result.Map` transforms only a success, and `result.Bind` chains an operation that can also fail. `MapError` transforms only the error, which is useful when adapting one error type to another.

`OrElse` and `OrElseWith` provide alternate results on failure; `DefaultValue` and `DefaultWith` extract a success or use a fallback value. Use these when discarding or recovering from the error is intentional. `Result` does not replace Go's `error` interface: use `E = error` when you want standard error wrapping and inspection.

## Parsing and conversion

The `parse` package wraps Go's standard parsing functions for booleans, numeric and complex values, durations, and times. Each function returns a `result.Result[T, error]`, preserving the parse error when the input is invalid. For example:

```go
package main

import (
    "errors"
    "fmt"

    "github.com/eastcitysoftware/danom/convert"
    "github.com/eastcitysoftware/danom/parse"
    "github.com/eastcitysoftware/danom/result"
)

func main() {
    port := parse.Int("8080").
        Filter(func(port int) bool { return port >= 1 && port <= 65535 },
            errors.New("port must be between 1 and 65535")).
        Tap(func(port int) { fmt.Println("validated port:", port) })
    result.Do(port,
        func(n int) { fmt.Println(n) },
        func(e error) { fmt.Println(e) })

    retries := convert.Int("not a number").
        Filter(func(n int) bool { return n >= 0 }).
        DefaultValue(0)
    fmt.Println("retry count:", retries)
}
```

The `convert` package provides matching functions that return `option.Option[T]` instead. A valid input produces `Some(value)`; invalid input produces `None`, which is useful when the parse error itself is not needed, as in the retry-count fallback above:

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/convert"
)

func main() {
    port := convert.Uint16("8080")
    fmt.Println(port.DefaultValue(0))
}
```

Both packages provide functions for `Bool`, `Float32`, `Float64`, `Int8`, `Int16`, `Int32`, `Int64`, `Uint8`, `Uint16`, `Uint32`, `Uint64`, `Complex64`, `Complex128`, `Duration`, `Time`, and `TimeRFC3339`. `Time` accepts a Go time layout; `TimeRFC3339` uses the RFC3339 layout.

## Validation

The `check` package builds reusable validators from named rules. Each rule selects a field or derived value from the input, then applies one or more facts to it. A validator evaluates every configured rule, collects all failures, and returns either the original input or the collected errors as a `result.Result`. This makes it possible to report multiple field errors in one pass rather than stopping at the first one.

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/check"
    "github.com/eastcitysoftware/danom/option"
)

type User struct {
    Name     string
    Age      int
    Nickname option.Option[string]
}

func main() {
    validate := check.NewValidator(
        check.NewRule("name", func(user User) string { return user.Name },
            check.StrNotEmpty(),
            check.StrLengthBetween(2, 40)),
        check.NewRule("age", func(user User) int { return user.Age },
            check.Between(18, 120)),
        check.NewRule("nickname",
            func(user User) option.Option[string] { return user.Nickname },
            check.Optional(check.StrLengthBetween(2, 20))),
    )

    input := User{Name: "", Age: 16, Nickname: option.Some("A")}
    outcome := validate.Validate(input)
    if user, ok := outcome.Value(); ok {
        fmt.Printf("Valid user: %+v\n", user)
    } else if errs, ok := outcome.ErrorValue(); ok {
        fmt.Println(errs)
    }
}
```

Create a `Rule` with `NewRule(name, selector, facts...)`, then pass one or more rules to `NewValidator`. Facts attached to a rule are all evaluated, so a field can report multiple reasons it is invalid. Use `When` for conditional checks, `AnyOf` when at least one fact must pass, and `NoneOf` when all listed conditions must fail.

Built-in facts cover equality, booleans, numeric and ordered values, strings, slices, maps, optional values, times, and durations. `MapAll` and `MapAny` apply a predicate to map values, requiring respectively every value or at least one value to pass. For optional fields, `Required` checks that an `Option` is present and validates its value; `Optional` validates it only when present. To define a project-specific rule, use `NewFact` with a predicate and a `NewMessage` format string; the rule name is supplied as the format's first argument.

### Custom validation facts

Domain-specific validators—email addresses, phone numbers, or any complex format—are best owned by the code that uses them. Use `NewFact` to wrap validation logic into a composable fact:

```go
package main

import (
    "net"
    "github.com/eastcitysoftware/danom/check"
)

type User struct {
    IP string
}

func main() {
    // Wrap net.ParseIP to validate IPv4 addresses strictly
    isValidIPv4 := check.NewFact(
        func(s string) bool {
            ip := net.ParseIP(s)
            return ip != nil && ip.To4() != nil
        },
        check.NewMessage("'%s' must be a valid IPv4 address"))

    validate := check.NewValidator(
        check.NewRule("ip",
            func(u User) string { return u.IP },
            isValidIPv4))

    outcome := validate.Validate(User{IP: "192.168.1.1"})
    if _, ok := outcome.Value(); ok {
        // Valid
    }
}
```

This pattern scales to any third-party validator: `regexp.Compile` for patterns, `libphonenumber` for E.164 numbers, or custom domain logic.

On failure, `Errors` contains named `Error` values and implements Go's `error` interface, so it can be printed, returned, or handled like other errors. On success, the result contains the original input.

## Find a bug?

There's an [issue](https://github.com/eastcitysoftware/danom/issues) for that.

## License

Licensed under [MIT](https://github.com/eastcitysoftware/danom/blob/main/LICENSE).
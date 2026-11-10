# danom

[![build](https://github.com/eastcitysoftware/danom/actions/workflows/build.yml/badge.svg)](https://github.com/eastcitysoftware/danom/actions/workflows/build.yml)

danom provides generic `Option` and `Result` types, along with composable validation, for Go. It offers a type-safe way to represent optional values and expected errors, with functions for matching and transforming them.

## Key Features

- `option.Option[T]` makes the presence or absence of a value explicit.
- `result.Result[T, E]` represents a successful value or a typed error.
- `Map`, `Bind`, `Match`, and related helpers compose transformations and handle each case.
- `check` provides reusable validation rules, composable facts, and collected errors.

## Getting Started

Install the module with Go:

```sh
go get github.com/eastcitysoftware/danom
```

danom requires Go 1.25.0 or later. Import only the packages you need: `option`, `result`, and `check` are separate packages. See the [API reference](api.md) for the complete list of types, functions, and built-in validation facts.

## Option

Use `Option[T]` when a value may legitimately be absent, for example, a lookup that might not find a record. This keeps absence distinct from a type's zero value. `option.Some(value)` represents a present value even when `value` is nil; `option.None[T]()` represents no value. The zero value of `Option[T]` is also `None`.

Use `Value` when ordinary Go control flow is clearest, or `Match` when both cases should produce the same result type:

```go
package main

import (
    "fmt"

    "github.com/eastcitysoftware/danom/option"
)

func main() {
    value := option.Some(5)
    doubled := option.Map(value, func(n int) int { return n * 2 })
    fmt.Println(doubled.DefaultValue(0))

    description := option.Match(
        option.None[int](),
        func(n int) string { return fmt.Sprintf("Value: %d", n) },
        func() string { return "No value" },
    )
    fmt.Println(description)
}
```

`option.Map` transforms a present value and leaves `None` absent. Use `option.Bind` when the transformation already returns an `Option`, so you don't end up with a nested `Option[Option[T]]`. `Filter` keeps a present value only when a predicate succeeds; `Flatten` removes one layer of nesting.

When you need a fallback, `DefaultValue` uses a value you already have, while `DefaultWith` computes one only if the option is `None`. `OrElse` and `OrElseWith` provide the corresponding fallback as another `Option`. `Tap` and `Iter` let you perform an action for a present value without changing the value.

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
    if port < 1 || port > 65535 {
        return result.Error[int, error](errors.New("port must be between 1 and 65535"))
    }
    return result.Ok[int, error](port)
}

func main() {
    message := result.Match(
        parsePort("8080"),
        func(port int) string { return fmt.Sprintf("Listening on port %d", port) },
        func(err error) string { return fmt.Sprintf("Invalid port: %v", err) },
    )
    fmt.Println(message)
}
```

`Value` and `ErrorValue` return a value and a boolean for conventional branching; `Match` and `Do` handle both cases through callbacks. `result.Map` transforms only a success, and `result.Bind` chains an operation that can also fail. `MapError` transforms only the error, which is useful when adapting one error type to another.

`OrElse` and `OrElseWith` provide alternate results on failure; `DefaultValue` and `DefaultWith` extract a success or use a fallback value. Use these when discarding or recovering from the error is intentional. `Result` does not replace Go's `error` interface: use `E = error` when you want standard error wrapping and inspection.

## Validation

The `check` package builds reusable validators from named rules. Each rule selects a field or derived value from the input, then applies one or more facts to it. A validator evaluates every configured rule, collects all failures, and returns either the original input or the collected errors as a `result.Result`. This makes it possible to report multiple field errors in one pass rather than stopping at the first one.

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
            check.Between(18, 120)),
    )

    input := User{Name: "", Age: 16}
    outcome := validate.Validate(input)
    if user, ok := outcome.Value(); ok {
        fmt.Printf("Valid user: %+v\n", user)
    } else if errs, ok := outcome.ErrorValue(); ok {
        fmt.Println(errs)
    }
}
```

Create a `Rule` with `NewRule(name, selector, facts...)`, then pass one or more rules to `NewValidator`. Facts attached to a rule are all evaluated, so a field can report multiple reasons it is invalid. Use `When` for conditional checks, `AnyOf` when at least one fact must pass, and `NoneOf` when all listed conditions must fail.

Built-in facts cover equality, booleans, numeric and ordered values, strings, slices, maps, optional values, times, and durations. For optional fields, `Required` checks that an `Option` is present and validates its value; `Optional` validates it only when present. To define a project-specific rule, use `NewFact` with a predicate and a `NewMessage` format string; the rule name is supplied as the format's first argument.

On failure, `Errors` contains named `Error` values and implements Go's `error` interface, so it can be printed, returned, or handled like other errors. On success, the result contains the original input.

## Find a bug?

There's an [issue](https://github.com/eastcitysoftware/danom/issues) for that.

## License

Licensed under [MIT](https://github.com/eastcitysoftware/danom/blob/master/LICENSE).
# Danom API Reference

Danom is organized into packages for optional values, success-or-error results, two-way choices, and validation.

## Packages

- [`option`](#option) — optional values.
- [`result`](#result) — success or error outcomes.
- [`choice`](#choice) — a value of one of two types.
- [`check`](#check) — validation rules, facts, and errors.

## `option`

Import path: `github.com/eastcitysoftware/danom/option`

### Types

```go
type Option[T any] struct { /* unexported fields */ }
```

`Option[T]` represents a value that may or may not be present. Its zero value is `None`.

#### Methods

- `func (o Option[T]) IsSome() bool` reports whether `o` contains a value.
- `func (o Option[T]) IsNone() bool` reports whether `o` has no value.
- `func (o Option[T]) Value() (T, bool)` returns the value and `true` when `o` is `Some`; otherwise returns the zero value and `false`.
- `func (o Option[T]) DefaultValue(d T) T` returns the value or `d` when `o` is `None`.
- `func (o Option[T]) DefaultWith(fn func() T) T` returns the value or calls `fn` when `o` is `None`.

### Functions

- `func Some[T any](val T) Option[T]` creates an `Option` containing `val`.
- `func None[T any]() Option[T]` creates an `Option` with no value.
- `func IsSome[T any](o Option[T]) bool` reports whether `o` contains a value.
- `func IsNone[T any](o Option[T]) bool` reports whether `o` has no value.
- `func Match[T any, U any](o Option[T], some func(T) U, none func() U) U` calls the handler for the state of `o` and returns its result.
- `func Do[T any](o Option[T], some func(T), none func())` calls the handler for the state of `o`.
- `func Bind[T any, U any](o Option[T], fn func(T) Option[U]) Option[U]` applies `fn` to the value, or returns `None`.
- `func Map[T any, U any](o Option[T], fn func(T) U) Option[U]` transforms the value, or returns `None`.
- `func Apply[T any, U any](o Option[T], fn Option[func(T) U]) Option[U]` applies an optional function to an optional value.
- `func OrElse[T any](o Option[T], ifNone Option[T]) Option[T]` returns `o` when it is `Some`, or `ifNone` otherwise.
- `func OrElseWith[T any](o Option[T], ifNone func() Option[T]) Option[T]` returns `o` when it is `Some`, or calls `ifNone` otherwise.
- `func Tap[T any](o Option[T], fn func(T)) Option[T]` calls `fn` for a present value and returns `o`.
- `func Iter[T any](o Option[T], fn func(T))` calls `fn` when `o` is `Some`.
- `func Filter[T any](o Option[T], pred func(T) bool) Option[T]` returns `o` only when its value satisfies `pred`.
- `func Flatten[T any](o Option[Option[T]]) Option[T]` converts a nested `Option` into one `Option`.

## `result`

Import path: `github.com/eastcitysoftware/danom/result`

### Types

```go
type Result[T any, E any] struct { /* unexported fields */ }
```

`Result[T, E]` represents either a value of type `T` or an error value of type `E`. `E` can be any type; it does not have to implement Go's `error` interface. The zero value is an `Error` containing the zero value of `E`.

#### Methods

- `func (r Result[T, E]) IsOk() bool` reports whether `r` contains a success value.
- `func (r Result[T, E]) IsError() bool` reports whether `r` contains an error.
- `func (r Result[T, E]) Value() (T, bool)` returns the success value and `true` when `r` is `Ok`; otherwise returns the zero value and `false`.
- `func (r Result[T, E]) ErrorValue() (E, bool)` returns the error value and `true` when `r` is an `Error`; otherwise returns the zero value and `false`.
- `func (r Result[T, E]) DefaultValue(d T) T` returns the success value or `d` when `r` is an `Error`.
- `func (r Result[T, E]) DefaultWith(fn func() T) T` returns the success value or calls `fn` when `r` is an `Error`.

### Functions

- `func Ok[T any, E any](t T) Result[T, E]` creates a successful `Result`.
- `func Error[T any, E any](e E) Result[T, E]` creates an error `Result`.
- `func Match[T any, E any, U any](r Result[T, E], ok func(T) U, error func(E) U) U` calls the handler for the state of `r` and returns its result.
- `func Do[T any, E any](r Result[T, E], ok func(T), error func(E))` calls the handler for the state of `r`.
- `func Bind[T any, E any, U any](r Result[T, E], fn func(T) Result[U, E]) Result[U, E]` applies `fn` to the success value, or preserves the error.
- `func Map[T any, E any, U any](r Result[T, E], fn func(T) U) Result[U, E]` transforms the success value, or preserves the error.
- `func MapError[T any, E any, U any](r Result[T, E], fn func(E) U) Result[T, U]` transforms the error value, or preserves the success value.
- `func Apply[T any, E any, U any](r Result[T, E], fn Result[func(T) U, E]) Result[U, E]` applies a successful function to a successful value.
- `func OrElse[T any, E any](r Result[T, E], ifError Result[T, E]) Result[T, E]` returns `r` when it is `Ok`, or `ifError` otherwise.
- `func OrElseWith[T any, E any](r Result[T, E], ifError func() Result[T, E]) Result[T, E]` returns `r` when it is `Ok`, or calls `ifError` otherwise.
- `func Tap[T any, E any](r Result[T, E], fn func(T)) Result[T, E]` calls `fn` for a success value and returns `r`.
- `func Iter[T any, E any](r Result[T, E], fn func(T))` calls `fn` when `r` is `Ok`.

## `choice`

Import path: `github.com/eastcitysoftware/danom/choice`

### Types

```go
type Choice[T1 any, T2 any] struct { /* unexported fields */ }
```

`Choice[T1, T2]` holds a value of type `T1` or `T2`. The zero value holds the zero value of `T1`.

#### Methods

- `func (c Choice[T1, T2]) IsT1() bool` reports whether `c` holds a `T1`.
- `func (c Choice[T1, T2]) IsT2() bool` reports whether `c` holds a `T2`.
- `func (c Choice[T1, T2]) TryGetT1() (T1, bool)` returns the `T1` value and `true` when `c` holds `T1`; otherwise returns the zero value and `false`.
- `func (c Choice[T1, T2]) TryGetT2() (T2, bool)` returns the `T2` value and `true` when `c` holds `T2`; otherwise returns the zero value and `false`.
- `func (c Choice[T1, T2]) Swap() Choice[T2, T1]` returns a choice with the type parameters and held value reversed.

### Functions

- `func OfT1[T1 any, T2 any](t T1) Choice[T1, T2]` creates a `Choice` holding a `T1`.
- `func OfT2[T1 any, T2 any](t T2) Choice[T1, T2]` creates a `Choice` holding a `T2`.
- `func Match[T1 any, T2, U any](c Choice[T1, T2], t1Fn func(T1) U, t2Fn func(T2) U) U` calls the handler for the held type and returns its result.
- `func Do[T1 any, T2 any](c Choice[T1, T2], t1Fn func(T1), t2Fn func(T2))` calls the handler for the held type.
- `func BindT1[T1 any, T2 any, U any](c Choice[T1, T2], fn func(T1) Choice[U, T2]) Choice[U, T2]` applies `fn` to a `T1`, leaving a `T2` unchanged.
- `func MapT1[T1 any, T2 any, U any](c Choice[T1, T2], fn func(T1) U) Choice[U, T2]` transforms a `T1`, leaving a `T2` unchanged.
- `func BindT2[T1 any, T2 any, U any](c Choice[T1, T2], fn func(T2) Choice[T1, U]) Choice[T1, U]` applies `fn` to a `T2`, leaving a `T1` unchanged.
- `func MapT2[T1 any, T2 any, U any](c Choice[T1, T2], fn func(T2) U) Choice[T1, U]` transforms a `T2`, leaving a `T1` unchanged.

## `check`

Import path: `github.com/eastcitysoftware/danom/check`

Validation rules select values from an input and evaluate facts against them. A validator evaluates every configured rule, collects all errors, and returns either the original input or the collected errors as a `result.Result`.

### Types

- `type Error struct { Name string; Errors []string }` contains messages for one named rule. Its `String` method formats the rule name and messages.
- `type Errors []Error` is a collection of validation errors. It implements `error`; `String` formats the collection.
- `type Message func(string) string` formats a message using the name of the value being checked.
- `type Fact[T any] func(T) []Message` checks a value and returns messages for failed conditions.
- `type Rule[T any] func(T) Error` validates an input and returns a named error.
- `type Validator[T any] func(T) result.Result[T, Errors]` validates an input and returns it or its errors.
- `type Number interface { ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 }` is the type constraint used by numeric facts.

#### Methods

- `func (v Validator[T]) Validate(input T) result.Result[T, Errors]` applies the validator to `input`.
- `func (e Error) String() string` formats one rule error.
- `func (e Errors) String() string` formats the collection of rule errors.
- `func (e Errors) Error() string` implements Go's `error` interface.

### Rule and message constructors

- `func NewValidator[T any](rules ...Rule[T]) Validator[T]` creates a validator that evaluates the rules and collects their errors.
- `func NewRule[T, U any](name string, selector func(T) U, preds ...Fact[U]) Rule[T]` creates a named rule that selects a value and checks it against facts.
- `func NewFact[T any](pred func(T) bool, message Message) Fact[T]` creates a fact that returns `message` when `pred` is false.
- `func NewMessage(format string, args ...any) Message` creates a `Message` using `fmt.Sprintf`, with the rule name as its first argument.

### Fact combinators

- `func When[T any](cond Fact[T], facts ...Fact[T]) Fact[T]` evaluates `facts` when `cond` succeeds and returns their messages.
- `func AnyOf[T any](facts ...Fact[T]) Fact[T]` succeeds when at least one fact succeeds.
- `func NoneOf[T any](facts ...Fact[T]) Fact[T]` succeeds when every fact fails.

### Equality and boolean facts

- `func Equals[C comparable](e C) Fact[C]` succeeds when a value equals `e`.
- `func NotEquals[C comparable](e C) Fact[C]` succeeds when a value does not equal `e`.
- `func True() Fact[bool]` succeeds when the value is true.
- `func False() Fact[bool]` succeeds when the value is false.

### Numeric facts

These functions use the `Number` constraint.

- `func Zero[N Number]() Fact[N]` succeeds when a value is zero.
- `func Positive[N Number]() Fact[N]` succeeds when a value is greater than zero.
- `func Negative[N Number]() Fact[N]` succeeds when a value is less than zero.
- `func NonZero[N Number]() Fact[N]` succeeds when a value is not zero.
- `func NonPositive[N Number]() Fact[N]` succeeds when a value is at most zero.
- `func NonNegative[N Number]() Fact[N]` succeeds when a value is at least zero.

### Ordered-value facts

These functions accept types supported by `cmp.Ordered`.

- `func Between[N cmp.Ordered](min, max N) Fact[N]` succeeds when a value is in the inclusive range `[min, max]`.
- `func NotBetween[N cmp.Ordered](min, max N) Fact[N]` succeeds when a value is outside the inclusive range `[min, max]`.
- `func GreaterThan[N cmp.Ordered](min N) Fact[N]` succeeds when a value is greater than `min`.
- `func GreaterThanOrEqual[N cmp.Ordered](minimum N) Fact[N]` succeeds when a value is at least `minimum`.
- `func LessThan[N cmp.Ordered](min N) Fact[N]` succeeds when a value is less than `min`.
- `func LessThanOrEqual[N cmp.Ordered](minimum N) Fact[N]` succeeds when a value is at most `minimum`.

### String facts

- `func StrEmpty() Fact[string]` succeeds when a string is empty or contains only whitespace.
- `func StrNotEmpty() Fact[string]` succeeds when a string contains a non-whitespace character.
- `func StrEqualsFold(s string) Fact[string]` succeeds when a string equals `s`, ignoring case.
- `func StrStartsWith(s string) Fact[string]` succeeds when a string starts with `s`.
- `func StrEndsWith(s string) Fact[string]` succeeds when a string ends with `s`.
- `func StrContains(s string) Fact[string]` succeeds when a string contains `s`.
- `func StrNotContains(s string) Fact[string]` succeeds when a string does not contain `s`.
- `func StrLength(n int) Fact[string]` succeeds when a string contains exactly `n` Unicode code points.
- `func StrLengthBetween(min, max int) Fact[string]` succeeds when a string contains between `min` and `max` Unicode code points, inclusive.
- `func StrGreaterThan(min int) Fact[string]` succeeds when a string contains more than `min` Unicode code points.
- `func StrGreaterThanOrEqual(min int) Fact[string]` succeeds when a string contains at least `min` Unicode code points.
- `func StrLessThan(max int) Fact[string]` succeeds when a string contains fewer than `max` Unicode code points.
- `func StrLessThanOrEqual(max int) Fact[string]` succeeds when a string contains at most `max` Unicode code points.
- `func StrMatch(r *regexp.Regexp) Fact[string]` succeeds when `r` matches the string.
- `func StrAlpha() Fact[string]` succeeds when a string contains only ASCII letters.
- `func StrAlphanumeric() Fact[string]` succeeds when a string contains only ASCII letters and digits.
- `func StrIsUrl() Fact[string]` succeeds when a string is an HTTP or HTTPS URL with a host.
- `func StrIsEmail() Fact[string]` succeeds when a string matches the package's email pattern.
- `func StrIsE164() Fact[string]` succeeds when a string matches the package's E.164 phone number pattern.
- `func StrIsIPV4() Fact[string]` succeeds when a string parses as an IPv4 address.
- `func StrIsIPV6() Fact[string]` succeeds when a string parses as an IPv6 address.

### Slice facts

- `func Empty[T any]() Fact[[]T]` succeeds when a slice has no elements.
- `func NotEmpty[T any]() Fact[[]T]` succeeds when a slice has at least one element.
- `func Length[T any](n int) Fact[[]T]` succeeds when a slice contains exactly `n` elements.
- `func LengthBetween[T any](min, max int) Fact[[]T]` succeeds when a slice contains between `min` and `max` elements, inclusive.
- `func Min[T any](n int) Fact[[]T]` succeeds when a slice contains at least `n` elements.
- `func Max[T any](n int) Fact[[]T]` succeeds when a slice contains no more than `n` elements.
- `func All[T any](pred func(T) bool) Fact[[]T]` succeeds when `pred` returns true for every element.
- `func Any[T any](pred func(T) bool) Fact[[]T]` succeeds when `pred` returns true for at least one element.
- `func UniqueBy[T any, U comparable](selector func(T) U) Fact[[]T]` succeeds when `selector` returns a distinct key for every element.
- `func Unique[S comparable]() Fact[[]S]` succeeds when every element is distinct.
- `func Contains[T comparable](item T) Fact[[]T]` succeeds when the slice contains `item`.
- `func NotContains[T comparable](item T) Fact[[]T]` succeeds when the slice does not contain `item`.

### Map facts

- `func MapEmpty[K comparable, V any]() Fact[map[K]V]` succeeds when a map has no entries.
- `func MapNotEmpty[K comparable, V any]() Fact[map[K]V]` succeeds when a map has at least one entry.
- `func MapLength[K comparable, V any](n int) Fact[map[K]V]` succeeds when a map has exactly `n` entries.
- `func MapLengthBetween[K comparable, V any](min, max int) Fact[map[K]V]` succeeds when a map has between `min` and `max` entries, inclusive.
- `func MapContains[K comparable, V any](key K) Fact[map[K]V]` succeeds when a map contains `key`.
- `func MapNotContains[K comparable, V any](key K) Fact[map[K]V]` succeeds when a map does not contain `key`.

### Option facts

- `func Some[T any]() Fact[option.Option[T]]` succeeds when an `Option` contains a value.
- `func None[T any]() Fact[option.Option[T]]` succeeds when an `Option` has no value.
- `func Required[T any](facts ...Fact[T]) Fact[option.Option[T]]` requires a value and checks it against each fact.
- `func Optional[T any](facts ...Fact[T]) Fact[option.Option[T]]` checks a value when present and succeeds when it is `None`.

### Time and duration facts

- `func TimeBetween(s, e time.Time) Fact[time.Time]` succeeds when a time is between `s` and `e`, inclusive.
- `func TimeLessThan(t time.Time) Fact[time.Time]` succeeds when a time is before `t`.
- `func TimeGreaterThan(t time.Time) Fact[time.Time]` succeeds when a time is after `t`.
- `func Future(from time.Time) Fact[time.Time]` succeeds when a time is after `from`.
- `func Past(from time.Time) Fact[time.Time]` succeeds when a time is before `from`.
- `func DurationBetween(s, e time.Duration) Fact[time.Duration]` succeeds when a duration is between `s` and `e`, inclusive.
- `func DurationLessThan(t time.Duration) Fact[time.Duration]` succeeds when a duration is less than `t`.
- `func DurationGreaterThan(t time.Duration) Fact[time.Duration]` succeeds when a duration is greater than `t`.
- `func DurationPositive() Fact[time.Duration]` succeeds when a duration is greater than zero.

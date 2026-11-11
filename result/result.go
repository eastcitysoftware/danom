// Package result provides Result[T, E], a type-safe representation of an
// operation that either succeeds with a value (Ok) or fails with an expected
// error (Error), along with functions for matching, mapping, and binding over it.
package result

// Result represents either a value of type T or an error value of type E.
type Result[T any, E any] struct {
	ok    bool
	value T
	error E
}

// Ok creates a Result that contains t.
func Ok[T any, E any](t T) Result[T, E] {
	var zero E
	return Result[T, E]{true, t, zero}
}

// Error creates a Result that contains e.
func Error[T any, E any](e E) Result[T, E] {
	var zero T
	return Result[T, E]{false, zero, e}
}

// IsOk reports whether the Result contains a success value.
func (r Result[T, E]) IsOk() bool {
	return r.ok == true
}

// IsError reports whether the Result contains an error.
func (r Result[T, E]) IsError() bool {
	return r.ok == false
}

// Value returns the success value and true when the Result is Ok.
// Otherwise, it returns the zero value of T and false.
func (r Result[T, E]) Value() (T, bool) {
	if r.IsOk() {
		return r.value, true
	}

	var zero T
	return zero, false
}

// ErrorValue returns the error value and true when the Result is an Error.
// Otherwise, it returns the zero value of E and false.
func (r Result[T, E]) ErrorValue() (E, bool) {
	if r.IsError() {
		return r.error, true
	}

	var zero E
	return zero, false
}

// DefaultValue returns the success value, or d when the Result is an Error.
func (r Result[T, E]) DefaultValue(d T) T {
	if value, ok := r.Value(); ok {
		return value
	}
	return d
}

// DefaultWith returns the success value, or calls fn when the Result is an Error.
func (r Result[T, E]) DefaultWith(fn func() T) T {
	if value, ok := r.Value(); ok {
		return value
	}
	return fn()
}

// OrElse returns o when it is Ok, or ifError otherwise.
func (r Result[T, E]) OrElse(ifError Result[T, E]) Result[T, E] {
	if r.IsError() {
		return ifError
	}

	return r
}

// OrElseWith returns o when it is Ok, or calls ifError otherwise.
func (r Result[T, E]) OrElseWith(ifError func() Result[T, E]) Result[T, E] {
	if r.IsError() {
		return ifError()
	}

	return r
}

// Tap calls fn with the success value, then returns o unchanged.
func (r Result[T, E]) Tap(fn func(T)) Result[T, E] {
	if value, ok := r.Value(); ok {
		fn(value)
	}

	return r
}

// Filter returns o when it is Some and its value satisfies pred; otherwise it returns None.
func (r Result[T, E]) Filter(pred func(T) bool, ifError E) Result[T, E] {
	if errValue, err := r.ErrorValue(); err {
		return Error[T](errValue)
	}

	if value, ok := r.Value(); ok && pred(value) {
		return Ok[T, E](value)
	}

	return Error[T](ifError)
}

// Match calls the handler for the state of r and returns its result.
func Match[T any, E any, U any](r Result[T, E], ok func(T) U, error func(E) U) U {
	if value, isOk := r.Value(); isOk {
		return ok(value)
	}

	err, _ := r.ErrorValue()
	return error(err)
}

// Bind applies fn to the success value, or preserves the error value.
func Bind[T any, E any, U any](r Result[T, E], fn func(T) Result[U, E]) Result[U, E] {
	return Match(r, fn, Error[U, E])
}

// Map transforms the success value, or preserves the error value.
func Map[T any, E any, U any](r Result[T, E], fn func(T) U) Result[U, E] {
	return Bind(r, func(x T) Result[U, E] {
		return Ok[U, E](fn(x))
	})
}

// MapError transforms the error value, or preserves the success value.
func MapError[T any, E any, U any](r Result[T, E], fn func(E) U) Result[T, U] {
	return Match(r, Ok[T, U], func(e E) Result[T, U] {
		return Error[T, U](fn(e))
	})
}

// Apply applies a successful function to a successful value.
func Apply[T any, E any, U any](r Result[T, E], fn Result[func(T) U, E]) Result[U, E] {
	return Bind(fn, func(f func(T) U) Result[U, E] {
		return Map(r, f)
	})
}

// Flatten converts a nested Result into a single Result.
func Flatten[T any, E any](r Result[Result[T, E], E]) Result[T, E] {
	if value, ok := r.Value(); ok {
		return value
	}

	err, _ := r.ErrorValue()
	return Error[T, E](err)
}

// Do calls the handler for the state of r.
func Do[T any, E any](r Result[T, E], ok func(T), error func(E)) {
	if value, isOk := r.Value(); isOk {
		ok(value)
	} else {
		err, _ := r.ErrorValue()
		error(err)
	}
}

// Iter calls fn with the success value when r is Ok.
func Iter[T any, E any](o Result[T, E], fn func(T)) {
	Do(o, fn, func(_ E) {})
}

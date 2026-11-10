package result

const (
	error uint8 = iota
	ok    uint8 = 1
)

// Result represents either a value of type T or an error value of type E.
type Result[T any, E any] struct {
	state uint8
	ok    T
	error E
}

// IsOk reports whether the Result contains a success value.
func (o Result[T, E]) IsOk() bool {
	return o.state == ok
}

// IsError reports whether the Result contains an error.
func (o Result[T, E]) IsError() bool {
	return o.state == error
}

// Value returns the success value and true when the Result is Ok.
// Otherwise, it returns the zero value of T and false.
func (o Result[T, E]) Value() (T, bool) {
	if o.IsOk() {
		return o.ok, true
	}

	var zero T
	return zero, false
}

// ErrorValue returns the error value and true when the Result is an Error.
// Otherwise, it returns the zero value of E and false.
func (o Result[T, E]) ErrorValue() (E, bool) {
	if o.IsError() {
		return o.error, true
	}

	var zero E
	return zero, false
}

// DefaultValue returns the success value, or d when the Result is an Error.
func (o Result[T, E]) DefaultValue(d T) T {
	if value, ok := o.Value(); ok {
		return value
	}
	return d
}

// DefaultWith returns the success value, or calls fn when the Result is an Error.
func (o Result[T, E]) DefaultWith(fn func() T) T {
	if value, ok := o.Value(); ok {
		return value
	}
	return fn()
}

// Ok creates a Result that contains t.
func Ok[T any, E any](t T) Result[T, E] {
	var zero E
	return Result[T, E]{ok, t, zero}
}

// Error creates a Result that contains e.
func Error[T any, E any](e E) Result[T, E] {
	var zero T
	return Result[T, E]{error, zero, e}
}

// Match calls the handler for the state of r and returns its result.
func Match[T any, E any, U any](r Result[T, E], ok func(T) U, error func(E) U) U {
	if value, isOk := r.Value(); isOk {
		return ok(value)
	}

	err, _ := r.ErrorValue()
	return error(err)
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
func Apply[T any, E any, U any](o Result[T, E], fn Result[func(T) U, E]) Result[U, E] {
	return Bind(fn, func(f func(T) U) Result[U, E] {
		return Map(o, f)
	})
}

// OrElse returns o when it is Ok, or ifError otherwise.
func OrElse[T any, E any](o Result[T, E], ifError Result[T, E]) Result[T, E] {
	return Match(o, Ok[T, E], func(_ E) Result[T, E] { return ifError })
}

// OrElseWith returns o when it is Ok, or calls ifError otherwise.
func OrElseWith[T any, E any](o Result[T, E], ifError func() Result[T, E]) Result[T, E] {
	return Match(o, Ok[T, E], func(_ E) Result[T, E] { return ifError() })
}

// Tap calls fn with the success value, then returns o unchanged.
func Tap[T any, E any](o Result[T, E], fn func(T)) Result[T, E] {
	return Map(o, func(t T) T {
		fn(t)
		return t
	})
}

// Iter calls fn with the success value when r is Ok.
func Iter[T any, E any](o Result[T, E], fn func(T)) {
	Do(o, fn, func(_ E) {})
}

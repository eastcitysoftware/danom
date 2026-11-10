// Package option provides Option[T], a type-safe representation of a value
// that may be present (Some) or absent (None), along with functions for
// matching, mapping, and binding over it.
package option

const (
	none uint8 = iota
	some uint8 = 1
)

// Option represents a value that may or may not be present.
type Option[T any] struct {
	state uint8
	value T
}

// IsSome reports whether the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.state == some
}

// IsNone reports whether the Option has no value.
func (o Option[T]) IsNone() bool {
	return o.state == none
}

// Value returns the value and true when the Option is Some.
// Otherwise, it returns the zero value of T and false.
func (o Option[T]) Value() (T, bool) {
	if o.IsSome() {
		return o.value, true
	}

	var zero T
	return zero, false
}

// DefaultValue returns the Option's value, or d when it is None.
func (o Option[T]) DefaultValue(d T) T {
	if value, ok := o.Value(); ok {
		return value
	}
	return d
}

// DefaultWith returns the Option's value, or calls fn when it is None.
func (o Option[T]) DefaultWith(fn func() T) T {
	if value, ok := o.Value(); ok {
		return value
	}
	return fn()
}

// IsSome reports whether o contains a value.
func IsSome[T any](o Option[T]) bool {
	return o.IsSome()
}

// IsNone reports whether o has no value.
func IsNone[T any](o Option[T]) bool {
	return o.IsNone()
}

// Some creates an Option that contains val.
func Some[T any](val T) Option[T] {
	return Option[T]{state: some, value: val}
}

// None creates an Option with no value.
func None[T any]() Option[T] {
	return Option[T]{state: none}
}

// Match calls some with the value when o is Some, or calls none otherwise.
func Match[T any, U any](o Option[T], some func(T) U, none func() U) U {
	if value, ok := o.Value(); ok {
		return some(value)
	}

	return none()
}

// Do calls the handler for the state of o.
func Do[T any](o Option[T], some func(T), none func()) {
	if o.IsSome() {
		some(o.value)
	} else {
		none()
	}
}

// Bind applies fn to the value in o, or returns None when o is None.
func Bind[T any, U any](o Option[T], fn func(T) Option[U]) Option[U] {
	return Match(o, fn, None[U])
}

// Map transforms the value in o, or returns None when o is None.
func Map[T any, U any](o Option[T], fn func(T) U) Option[U] {
	return Bind(o, func(x T) Option[U] {
		return Some(fn(x))
	})
}

// Apply applies an optional function to an optional value.
func Apply[T any, U any](o Option[T], fn Option[func(T) U]) Option[U] {
	return Bind(fn, func(f func(T) U) Option[U] {
		return Map(o, f)
	})
}

// OrElse returns o when it is Some, or ifNone otherwise.
func OrElse[T any](o Option[T], ifNone Option[T]) Option[T] {
	return Match(o, Some, func() Option[T] { return ifNone })
}

// OrElseWith returns o when it is Some, or calls ifNone otherwise.
func OrElseWith[T any](o Option[T], ifNone func() Option[T]) Option[T] {
	return Match(o, Some, ifNone)
}

// Tap calls fn with the value when o is Some, then returns o unchanged.
func Tap[T any](o Option[T], fn func(T)) Option[T] {
	return Map(o, func(t T) T {
		fn(t)
		return t
	})
}

// Iter calls fn with the value when o is Some.
func Iter[T any](o Option[T], fn func(T)) {
	Do(o, fn, func() {})
}

// Filter returns o when it is Some and its value satisfies pred; otherwise it returns None.
func Filter[T any](o Option[T], pred func(T) bool) Option[T] {
	if value, ok := o.Value(); ok && pred(value) {
		return Some(value)
	}
	return None[T]()
}

// Flatten converts a nested Option into a single Option.
func Flatten[T any](o Option[Option[T]]) Option[T] {
	return Bind(o, func(inner Option[T]) Option[T] {
		return inner
	})
}

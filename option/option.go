package option

const (
	none uint8 = iota
	some uint8 = 1
)

type Option[T any] struct {
	state uint8
	value T
}

func (o Option[T]) IsSome() bool {
	return o.state == some
}

func (o Option[T]) IsNone() bool {
	return o.state == none
}

func (o Option[T]) TryGet() (T, bool) {
	if o.IsSome() {
		return o.value, true
	}

	var zero T
	return zero, false
}

func (o Option[T]) DefaultValue(d T) T {
	if value, ok := o.TryGet(); ok {
		return value
	}
	return d
}

func (o Option[T]) DefaultWith(fn func() T) T {
	if value, ok := o.TryGet(); ok {
		return value
	}
	return fn()
}

func Some[T any](val T) Option[T] {
	return Option[T]{state: some, value: val}
}

func None[T any]() Option[T] {
	return Option[T]{state: none}
}

func Match[T any, U any](o Option[T], some func(T) U, none func() U) U {
	if value, ok := o.TryGet(); ok {
		return some(value)
	}

	return none()
}

func Do[T any](o Option[T], some func(T), none func()) {
	if o.IsSome() {
		some(o.value)
	} else {
		none()
	}
}

func Bind[T any, U any](o Option[T], fn func(T) Option[U]) Option[U] {
	return Match(o, fn, None[U])
}

func Map[T any, U any](o Option[T], fn func(T) U) Option[U] {
	return Bind(o, func(x T) Option[U] {
		return Some(fn(x))
	})
}

func Apply[T any, U any](o Option[T], fn Option[func(T) U]) Option[U] {
	return Bind(fn, func(f func(T) U) Option[U] {
		return Map(o, f)
	})
}

func OrElse[T any](o Option[T], ifNone Option[T]) Option[T] {
	return Match(o, Some, func() Option[T] { return ifNone })
}

func OrElseWith[T any](o Option[T], ifNone func() Option[T]) Option[T] {
	return Match(o, Some, ifNone)
}

func Tap[T any](o Option[T], fn func(T)) Option[T] {
	return Map(o, func(t T) T {
		fn(t)
		return t
	})
}

func Iter[T any](o Option[T], fn func(T)) {
	Do(o, fn, func() {})
}

func Filter[T any](o Option[T], pred func(T) bool) Option[T] {
	if value, ok := o.TryGet(); ok && pred(value) {
		return Some(value)
	}
	return None[T]()
}

func Flatten[T any](o Option[Option[T]]) Option[T] {
	return Bind(o, func(inner Option[T]) Option[T] {
		return inner
	})
}

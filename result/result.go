package result

const (
	error uint8 = iota
	ok    uint8 = 1
)

type Result[T any, E any] struct {
	state uint8
	ok    T
	error E
}

func (o Result[T, E]) IsOk() bool {
	return o.state == ok
}

func (o Result[T, E]) IsError() bool {
	return o.state == error
}

func (o Result[T, E]) TryGet() (T, E, bool) {
	var zeroT T
	var zeroE E

	if o.IsOk() {
		return o.ok, zeroE, true
	}

	return zeroT, o.error, false
}

func (o Result[T, E]) DefaultValue(d T) T {
	if value, _, ok := o.TryGet(); ok {
		return value
	}
	return d
}

func (o Result[T, E]) DefaultWith(fn func() T) T {
	if value, _, ok := o.TryGet(); ok {
		return value
	}
	return fn()
}

func Ok[T any, E any](t T) Result[T, E] {
	var zero E
	return Result[T, E]{ok, t, zero}
}

func Error[T any, E any](e E) Result[T, E] {
	var zero T
	return Result[T, E]{error, zero, e}
}

func Match[T any, E any, U any](r Result[T, E], ok func(T) U, error func(E) U) U {
	v, e, isOk := r.TryGet()
	if isOk {
		return ok(v)
	}
	return error(e)
}

func Do[T any, E any](r Result[T, E], ok func(T), error func(E)) {
	v, e, isOk := r.TryGet()
	if isOk {
		ok(v)
	} else {
		error(e)
	}
}

func Bind[T any, E any, U any](r Result[T, E], fn func(T) Result[U, E]) Result[U, E] {
	return Match(r, fn, Error[U, E])
}

func Map[T any, E any, U any](r Result[T, E], fn func(T) U) Result[U, E] {
	return Bind(r, func(x T) Result[U, E] {
		return Ok[U, E](fn(x))
	})
}

func MapError[T any, E any, U any](r Result[T, E], fn func(E) U) Result[T, U] {
	return Match(r, Ok[T, U], func(e E) Result[T, U] {
		return Error[T, U](fn(e))
	})
}

func Apply[T any, E any, U any](o Result[T, E], fn Result[func(T) U, E]) Result[U, E] {
	return Bind(fn, func(f func(T) U) Result[U, E] {
		return Map(o, f)
	})
}

func OrElse[T any, E any](o Result[T, E], ifError Result[T, E]) Result[T, E] {
	return Match(o, Ok[T, E], func(_ E) Result[T, E] { return ifError })
}

func OrElseWith[T any, E any](o Result[T, E], ifError func() Result[T, E]) Result[T, E] {
	return Match(o, Ok[T, E], func(_ E) Result[T, E] { return ifError() })
}

func Tap[T any, E any](o Result[T, E], fn func(T)) Result[T, E] {
	return Map(o, func(t T) T {
		fn(t)
		return t
	})
}

func Iter[T any, E any](o Result[T, E], fn func(T)) {
	Do(o, fn, func(_ E) {})
}

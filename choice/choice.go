package choice

const (
	t1 uint8 = iota
	t2 uint8 = 1
)

type Choice[T1 any, T2 any] struct {
	state uint8
	t1    T1
	t2    T2
}

func (c Choice[T1, T2]) IsT1() bool {
	return c.state == t1
}

func (c Choice[T1, T2]) IsT2() bool {
	return c.state == t2
}

func (o Choice[T1, T2]) TryGetT1() (T1, bool) {
	if o.IsT1() {
		return o.t1, true
	}

	var zero T1
	return zero, false
}

func (o Choice[T1, T2]) TryGetT2() (T2, bool) {
	if o.IsT2() {
		return o.t2, true
	}

	var zero T2
	return zero, false
}

func (o Choice[T1, T2]) Swap() Choice[T2, T1] {
	if o.IsT1() {
		return OfT2[T2, T1](o.t1)
	}

	return OfT1[T2, T1](o.t2)
}

func OfT1[T1 any, T2 any](t T1) Choice[T1, T2] {
	var zero T2
	return Choice[T1, T2]{t1, t, zero}
}

func OfT2[T1 any, T2 any](t T2) Choice[T1, T2] {
	var zero T1
	return Choice[T1, T2]{t2, zero, t}
}

func Match[T1 any, T2, U any](o Choice[T1, T2], t1Fn func(T1) U, t2Fn func(T2) U) U {
	if value, ok := o.TryGetT1(); ok {
		return t1Fn(value)
	}

	if value, ok := o.TryGetT2(); ok {
		return t2Fn(value)
	}

	var zero U
	return zero
}

func Do[T1 any, T2 any](o Choice[T1, T2], t1Fn func(T1), t2Fn func(T2)) {
	if value, ok := o.TryGetT1(); ok {
		t1Fn(value)
	}

	if value, ok := o.TryGetT2(); ok {
		t2Fn(value)
	}
}

func BindT1[T1 any, T2 any, U any](o Choice[T1, T2], fn func(T1) Choice[U, T2]) Choice[U, T2] {
	return Match(o, fn, OfT2)
}

func MapT1[T1 any, T2 any, U any](o Choice[T1, T2], fn func(T1) U) Choice[U, T2] {
	return BindT1(o, func(x T1) Choice[U, T2] {
		return OfT1[U, T2](fn(x))
	})
}

func BindT2[T1 any, T2 any, U any](o Choice[T1, T2], fn func(T2) Choice[T1, U]) Choice[T1, U] {
	return Match(o, OfT1, fn)
}

func MapT2[T1 any, T2 any, U any](o Choice[T1, T2], fn func(T2) U) Choice[T1, U] {
	return BindT2(o, func(x T2) Choice[T1, U] {
		return OfT2[T1, U](fn(x))
	})
}

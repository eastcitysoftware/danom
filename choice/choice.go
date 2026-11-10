package choice

const (
	t1 uint8 = iota
	t2 uint8 = 1
)

// Choice holds a value of either T1 or T2.
type Choice[T1 any, T2 any] struct {
	state uint8
	t1    T1
	t2    T2
}

// IsT1 reports whether the Choice holds a value of type T1.
func (c Choice[T1, T2]) IsT1() bool {
	return c.state == t1
}

// IsT2 reports whether the Choice holds a value of type T2.
func (c Choice[T1, T2]) IsT2() bool {
	return c.state == t2
}

// T1 returns the T1 value and true when the Choice holds T1.
// Otherwise, it returns the zero value of T1 and false.
func (o Choice[T1, T2]) T1() (T1, bool) {
	if o.IsT1() {
		return o.t1, true
	}

	var zero T1
	return zero, false
}

// T2 returns the T2 value and true when the Choice holds T2.
// Otherwise, it returns the zero value of T2 and false.
func (o Choice[T1, T2]) T2() (T2, bool) {
	if o.IsT2() {
		return o.t2, true
	}

	var zero T2
	return zero, false
}

// Swap returns a Choice with T1 and T2 reversed.
func (o Choice[T1, T2]) Swap() Choice[T2, T1] {
	if o.IsT1() {
		return OfT2[T2, T1](o.t1)
	}

	return OfT1[T2, T1](o.t2)
}

// OfT1 creates a Choice that holds a value of type T1.
func OfT1[T1 any, T2 any](t T1) Choice[T1, T2] {
	var zero T2
	return Choice[T1, T2]{t1, t, zero}
}

// OfT2 creates a Choice that holds a value of type T2.
func OfT2[T1 any, T2 any](t T2) Choice[T1, T2] {
	var zero T1
	return Choice[T1, T2]{t2, zero, t}
}

// Match calls the function for the type held by the Choice and returns its result.
func Match[T1 any, T2, U any](o Choice[T1, T2], t1Fn func(T1) U, t2Fn func(T2) U) U {
	if value, ok := o.T1(); ok {
		return t1Fn(value)
	}

	if value, ok := o.T2(); ok {
		return t2Fn(value)
	}

	var zero U
	return zero
}

// Do calls the function for the type held by the Choice.
func Do[T1 any, T2 any](o Choice[T1, T2], t1Fn func(T1), t2Fn func(T2)) {
	if value, ok := o.T1(); ok {
		t1Fn(value)
	}

	if value, ok := o.T2(); ok {
		t2Fn(value)
	}
}

// BindT1 applies fn to a T1 value and leaves a T2 value unchanged.
func BindT1[T1 any, T2 any, U any](o Choice[T1, T2], fn func(T1) Choice[U, T2]) Choice[U, T2] {
	return Match(o, fn, OfT2)
}

// MapT1 transforms a T1 value and leaves a T2 value unchanged.
func MapT1[T1 any, T2 any, U any](o Choice[T1, T2], fn func(T1) U) Choice[U, T2] {
	return BindT1(o, func(x T1) Choice[U, T2] {
		return OfT1[U, T2](fn(x))
	})
}

// BindT2 applies fn to a T2 value and leaves a T1 value unchanged.
func BindT2[T1 any, T2 any, U any](o Choice[T1, T2], fn func(T2) Choice[T1, U]) Choice[T1, U] {
	return Match(o, OfT1, fn)
}

// MapT2 transforms a T2 value and leaves a T1 value unchanged.
func MapT2[T1 any, T2 any, U any](o Choice[T1, T2], fn func(T2) U) Choice[T1, U] {
	return BindT2(o, func(x T2) Choice[T1, U] {
		return OfT2[T1, U](fn(x))
	})
}

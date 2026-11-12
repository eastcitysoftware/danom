package check

// Number includes the built-in integer, unsigned integer, and floating-point types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Zero succeeds when a numeric value is zero.
func Zero[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x == 0 },
		NewMessage("'%s' must be equal to 0"))
}

// Positive succeeds when a numeric value is greater than zero.
func Positive[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x > 0 },
		NewMessage("'%s' must be positive"))
}

// Negative succeeds when a numeric value is less than zero.
func Negative[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x < 0 },
		NewMessage("'%s' must be negative"))
}

// NonZero succeeds when a numeric value is not zero.
func NonZero[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x != 0 },
		NewMessage("'%s' must not be zero"))
}

// NonPositive succeeds when a numeric value is less than or equal to zero.
func NonPositive[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x <= 0 },
		NewMessage("'%s' must be non-positive"))
}

// NonNegative succeeds when a numeric value is greater than or equal to zero.
func NonNegative[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x >= 0 },
		NewMessage("'%s' must be non-negative"))
}

package check

import "cmp"

// Between succeeds when an ordered value is between min and max, inclusive.
func Between[N cmp.Ordered](min, max N) Fact[N] {
	return NewFact(
		func(x N) bool { return x >= min && x <= max },
		NewMessage("'%s' must be between %v and %v", min, max))
}

// NotBetween succeeds when an ordered value is outside the inclusive range [min, max].
func NotBetween[N cmp.Ordered](min, max N) Fact[N] {
	return NewFact(
		func(x N) bool { return x < min || x > max },
		NewMessage("'%s' must not be between %v and %v", min, max))
}

// GreaterThan succeeds when an ordered value is greater than min.
func GreaterThan[N cmp.Ordered](min N) Fact[N] {
	return NewFact(
		func(x N) bool { return x > min },
		NewMessage("'%s' must be greater than %v", min))
}

// GreaterThanOrEqual succeeds when an ordered value is at least minimum.
func GreaterThanOrEqual[N cmp.Ordered](minimum N) Fact[N] {
	return NewFact(
		func(x N) bool { return x >= minimum },
		NewMessage("'%s' must be greater than or equal to %v", minimum))
}

// LessThan succeeds when an ordered value is less than min.
func LessThan[N cmp.Ordered](min N) Fact[N] {
	return NewFact(
		func(x N) bool { return x < min },
		NewMessage("'%s' must be less than %v", min))
}

// LessThanOrEqual succeeds when an ordered value is at most minimum.
func LessThanOrEqual[N cmp.Ordered](minimum N) Fact[N] {
	return NewFact(
		func(x N) bool { return x <= minimum },
		NewMessage("'%s' must be less than or equal to %v", minimum))
}

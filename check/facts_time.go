package check

import "time"

// TimeBetween succeeds when a time is between s and e, inclusive.
func TimeBetween(s, e time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return (x.Before(e) || x.Equal(e)) && (x.After(s) || x.Equal(s)) },
		NewMessage("'%s' must be on or after %s and on or before %s",
			s.Format(time.RFC3339Nano), e.Format(time.RFC3339Nano)))
}

// TimeLessThan succeeds when a time is before t.
func TimeLessThan(t time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.Before(t) },
		NewMessage("'%s' must be before %v", t))
}

// TimeGreaterThan succeeds when a time is after t.
func TimeGreaterThan(t time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.After(t) },
		NewMessage("'%s' must be after %v", t))
}

// Future succeeds when a time is after from.
func Future(from time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.After(from) },
		NewMessage("'%s' must be in the future"))
}

// Past succeeds when a time is before from.
func Past(from time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.Before(from) },
		NewMessage("'%s' must be in the past"))
}

// DurationBetween succeeds when a duration is between s and e, inclusive.
func DurationBetween(s, e time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x >= s && x <= e },
		NewMessage("'%s' must be between %v and %v", s, e))
}

// DurationLessThan succeeds when a duration is less than t.
func DurationLessThan(t time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x < t },
		NewMessage("'%s' must be less than %v", t))
}

// DurationGreaterThan succeeds when a duration is greater than t.
func DurationGreaterThan(t time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x > t },
		NewMessage("'%s' must be greater than %v", t))
}

// DurationPositive succeeds when a duration is greater than zero.
func DurationPositive() Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x > 0 },
		NewMessage("'%s' must be positive"))
}

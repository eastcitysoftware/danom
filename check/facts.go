package check

// When runs facts when cond succeeds and returns their messages.
func When[T any](cond Fact[T], facts ...Fact[T]) Fact[T] {
	return func(x T) []Message {
		var m []Message
		if len(cond(x)) == 0 {
			for _, f := range facts {
				fm := f(x)
				if len(fm) > 0 {
					m = append(m, fm...)
				}
			}
		}

		if len(m) > 0 {
			return m
		}

		return nil
	}
}

// NoneOf succeeds when every fact fails.
func NoneOf[T any](facts ...Fact[T]) Fact[T] {
	return func(x T) []Message {
		for _, f := range facts {
			if len(f(x)) == 0 {
				return []Message{NewMessage("'%s' must not satisfy any of the conditions")}
			}
		}
		return nil
	}
}

// AnyOf succeeds when at least one fact succeeds.
func AnyOf[T any](facts ...Fact[T]) Fact[T] {
	return func(x T) []Message {
		for _, f := range facts {
			if len(f(x)) == 0 {
				return nil
			}
		}
		return []Message{NewMessage("'%s' must satisfy at least one condition")}
	}
}

// Equals succeeds when a value equals e.
func Equals[C comparable](e C) Fact[C] {
	return NewFact(
		func(x C) bool { return x == e },
		NewMessage("'%s' must be equal to %v", e))
}

// NotEquals succeeds when a value does not equal e.
func NotEquals[C comparable](e C) Fact[C] {
	return NewFact(
		func(x C) bool { return x != e },
		NewMessage("'%s' must not be equal to %v", e))
}

// True succeeds when a value is true.
func True() Fact[bool] {
	return NewFact(
		func(x bool) bool { return x == true },
		NewMessage("'%s' must be true"))
}

// False succeeds when a value is false.
func False() Fact[bool] {
	return NewFact(
		func(x bool) bool { return x != true },
		NewMessage("'%s' must be false"))
}

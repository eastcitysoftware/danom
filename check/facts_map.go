package check

// MapEmpty succeeds when a map has no entries.
func MapEmpty[K comparable, V any]() Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) == 0 },
		NewMessage("'%s' must be empty"))
}

// MapNotEmpty succeeds when a map has at least one entry.
func MapNotEmpty[K comparable, V any]() Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) > 0 },
		NewMessage("'%s' must not be empty"))
}

// MapLength succeeds when a map has exactly n entries.
func MapLength[K comparable, V any](n int) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) == n },
		NewMessage("'%s' must have %v keys", n))
}

// MapLengthBetween succeeds when a map has between min and max entries, inclusive.
func MapLengthBetween[K comparable, V any](min, max int) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) >= min && len(m) <= max },
		NewMessage("'%s' must have between %v and %v keys", min, max))
}

// MapContains succeeds when a map contains key.
func MapContains[K comparable, V any](key K) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { _, ok := m[key]; return ok },
		NewMessage("'%s' is missing key: %v", key))
}

// MapNotContains succeeds when a map does not contain key.
func MapNotContains[K comparable, V any](key K) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { _, ok := m[key]; return !ok },
		NewMessage("'%s' must not have key: %v", key))
}

// MapAll succeeds when all facts succeed.
func MapAll[K comparable, V any](facts ...Fact[map[K]V]) Fact[map[K]V] {
	return func(m map[K]V) []Message {
		var msgs []Message
		for _, f := range facts {
			msgs = append(msgs, f(m)...)
		}
		return msgs
	}
}

// MapAny succeeds when any fact succeeds.
func MapAny[K comparable, V any](facts ...Fact[map[K]V]) Fact[map[K]V] {
	return func(m map[K]V) []Message {
		for _, f := range facts {
			if len(f(m)) == 0 {
				return nil
			}
		}
		return []Message{NewMessage("'%s' must satisfy at least one condition")}
	}
}

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

// MapAll succeeds when pred returns true for every map value.
// It succeeds for an empty map.
func MapAll[K comparable, V any](pred func(V) bool) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool {
			for _, value := range m {
				if !pred(value) {
					return false
				}
			}
			return true
		},
		NewMessage("'%s' must contain only valid values"))
}

// MapAny succeeds when pred returns true for at least one map value.
// It fails for an empty map.
func MapAny[K comparable, V any](pred func(V) bool) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool {
			for _, value := range m {
				if pred(value) {
					return true
				}
			}
			return false
		},
		NewMessage("'%s' must contain at least one valid value"))
}

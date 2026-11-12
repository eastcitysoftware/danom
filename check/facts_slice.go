package check

import "slices"

// Empty succeeds when a slice has no elements.
func Empty[T any]() Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) == 0 },
		NewMessage("'%s' must be empty"))
}

// NotEmpty succeeds when a slice has at least one element.
func NotEmpty[T any]() Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) > 0 },
		NewMessage("'%s' must not be empty"))
}

// Length succeeds when a slice contains exactly n elements.
func Length[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) == n },
		NewMessage("'%s' must contain %v items", n))
}

// LengthBetween succeeds when a slice contains between min and max elements, inclusive.
func LengthBetween[T any](min, max int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) >= min && len(items) <= max },
		NewMessage("'%s' must contain between %v and %v items", min, max))
}

// Min succeeds when a slice contains at least n elements.
func Min[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) >= n },
		NewMessage("'%s' must contain at least %v items", n))
}

// Max succeeds when a slice contains no more than n elements.
func Max[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) <= n },
		NewMessage("'%s' must contain no more than %v items", n))
}

// All succeeds when pred returns true for every slice element.
func All[T any](pred func(T) bool) Fact[[]T] {
	return func(items []T) []Message {
		for i, item := range items {
			if !pred(item) {
				return []Message{
					NewMessage("'%s' contains an invalid item at index %v", i+1),
				}
			}
		}
		return nil
	}
}

// Any succeeds when pred returns true for at least one slice element.
func Any[T any](pred func(T) bool) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			return slices.ContainsFunc(items, pred)
		},
		NewMessage("'%s' must contain one valid item"))
}

// UniqueBy succeeds when selector returns a distinct key for every slice element.
func UniqueBy[T any, U comparable](selector func(T) U) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			seen := make(map[U]struct{}, len(items))

			for _, item := range items {
				key := selector(item)
				if _, exists := seen[key]; exists {
					return false
				}
				seen[key] = struct{}{}
			}

			return true
		},
		NewMessage("'%s' must contain unique items"))
}

// Unique succeeds when every element in a slice is distinct.
func Unique[S comparable]() Fact[[]S] {
	return UniqueBy(func(x S) S { return x })
}

// Contains succeeds when a slice contains item.
func Contains[T comparable](item T) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			for _, i := range items {
				if i == item {
					return true
				}
			}
			return false
		},
		NewMessage("'%s' is missing the item: %v", item))
}

// NotContains succeeds when a slice does not contain item.
func NotContains[T comparable](item T) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			for _, i := range items {
				if i == item {
					return false
				}
			}
			return true
		},
		NewMessage("'%s' must not contain the item: %v", item))
}

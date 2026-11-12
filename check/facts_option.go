package check

import "github.com/eastcitysoftware/danom/option"

// Some succeeds when an Option contains a value.
func Some[T any]() Fact[option.Option[T]] {
	return NewFact[option.Option[T]](
		func(o option.Option[T]) bool { return o.IsSome() },
		NewMessage("'%s' should have a value"))
}

// None succeeds when an Option has no value.
func None[T any]() Fact[option.Option[T]] {
	return NewFact[option.Option[T]](
		func(o option.Option[T]) bool { return o.IsNone() },
		NewMessage("'%s' should have no value"))
}

// Required succeeds when an Option has a value and that value passes every fact.
func Required[T any](facts ...Fact[T]) Fact[option.Option[T]] {
	return func(x option.Option[T]) []Message {
		return option.Match(x,
			func(some T) []Message {
				var m []Message
				for _, fact := range facts {
					for _, message := range fact(some) {
						m = append(m, message)
					}
				}
				return m
			},
			func() []Message { return []Message{NewMessage("'%s' is required")} })
	}
}

// Optional validates an Option's value when present and succeeds when it is None.
func Optional[T any](facts ...Fact[T]) Fact[option.Option[T]] {
	return func(x option.Option[T]) []Message {
		return option.Match(x,
			func(some T) []Message {
				var m []Message
				for _, fact := range facts {
					for _, message := range fact(some) {
						m = append(m, message)
					}
				}
				return m
			},
			func() []Message { return []Message{} })
	}
}

package option

import (
	"testing"
)

func TestSomeNone(t *testing.T) {
	expectNone(t, None[uint8]())
	expectSome(t, Some[uint8](1))

	expectNone(t, None[uint16]())
	expectSome(t, Some[uint16](1))

	expectNone(t, None[uint32]())
	expectSome(t, Some[uint32](1))

	expectNone(t, None[uint64]())
	expectSome(t, Some[uint64](1))

	expectNone(t, None[int8]())
	expectSome(t, Some[int8](1))

	expectNone(t, None[int16]())
	expectSome(t, Some[int16](1))

	expectNone(t, None[int32]())
	expectSome(t, Some[int32](1))

	expectNone(t, None[int64]())
	expectSome(t, Some[int64](1))

	expectNone(t, None[float32]())
	expectSome(t, Some[float32](1))

	expectNone(t, None[float64]())
	expectSome(t, Some[float64](1))

	expectNone(t, None[complex64]())
	expectSome(t, Some[complex64](1))

	expectNone(t, None[complex128]())
	expectSome(t, Some[complex128](1))

	expectNone(t, None[byte]())
	expectSome(t, Some[byte](1))

	expectNone(t, None[rune]())
	expectSome(t, Some('d'))

	expectNone(t, None[string]())
	expectSome(t, Some("danom"))

	expectNone(t, None[mockType]())
	expectSome(t, Some(mockType{name: "danom", age: Some(86)}))

	fieldNone := mockType{name: "danom", age: None[int]()}
	expectNone(t, fieldNone.age)

	fieldSome := mockType{name: "danom", age: Some(86)}
	expectSome(t, fieldSome.age)
}

func TestSomeNil(t *testing.T) {
	var ptr *int

	o := Some(ptr)

	value, ok := o.TryGet()

	if !ok {
		t.Error("Some(nil) should still be Some")
	}

	if value != nil {
		t.Error("Expected nil pointer value")
	}
}

func TestTryGet(t *testing.T) {
	if _, ok := Some(1).TryGet(); !ok {
		t.Error("Expected value, but found none")
	}

	if _, ok := None[int]().TryGet(); ok {
		t.Error("Expected none, but found value")
	}
}

func TestTryGetZeroValue(t *testing.T) {
	var o Option[int]

	value, ok := o.TryGet()

	if ok {
		t.Error("Zero-value Option should be None")
	}

	if value != 0 {
		t.Error("TryGet should return zero value for None")
	}
}

func TestDefaultValue(t *testing.T) {
	value := Some(1).DefaultValue(0)
	if value != 1 {
		t.Error("Default value should not have been used")
	}

	value = None[int]().DefaultValue(1)

	if value != 1 {
		t.Error("Default value should have been used")
	}
}

func TestDefaultWith(t *testing.T) {
	value := Some(1).DefaultWith(func() int { return 0 })
	if value != 1 {
		t.Error("Default value should not have been used")
	}

	value = None[int]().DefaultWith(func() int { return 1 })

	if value != 1 {
		t.Error("Default value should have been used")
	}
}

func TestDefaultWithLazy(t *testing.T) {
	called := false

	value := Some(1).DefaultWith(func() int {
		called = true
		return 2
	})

	if called {
		t.Error("DefaultWith should not invoke fallback when Some")
	}

	if value != 1 {
		t.Error("Expected original value")
	}
}

func TestMatch(t *testing.T) {
	if Match(None[string](),
		func(x string) bool { return true },
		func() bool { return false },
	) {
		t.Error("Should not have invoked Some branch in match")
	}

	if Match(Some("danom"),
		func(x string) bool { return false },
		func() bool { return true },
	) {
		t.Error("Should not have invoked Some branch in match")
	}
}

func TestDo(t *testing.T) {
	state := 0
	Do(None[string](),
		func(x string) { state = 1 },
		func() {})

	if state == 1 {
		t.Error("Should not have invoked Some branch in match")
	}

	Do(Some[string]("danom"),
		func(x string) {},
		func() { state = 2 })

	if state == 2 {
		t.Error("Should not have invoked Some branch in match")
	}
}

func TestBind(t *testing.T) {
	result := Bind(Some(1), func(n int) Option[int] { return Some(n + 1) })
	expectSome(t, result)

	value := Match(result, func(n int) int { return n }, func() int { return 0 })
	if value != 2 {
		t.Error("Bind function was not executed as expected")
	}
}

func TestBindNone(t *testing.T) {
	called := false

	result := Bind(None[int](), func(n int) Option[int] {
		called = true
		return Some(n + 1)
	})

	expectNone(t, result)

	if called {
		t.Error("Bind should not invoke callback for None")
	}
}

func TestMap(t *testing.T) {
	result := Map(Some(1), func(n int) int { return n + 1 })
	expectSome(t, result)

	value := Match(result, func(n int) int { return n }, func() int { return 0 })
	if value != 2 {
		t.Error("Bind function was not executed as expected")
	}
}

func TestMapNone(t *testing.T) {
	called := false

	result := Map(None[int](), func(n int) int {
		called = true
		return n + 1
	})

	expectNone(t, result)

	if called {
		t.Error("Map should not invoke callback for None")
	}
}
func TestApply(t *testing.T) {
	result := Apply(Some(1), Some(func(n int) int { return n + 1 }))
	expectSome(t, result)

	value := Match(result, func(n int) int { return n }, func() int { return 0 })
	if value != 2 {
		t.Error("Bind function was not executed as expected")
	}

	result = Apply(Some(1), None[func(n int) int]())
	expectNone(t, result)

	value = Match(result, func(n int) int { return n }, func() int { return 0 })
	if value != 0 {
		t.Error("Apply function should not return Some")
	}
}

func TestOrElse(t *testing.T) {
	result := OrElse(Some(1), Some(2))
	expectSome(t, result)

	value := result.DefaultValue(0)
	if value != 1 {
		t.Error("OrElse should preserve existing Some value")
	}

	result = OrElse(None[int](), Some(2))
	expectSome(t, result)

	value = result.DefaultValue(0)
	if value != 2 {
		t.Error("OrElse should return fallback when None")
	}

	expectNone(t, OrElse(None[int](), None[int]()))
}

func TestOrElseWith(t *testing.T) {
	called := false

	result := OrElseWith(Some(1), func() Option[int] {
		called = true
		return Some(2)
	})

	expectSome(t, result)

	if called {
		t.Error("OrElseWith should not invoke fallback when Some")
	}

	value := result.DefaultValue(0)
	if value != 1 {
		t.Error("OrElseWith should preserve existing Some value")
	}

	result = OrElseWith(None[int](), func() Option[int] {
		called = true
		return Some(2)
	})

	expectSome(t, result)

	if !called {
		t.Error("OrElseWith should invoke fallback when None")
	}

	value = result.DefaultValue(0)
	if value != 2 {
		t.Error("OrElseWith should return fallback value")
	}
}

func TestForEach(t *testing.T) {
	state := 0

	ForEach(None[int](), func(n int) {
		state = n
	})

	if state != 0 {
		t.Error("ForEach should not invoke callback for None")
	}

	ForEach(Some(42), func(n int) {
		state = n
	})

	if state != 42 {
		t.Error("ForEach should invoke callback with Some value")
	}
}

func TestFilter(t *testing.T) {
	result := Filter(Some(10), func(n int) bool {
		return n > 5
	})

	expectSome(t, result)

	value := result.DefaultValue(0)
	if value != 10 {
		t.Error("Filter should preserve value when predicate matches")
	}

	result = Filter(Some(3), func(n int) bool {
		return n > 5
	})

	expectNone(t, result)

	called := false

	result = Filter(None[int](), func(n int) bool {
		called = true
		return true
	})

	expectNone(t, result)

	if called {
		t.Error("Filter should not invoke predicate for None")
	}
}

func TestFilterPreservesZeroValue(t *testing.T) {
	result := Filter(Some(0), func(n int) bool {
		return true
	})

	expectSome(t, result)

	value, ok := result.TryGet()
	if !ok || value != 0 {
		t.Error("Filter should preserve Some(0)")
	}
}

type mockType struct {
	name string
	age  Option[int]
}

func expectNone[T any](t *testing.T, o Option[T], msg ...string) {
	t.Helper()

	if o.IsSome() {
		t.Fatalf("Expected None, got Some: %v", msg)
	}
}

func expectSome[T any](t *testing.T, o Option[T], msg ...string) {
	t.Helper()

	if o.IsNone() {
		t.Fatalf("Expected Some, got None: %v", msg)
	}
}

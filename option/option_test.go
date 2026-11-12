package option

import (
	"testing"
)

func TestSomeNone(t *testing.T) {
	expectSome(t, Some[uint8](1))
	expectSome(t, Some[uint16](1))
	expectSome(t, Some[uint32](1))
	expectSome(t, Some[uint64](1))
	expectSome(t, Some[int8](1))
	expectSome(t, Some[int16](1))
	expectSome(t, Some[int32](1))
	expectSome(t, Some[int64](1))
	expectSome(t, Some[float32](1))
	expectSome(t, Some[float64](1))
	expectSome(t, Some[complex64](1))
	expectSome(t, Some[complex128](1))
	expectSome(t, Some[byte](1))
	expectSome(t, Some('d'))
	expectSome(t, Some("danom"))
	expectSome(t, Some(mockType{name: "danom", age: Some(86)}))
	expectSome(t, mockType{name: "danom", age: Some(86)}.age)
}

func TestNone(t *testing.T) {
	expectNone(t, None[uint8]())
	expectNone(t, None[uint16]())
	expectNone(t, None[uint32]())
	expectNone(t, None[uint64]())
	expectNone(t, None[int8]())
	expectNone(t, None[int16]())
	expectNone(t, None[int32]())
	expectNone(t, None[int64]())
	expectNone(t, None[float32]())
	expectNone(t, None[float64]())
	expectNone(t, None[complex64]())
	expectNone(t, None[complex128]())
	expectNone(t, None[byte]())
	expectNone(t, None[rune]())
	expectNone(t, None[string]())
	expectNone(t, None[mockType]())
	expectNone(t, mockType{name: "danom", age: None[int]()}.age)
}

func TestSomeNil(t *testing.T) {
	var ptr *int

	o := Some(ptr)

	value, ok := o.Value()

	if !ok {
		t.Error("Some(nil) should still be Some")
	}

	if value != nil {
		t.Error("Expected nil pointer value")
	}
}

func TestValue(t *testing.T) {
	if _, ok := Some(1).Value(); !ok {
		t.Error("Expected value, but found none")
	}

	if _, ok := None[int]().Value(); ok {
		t.Error("Expected none, but found value")
	}
}

func TestValueZeroValue(t *testing.T) {
	var o Option[int]

	value, ok := o.Value()

	if ok {
		t.Error("Zero-value Option should be None")
	}

	if value != 0 {
		t.Error("Value should return zero value for None")
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

func TestOrElse(t *testing.T) {
	result := Some(1).OrElse(Some(2))
	expectSome(t, result)

	value := result.DefaultValue(0)
	if value != 1 {
		t.Error("OrElse should preserve existing Some value")
	}

	result = None[int]().OrElse(Some(2))
	expectSome(t, result)

	value = result.DefaultValue(0)
	if value != 2 {
		t.Error("OrElse should return fallback when None")
	}

	expectNone(t, None[int]().OrElse(None[int]()))
}

func TestOrElseWith(t *testing.T) {
	called := false

	result := Some(1).OrElseWith(func() Option[int] {
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

	result = None[int]().OrElseWith(func() Option[int] {
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

func TestTap(t *testing.T) {
	state := 0

	o := None[int]().Tap(func(n int) {
		state = n
	})

	expectNone(t, o)

	if state != 0 {
		t.Error("Tap should not invoke callback for None")
	}

	o = Some(86).Tap(func(n int) {
		state = n
	})

	expectSome(t, o)

	if state != 86 {
		t.Error("Tap should invoke callback with Some value")
	}
}

func TestFilter(t *testing.T) {
	result := Some(10).Filter(func(n int) bool {
		return n > 5
	})

	expectSome(t, result)

	value := result.DefaultValue(0)
	if value != 10 {
		t.Error("Filter should preserve value when predicate matches")
	}

	result = Some(3).Filter(func(n int) bool {
		return n > 5
	})

	expectNone(t, result)

	called := false

	result = None[int]().Filter(func(n int) bool {
		called = true
		return true
	})

	expectNone(t, result)

	if called {
		t.Error("Filter should not invoke predicate for None")
	}
}

func TestMatch(t *testing.T) {
	if Match(None[string](),
		func(x string) bool { return true },
		func() bool { return false }) {
		t.Error("Should not have invoked Some branch in match")
	}

	if Match(Some("danom"),
		func(x string) bool { return false },
		func() bool { return true }) {
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
		t.Error("Map function was not executed as expected")
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

func TestDo(t *testing.T) {
	state := 0
	Do(None[string](),
		func(x string) { state = 1 },
		func() {})

	if state == 1 {
		t.Error("Should not have invoked Some branch in do")
	}

	Do(Some[string]("danom"),
		func(x string) {},
		func() { state = 2 })

	if state == 2 {
		t.Error("Should not have invoked Some branch in do")
	}
}

func TestIter(t *testing.T) {
	state := 0

	Iter(None[int](), func(n int) {
		state = n
	})

	if state != 0 {
		t.Error("Iter should not invoke callback for None")
	}

	Iter(Some(86), func(n int) {
		state = n
	})

	if state != 86 {
		t.Error("Iter should invoke callback with Some value")
	}
}

func TestFilterPreservesZeroValue(t *testing.T) {
	result := Some(0).Filter(func(n int) bool {
		return true
	})

	expectSome(t, result)

	value, ok := result.Value()
	if !ok || value != 0 {
		t.Error("Filter should preserve Some(0)")
	}
}

func TestString(t *testing.T) {
	if Some(1).String() != "Some(1)" {
		t.Error("String failed")
	}

	if None[int]().String() != "None" {
		t.Error("String failed")
	}

	if Some(0).String() != "Some(0)" {
		t.Error("String failed")
	}

	if Some("a").String() != "Some(a)" {
		t.Error("String failed")
	}

	if Some("").String() != "Some()" {
		t.Error("String failed")
	}
}

func TestIsZero(t *testing.T) {
	if !(Option[int]{}).IsZero() {
		t.Error("zero-value Option should be zero")
	}
	if !None[int]().IsZero() {
		t.Error("None should be zero")
	}
	if Some(0).IsZero() {
		t.Error("Some(zero) should not be zero")
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

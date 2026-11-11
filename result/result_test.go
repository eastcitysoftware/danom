package result

import (
	"testing"
)

func TestOk(t *testing.T) {
	expectOk(t, Ok[uint8, byte](1))
	expectOk(t, Ok[uint16, byte](1))
	expectOk(t, Ok[uint32, byte](1))
	expectOk(t, Ok[uint64, byte](1))
	expectOk(t, Ok[int8, byte](1))
	expectOk(t, Ok[int16, byte](1))
	expectOk(t, Ok[int32, byte](1))
	expectOk(t, Ok[int64, byte](1))
	expectOk(t, Ok[float32, byte](1))
	expectOk(t, Ok[float64, byte](1))
	expectOk(t, Ok[complex64, byte](1))
	expectOk(t, Ok[complex128, byte](1))
	expectOk(t, Ok[byte, byte](1))
	expectOk(t, Ok[rune, byte]('d'))
	expectOk(t, Ok[string, byte]("danom"))
	// expectOk(t, Ok[struct{ name string }, byte](struct{name: "danom"}))
}

func TestError(t *testing.T) {
	expectError(t, Error[uint8, byte](0))
	expectError(t, Error[uint16, byte](0))
	expectError(t, Error[uint32, byte](0))
	expectError(t, Error[uint64, byte](0))
	expectError(t, Error[int8, byte](0))
	expectError(t, Error[int16, byte](0))
	expectError(t, Error[int32, byte](0))
	expectError(t, Error[int64, byte](0))
	expectError(t, Error[float32, byte](0))
	expectError(t, Error[float64, byte](0))
	expectError(t, Error[complex64, byte](0))
	expectError(t, Error[complex128, byte](0))
	expectError(t, Error[byte, byte](0))
	expectError(t, Error[rune, byte](0))
	expectError(t, Error[string, byte](0))
	// expectOk(t, Ok[struct{}, byte](struct{name: "danom"}))
}

// func TestOk(t *testing.T) {
// 	expectOk(t, Ok[int, string](1))
// }

// func TestError(t *testing.T) {
// 	expectError(t, Error[int, string]("danom"))
// }

func TestIsOk(t *testing.T) {
	if !Ok[int, string](1).IsOk() {
		t.Error("Expected Ok got Error")
	}
}

func TestIsError(t *testing.T) {
	if !Error[int, string]("danom").IsError() {
		t.Error("Expected Error got Ok")
	}
}

func TestValue(t *testing.T) {
	if value, ok := Ok[int, string](1).Value(); !ok || value != 1 {
		t.Error("Expected Ok got Error")
	}

	if value, ok := Error[int, string]("danom").Value(); ok || value != 0 {
		t.Error("Expected Error got Ok")
	}
}

func TestErrorValue(t *testing.T) {
	if value, ok := Error[int, string]("danom").ErrorValue(); !ok || value != "danom" {
		t.Error("Expected Error value")
	}

	if value, ok := Ok[int, string](1).ErrorValue(); ok || value != "" {
		t.Error("Expected no Error value")
	}
}

func TestDefaultValue(t *testing.T) {
	value := Ok[int, byte](1).DefaultValue(0)
	if value != 1 {
		t.Error("Default value should not have been used")
	}

	value = Error[int, byte](0).DefaultValue(1)

	if value != 1 {
		t.Error("Default value should have been used")
	}
}

func TestDefaultWith(t *testing.T) {
	value := Ok[int, byte](1).DefaultWith(func() int { return 0 })
	if value != 1 {
		t.Error("Default value should not have been used")
	}

	value = Error[int, byte](0).DefaultWith(func() int { return 1 })

	if value != 1 {
		t.Error("Default value should have been used")
	}
}

func TestDefaultWithLazy(t *testing.T) {
	called := false

	value := Ok[int, byte](1).DefaultWith(func() int {
		called = true
		return 2
	})

	if called {
		t.Error("DefaultWith should not invoke fallback when Ok")
	}

	if value != 1 {
		t.Error("Expected original value")
	}
}

func TestOrElse(t *testing.T) {
	result := Ok[int, string](1).OrElse(Ok[int, string](2))
	expectOk(t, result)

	value := result.DefaultValue(0)
	if value != 1 {
		t.Error("OrElse should preserve existing Ok value")
	}

	result = Error[int, string]("danom").OrElse(Ok[int, string](2))
	expectOk(t, result)

	value = result.DefaultValue(0)
	if value != 2 {
		t.Error("OrElse should return fallback when Error")
	}

	expectError(t, Error[int, string]("danom").OrElse(Error[int, string]("here")))
}

func TestOrElseWith(t *testing.T) {
	called := false

	result := Ok[int, string](1).OrElseWith(func() Result[int, string] {
		called = true
		return Ok[int, string](2)
	})

	expectOk(t, result)

	if called {
		t.Error("OrElseWith should not invoke fallback when Ok")
	}

	value := result.DefaultValue(0)
	if value != 1 {
		t.Error("OrElseWith should preserve existing Ok value")
	}

	result = Error[int, string]("danom").OrElseWith(func() Result[int, string] {
		called = true
		return Ok[int, string](2)
	})

	expectOk(t, result)

	if !called {
		t.Error("OrElseWith should invoke fallback when Error")
	}

	value = result.DefaultValue(0)
	if value != 2 {
		t.Error("OrElseWith should return fallback value")
	}
}

func TestTap(t *testing.T) {
	state := 0

	o := Error[int, string]("danom").Tap(func(n int) {
		state = n
	})

	expectError(t, o)

	if state != 0 {
		t.Error("Tap should not invoke callback for Error")
	}

	o = Ok[int, string](1).Tap(func(n int) {
		state = n
	})

	expectOk(t, o)

	if state != 1 {
		t.Error("Tap should invoke callback with Ok value")
	}
}

func TestFilter(t *testing.T) {
	result := Ok[int, byte](10).Filter(func(n int) bool {
		return n > 5
	}, 0)

	expectOk(t, result)

	value := result.DefaultValue(0)
	if value != 10 {
		t.Error("Filter should preserve value when predicate matches")
	}

	result = Ok[int, byte](3).Filter(func(n int) bool {
		return n > 5
	}, 0)

	expectError(t, result)

	called := false

	result = Error[int, byte](0).Filter(func(n int) bool {
		called = true
		return true
	}, 1)

	expectError(t, result)

	if called {
		t.Error("Filter should not invoke predicate for Error")
	}
}

func TestMatch(t *testing.T) {
	if Match(Ok[int, string](1),
		func(x int) bool { return false },
		func(e string) bool { return true }) {
		t.Error("Should not have invoked Error branch in match")
	}

	if Match(Error[int, string]("danom"),
		func(x int) bool { return true },
		func(e string) bool { return false }) {
		t.Error("Should not have invoked Ok branch in match")
	}
}

func TestBind(t *testing.T) {
	result := Bind(Ok[int, string](1), func(n int) Result[int, string] { return Ok[int, string](n + 1) })
	expectOk(t, result)

	value := Match(result, func(n int) int { return n }, func(e string) int { return 0 })
	if value != 2 {
		t.Error("Bind function was not executed as expected")
	}
}

func TestMap(t *testing.T) {
	result := Map(Ok[int, string](1), func(n int) int { return n + 1 })
	expectOk(t, result)

	value := Match(result, func(n int) int { return n }, func(e string) int { return 0 })
	if value != 2 {
		t.Error("Map function was not executed as expected")
	}
}

func TestMapError(t *testing.T) {
	result := MapError(Error[int, string]("danom"), func(s string) int { return -1 })
	expectError(t, result)

	value := Match(result, func(n int) int { return 0 }, func(e int) int { return e })
	if value != -1 {
		t.Error("MapError function was not executed as expected")
	}
}

func TestApply(t *testing.T) {
	result := Apply(Ok[int, string](1), Ok[func(n int) int, string](func(n int) int { return n + 1 }))
	expectOk(t, result)

	value := Match(result, func(n int) int { return n }, func(e string) int { return 0 })
	if value != 2 {
		t.Error("Bind function was not executed as expected")
	}

	result = Apply(Ok[int, string](1), Error[func(n int) int, string]("here"))
	expectError(t, result)

	value = Match(result, func(n int) int { return n }, func(e string) int { return 0 })
	if value != 0 {
		t.Error("Apply function should not return Ok")
	}
}

func TestDo(t *testing.T) {
	state := 0
	Do(Ok[int, string](1),
		func(x int) { state = 1 },
		func(e string) { state = 0 })

	if state != 1 {
		t.Error("Should not have invoked Error branch in do")
	}

	Do(Error[int, string]("danom"),
		func(x int) { state = 0 },
		func(e string) { state = 2 })

	if state != 2 {
		t.Error("Should not have invoked Ok branch in do")
	}
}

func TestIter(t *testing.T) {
	state := 0

	Iter(Error[int, string]("danom"), func(n int) {
		state = n
	})

	if state != 0 {
		t.Error("Iter should not invoke callback for Error")
	}

	Iter(Ok[int, string](1), func(n int) {
		state = n
	})

	if state != 1 {
		t.Error("Iter should invoke callback with Ok value")
	}
}

func expectOk[T any, E any](t *testing.T, r Result[T, E], msg ...string) {
	t.Helper()

	if r.IsError() {
		t.Fatalf("Expected Ok, got Error: %v", msg)
	}
}

func expectError[T any, E any](t *testing.T, r Result[T, E], msg ...string) {
	t.Helper()

	if r.IsOk() {
		t.Fatalf("Expected Error, got Ok: %v", msg)
	}
}

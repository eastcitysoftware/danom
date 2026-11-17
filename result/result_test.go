package result

import (
	"testing"
)

func TestOk(t *testing.T) {
	expectOk(t, Ok[string](1))
	expectOk(t, Ok[string]('d'))
	expectOk(t, Ok[string]("danom"))
}

func TestError(t *testing.T) {
	expectError(t, Error[int]("error"))
	expectError(t, Error[rune]("error"))
	expectError(t, Error[string]("error"))
}

func TestIsOk(t *testing.T) {
	if !Ok[string](1).IsOk() {
		t.Error("Expected Ok got Error")
	}
}

func TestIsError(t *testing.T) {
	if !Error[int]("danom").IsError() {
		t.Error("Expected Error got Ok")
	}
}

func TestValue(t *testing.T) {
	if value, ok := Ok[string](1).Value(); !ok || value != 1 {
		t.Error("Expected Ok got Error")
	}

	if value, ok := Error[int]("danom").Value(); ok || value != 0 {
		t.Error("Expected Error got Ok")
	}
}

func TestErrorValue(t *testing.T) {
	if value, ok := Error[int]("danom").ErrorValue(); !ok || value != "danom" {
		t.Error("Expected Error value")
	}

	if value, ok := Ok[string](1).ErrorValue(); ok || value != "" {
		t.Error("Expected no Error value")
	}
}

func TestDefaultValue(t *testing.T) {
	value := Ok[string](1).DefaultValue(0)
	if value != 1 {
		t.Error("Default value should not have been used")
	}

	value = Error[int](0).DefaultValue(1)

	if value != 1 {
		t.Error("Default value should have been used")
	}
}

func TestDefaultWith(t *testing.T) {
	value := Ok[string](1).DefaultWith(func() int { return 0 })
	if value != 1 {
		t.Error("Default value should not have been used")
	}

	value = Error[int](0).DefaultWith(func() int { return 1 })

	if value != 1 {
		t.Error("Default value should have been used")
	}
}

func TestDefaultWithLazy(t *testing.T) {
	called := false

	value := Ok[string](1).DefaultWith(func() int {
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
	result := Ok[string](1).OrElse(Ok[string](2))
	expectOk(t, result)

	value := result.DefaultValue(0)
	if value != 1 {
		t.Error("OrElse should preserve existing Ok value")
	}

	result = Error[int]("danom").OrElse(Ok[string](2))
	expectOk(t, result)

	value = result.DefaultValue(0)
	if value != 2 {
		t.Error("OrElse should return fallback when Error")
	}

	expectError(t, Error[int]("danom").OrElse(Error[int]("here")))
}

func TestOrElseWith(t *testing.T) {
	called := false

	result := Ok[string](1).OrElseWith(func() Result[int, string] {
		called = true
		return Ok[string](2)
	})

	expectOk(t, result)

	if called {
		t.Error("OrElseWith should not invoke fallback when Ok")
	}

	value := result.DefaultValue(0)
	if value != 1 {
		t.Error("OrElseWith should preserve existing Ok value")
	}

	result = Error[int]("danom").OrElseWith(func() Result[int, string] {
		called = true
		return Ok[string](2)
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

	o := Error[int]("danom").Tap(func(n int) {
		state = n
	})

	expectError(t, o)

	if state != 0 {
		t.Error("Tap should not invoke callback for Error")
	}

	o = Ok[string](1).Tap(func(n int) {
		state = n
	})

	expectOk(t, o)

	if state != 1 {
		t.Error("Tap should invoke callback with Ok value")
	}
}

func TestFilter(t *testing.T) {
	result := Ok[string](10).Filter(func(n int) bool {
		return n > 5
	}, "error")

	expectOk(t, result)

	value := result.DefaultValue(0)
	if value != 10 {
		t.Error("Filter should preserve value when predicate matches")
	}

	result = Ok[string](3).Filter(func(n int) bool {
		return n > 5
	}, "error")

	expectError(t, result)

	called := false

	result = Error[int]("error").Filter(func(n int) bool {
		called = true
		return true
	}, "error")

	expectError(t, result)

	if called {
		t.Error("Filter should not invoke predicate for Error")
	}
}

func TestDo(t *testing.T) {
	state := 0
	Ok[string](1).Do(
		func(x int) { state = 1 },
		func(e string) { state = 0 })

	if state != 1 {
		t.Error("Should not have invoked Error branch in do")
	}

	Error[int]("danom").Do(
		func(x int) { state = 0 },
		func(e string) { state = 2 })

	if state != 2 {
		t.Error("Should not have invoked Ok branch in do")
	}
}

func TestIter(t *testing.T) {
	state := 0

	Error[int]("danom").Iter(func(n int) {
		state = n
	})

	if state != 0 {
		t.Error("Iter should not invoke callback for Error")
	}

	Ok[string](1).Iter(func(n int) {
		state = n
	})

	if state != 1 {
		t.Error("Iter should invoke callback with Ok value")
	}
}

func TestMatch(t *testing.T) {
	if Match(Ok[string](1),
		func(x int) bool { return false },
		func(e string) bool { return true }) {
		t.Error("Should not have invoked Error branch in match")
	}

	if Match(Error[int]("danom"),
		func(x int) bool { return true },
		func(e string) bool { return false }) {
		t.Error("Should not have invoked Ok branch in match")
	}
}

func TestBind(t *testing.T) {
	result := Bind(Ok[string](1), func(n int) Result[int, string] { return Ok[string](n + 1) })
	expectOk(t, result)

	value := Match(result, func(n int) int { return n }, func(e string) int { return 0 })
	if value != 2 {
		t.Error("Bind function was not executed as expected")
	}
}

func TestMap(t *testing.T) {
	result := Map(Ok[string](1), func(n int) int { return n + 1 })
	expectOk(t, result)

	value := Match(result, func(n int) int { return n }, func(e string) int { return 0 })
	if value != 2 {
		t.Error("Map function was not executed as expected")
	}
}

func TestMapError(t *testing.T) {
	result := MapError(Error[int]("danom"), func(s string) int { return -1 })
	expectError(t, result)

	value := Match(result, func(n int) int { return 0 }, func(e int) int { return e })
	if value != -1 {
		t.Error("MapError function was not executed as expected")
	}
}

func TestApply(t *testing.T) {
	result := Apply(Ok[string](1), Ok[string](func(n int) int { return n + 1 }))
	expectOk(t, result)

	value := Match(result, func(n int) int { return n }, func(e string) int { return 0 })
	if value != 2 {
		t.Error("Bind function was not executed as expected")
	}

	result = Apply(Ok[string](1), Error[func(n int) int, string]("here"))
	expectError(t, result)

	value = Match(result, func(n int) int { return n }, func(e string) int { return 0 })
	if value != 0 {
		t.Error("Apply function should not return Ok")
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

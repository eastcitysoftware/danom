package result

import (
	"testing"
)

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

func TestTryGet(t *testing.T) {
	if _, _, ok := Ok[int, string](1).TryGet(); !ok {
		t.Error("Expected Ok got Error")
	}

	if _, _, ok := Error[int, string]("danom").TryGet(); ok {
		t.Error("Expected Error got Ok")
	}
}

func TestOk(t *testing.T) {
	expectOk(t, Ok[int, string](1))
}

func TestError(t *testing.T) {
	expectError(t, Error[int, string]("danom"))
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

func TestOrElse(t *testing.T) {
	result := OrElse(Ok[int, string](1), Ok[int, string](2))
	expectOk(t, result)

	value := result.DefaultValue(0)
	if value != 1 {
		t.Error("OrElse should preserve existing Ok value")
	}

	result = OrElse(Error[int, string]("danom"), Ok[int, string](2))
	expectOk(t, result)

	value = result.DefaultValue(0)
	if value != 2 {
		t.Error("OrElse should return fallback when Error")
	}

	expectError(t, OrElse(Error[int, string]("danom"), Error[int, string]("here")))
}

func TestOrElseWith(t *testing.T) {
	called := false

	result := OrElseWith(Ok[int, string](1), func() Result[int, string] {
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

	result = OrElseWith(Error[int, string]("danom"), func() Result[int, string] {
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

	o := Tap(Error[int, string]("danom"), func(n int) {
		state = n
	})

	expectError(t, o)

	if state != 0 {
		t.Error("Tap should not invoke callback for Error")
	}

	o = Tap(Ok[int, string](1), func(n int) {
		state = n
	})

	expectOk(t, o)

	if state != 1 {
		t.Error("Tap should invoke callback with Ok value")
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

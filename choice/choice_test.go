package choice

import (
	"testing"
)

func TestIsT1(t *testing.T) {
	if !OfT1[int, string](1).IsT1() {
		t.Error("Expected T1 got T2")
	}
}

func TestIsT2(t *testing.T) {
	if !OfT2[int, string]("danom").IsT2() {
		t.Error("Expected T2 got T1")
	}
}

func TestOfT1(t *testing.T) {
	expectT1(t, OfT1[int, string](1))
}

func TestOfT2(t *testing.T) {
	expectT2(t, OfT2[int, string]("danom"))
}

func TestTryGetT1(t *testing.T) {
	if _, ok := OfT1[int, string](1).TryGetT1(); !ok {
		t.Error("Expected T1, but found T2")
	}

	if _, ok := OfT2[int, string]("danom").TryGetT1(); ok {
		t.Error("Expected T2, but found T1")
	}
}

func TestTryGetT2(t *testing.T) {
	if _, ok := OfT2[int, string]("danom").TryGetT2(); !ok {
		t.Error("Expected T2, but found T1")
	}

	if _, ok := OfT1[int, string](1).TryGetT2(); ok {
		t.Error("Expected T1, but found T2")
	}
}

func TestSwap(t *testing.T) {
	if _, ok := OfT1[int, string](1).Swap().TryGetT2(); !ok {
		t.Error("After swap exptected T1 value")
	}

	if _, ok := OfT2[int, string]("danom").Swap().TryGetT1(); !ok {
		t.Error("After swap exptected T2 value")
	}
}

func TestMatch(t *testing.T) {
	if Match(OfT2[string](0),
		func(x string) bool { return true },
		func(x int) bool { return false },
	) {
		t.Error("Should not have invoked Some branch in match")
	}

	if Match(OfT1[string, int]("danom"),
		func(x string) bool { return false },
		func(x int) bool { return true },
	) {
		t.Error("Should not have invoked Some branch in match")
	}
}

func TestDo(t *testing.T) {
	state := 0
	Do(OfT2[string](0),
		func(x string) { state = 1 },
		func(x int) {})

	if state == 1 {
		t.Error("Should not have invoked Some branch in match")
	}

	Do(OfT1[string, int]("danom"),
		func(x string) {},
		func(x int) { state = 2 })

	if state == 2 {
		t.Error("Should not have invoked Some branch in match")
	}
}

func TestBindT1(t *testing.T) {
	result := BindT1(OfT1[int, string](1), func(n int) Choice[int, string] { return OfT1[int, string](n + 1) })
	expectT1(t, result)

	value := Match(result, func(n int) int { return n }, func(s string) int { return 0 })
	if value != 2 {
		t.Error("BindT1 function was not executed as expected")
	}
}

func TestMapT1(t *testing.T) {
	result := MapT1(OfT1[int, string](1), func(n int) int { return n + 1 })
	expectT1(t, result)

	value := Match(result, func(n int) int { return n }, func(s string) int { return 0 })
	if value != 2 {
		t.Error("MapT1 function was not executed as expected")
	}
}

func TestBindT2(t *testing.T) {
	result := BindT2(OfT2[string, int](1), func(n int) Choice[string, int] { return OfT2[string, int](n + 1) })
	expectT2(t, result)

	value := Match(result, func(s string) int { return 0 }, func(n int) int { return n })
	if value != 2 {
		t.Error("BindT2 function was not executed as expected")
	}
}

func TestMapT2(t *testing.T) {
	result := MapT2(OfT2[string, int](1), func(n int) int { return n + 1 })
	expectT2(t, result)

	value := Match(result, func(s string) int { return 0 }, func(n int) int { return n })
	if value != 2 {
		t.Error("MapT2 function was not executed as expected")
	}
}

func expectT1[T1 any, T2 any](t *testing.T, o Choice[T1, T2], msg ...string) {
	t.Helper()

	if o.IsT2() {
		t.Errorf("Expected T1, got T2: %v", msg)
	}
}

func expectT2[T1 any, T2 any](t *testing.T, o Choice[T1, T2], msg ...string) {
	t.Helper()

	if o.IsT1() {
		t.Errorf("Expected T2, got T1: %v", msg)
	}
}

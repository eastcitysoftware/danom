package check

import (
	"fmt"
	"testing"

	"github.com/eastcitysoftware/danom/option"
)

func TestValidate(t *testing.T) {
	validator := NewValidator[mockCheckType](
		NewRule("name", func(x mockCheckType) string { return x.name },
			StrNotEmpty(),
			StrGreaterThan(2),
			NoneOf(StrEqualsFold("here"), Equals("there"))),
		NewRule("age", func(x mockCheckType) int { return x.age },
			GreaterThanOrEqual(1)),
		NewRule("lucky number", func(x mockCheckType) option.Option[int] { return x.luckyNumber },
			Optional(GreaterThan(0))))

	// use method form
	if validator.Validate(mockCheckType{"danom", 1, option.None[int]()}).IsError() {
		t.Error("Provided valid input, but validation failed")
	}

	// use function form
	if validator(mockCheckType{"", 0, option.Some(-1)}).IsOk() {
		t.Error("Provided invalid input, but validation succeeded")
	}
}

func TestIt(t *testing.T) {
	it := NewRule("name", func(x mockCheckType) string { return x.name },
		StrNotEmpty())

	mock := mockCheckType{"danom", 1, option.Some(1)}
	result := it(mock)
	if len(result.Errors) > 0 {
		t.Error("Provided valid input, but got error")
	}

	mock = mockCheckType{"", 1, option.Some(-1)}
	result = it(mock)
	if len(result.Errors) == 0 {
		t.Error("Provided invalid input, but got no error")
	}
}

func TestNewPredicate(t *testing.T) {
	p := NewFact(func(x int) bool { return x > 1 }, NewMessage("Must be gt 1"))
	if len(p(0)) == 0 {
		t.Error("Predicate should produce one Message, got zero")
	}

	if len(p(2)) > 0 {
		t.Error("Predicate should produce zero Message, got some")
	}
}

func TestNewMessage(t *testing.T) {
	m := NewMessage("%s")

	if m("danom") != "danom" {
		t.Error("Message not working as expected")
	}

	m = NewMessage(fmt.Sprintf("%%s %s", "test"))
	expected := "danom test"
	if got := m("danom"); got != expected {
		t.Errorf("Message with escaping not working as expected, got '%s' expected '%s'", got, expected)
	}
}

type mockCheckType struct {
	name        string
	age         int
	luckyNumber option.Option[int]
}

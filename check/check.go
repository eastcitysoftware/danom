package check

import (
	"fmt"
	"strings"

	"github.com/eastcitysoftware/danom/result"
)

// Error contains the validation errors for a named rule.
type Error struct {
	Name   string
	Errors []string
}

// Errors is the collection of validation errors returned by a Validator.
type Errors []Error

// Message formats a validation message using the name of the value being checked.
type Message func(string) string

// Fact checks a value and returns messages for any failed conditions.
type Fact[T any] func(T) []Message

// Rule validates an input and returns its named Error.
type Rule[T any] func(T) Error

// Validator applies rules to an input and returns the input or its validation errors.
type Validator[T any] func(T) result.Result[T, Errors]

// NewValidator creates a Validator that evaluates each rule and collects its errors.
func NewValidator[T any](rules ...Rule[T]) Validator[T] {
	return func(input T) result.Result[T, Errors] {
		var errors Errors

		for _, rule := range rules {
			e := rule(input)

			if len(e.Errors) > 0 {
				errors = append(errors, e)
			}
		}

		if len(errors) > 0 {
			return result.Error[T, Errors](errors)
		}

		return result.Ok[T, Errors](input)
	}
}

// NewRule creates a named Rule that selects a value and checks it against facts.
func NewRule[T, U any](name string, selector func(T) U, preds ...Fact[U]) Rule[T] {
	return func(input T) Error {
		e := Error{Name: name}

		value := selector(input)

		for _, pred := range preds {
			for _, message := range pred(value) {
				e.Errors = append(e.Errors, message(name))
			}
		}

		return e
	}
}

// NewFact creates a Fact that returns message when pred is false.
func NewFact[T any](pred func(T) bool, message Message) Fact[T] {
	return func(input T) []Message {
		if !pred(input) {
			return []Message{message}
		}

		return nil
	}
}

// NewMessage creates a Message using fmt.Sprintf with the rule name as its first argument.
func NewMessage(format string, args ...any) Message {
	return func(name string) string {
		values := make([]any, 1, len(args)+1)
		values[0] = name
		values = append(values, args...)
		return fmt.Sprintf(format, values...)
	}
}

// Validate applies the Validator to input.
func (v Validator[T]) Validate(input T) result.Result[T, Errors] {
	return v(input)
}

// String formats the validation error with its rule name and messages.
func (e Error) String() string {
	return e.Name + " - " + strings.Join(e.Errors, " | ")
}

// String formats the validation errors as a bracketed list.
func (e Errors) String() string {
	parts := make([]string, len(e))

	for i, err := range e {
		parts[i] = err.String()
	}

	return "[ " + strings.Join(parts, "; ") + " ]"
}

// Error implements the error interface.
func (e Errors) Error() string {
	return e.String()
}

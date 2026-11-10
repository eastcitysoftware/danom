package check

import (
	"fmt"
	"strings"

	"github.com/eastcitysoftware/danom/result"
)

type Error struct {
	Name   string
	Errors []string
}

type Errors []Error

type Message func(string) string

type Fact[T any] func(T) []Message

type Rule[T any] func(T) Error

type Validator[T any] func(T) result.Result[T, Errors]

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

func NewFact[T any](pred func(T) bool, message Message) Fact[T] {
	return func(input T) []Message {
		if !pred(input) {
			return []Message{message}
		}

		return nil
	}
}

func NewMessage(format string, args ...any) Message {
	return func(name string) string {
		values := make([]any, 1, len(args)+1)
		values[0] = name
		values = append(values, args...)
		return fmt.Sprintf(format, values...)
	}
}

func (v Validator[T]) Validate(input T) result.Result[T, Errors] {
	return v(input)
}

func (e Error) String() string {
	return e.Name + " - " + strings.Join(e.Errors, " | ")
}

func (e Errors) String() string {
	parts := make([]string, len(e))

	for i, err := range e {
		parts[i] = err.String()
	}

	return "[ " + strings.Join(parts, "; ") + " ]"
}

func (e Errors) Error() string {
	return e.String()
}

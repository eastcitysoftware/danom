// Package convert parses strings into Option values. A valid input produces
// Some(value); invalid input produces None, discarding the parse error.
// Use package parse when the error is needed.
package convert

import (
	"time"

	"github.com/eastcitysoftware/danom/option"
	"github.com/eastcitysoftware/danom/parse"
	"github.com/eastcitysoftware/danom/result"
)

// Bool converts x to a boolean using strconv.ParseBool.
func Bool(x string) option.Option[bool] {
	return orNone(parse.Bool(x))
}

// Float32 converts x to a 32-bit floating-point number using strconv.ParseFloat.
func Float32(x string) option.Option[float32] {
	return orNone(parse.Float32(x))
}

// Float64 converts x to a 64-bit floating-point number using strconv.ParseFloat.
func Float64(x string) option.Option[float64] {
	return orNone(parse.Float64(x))
}

// Int8 converts x to a base-10, 8-bit signed integer.
func Int8(x string) option.Option[int8] {
	return orNone(parse.Int8(x))
}

// Int16 converts x to a base-10, 16-bit signed integer.
func Int16(x string) option.Option[int16] {
	return orNone(parse.Int16(x))
}

// Int32 converts x to a base-10, 32-bit signed integer.
func Int32(x string) option.Option[int32] {
	return orNone(parse.Int32(x))
}

// Int64 converts x to a base-10, 64-bit signed integer.
func Int64(x string) option.Option[int64] {
	return orNone(parse.Int64(x))
}

// Uint8 converts x to a base-10, 8-bit unsigned integer.
func Uint8(x string) option.Option[uint8] {
	return orNone(parse.Uint8(x))
}

// Uint16 converts x to a base-10, 16-bit unsigned integer.
func Uint16(x string) option.Option[uint16] {
	return orNone(parse.Uint16(x))
}

// Uint32 converts x to a base-10, 32-bit unsigned integer.
func Uint32(x string) option.Option[uint32] {
	return orNone(parse.Uint32(x))
}

// Uint64 converts x to a base-10, 64-bit unsigned integer.
func Uint64(x string) option.Option[uint64] {
	return orNone(parse.Uint64(x))
}

// Complex64 converts x to a 64-bit complex number using strconv.ParseComplex.
func Complex64(x string) option.Option[complex64] {
	return orNone(parse.Complex64(x))
}

// Complex128 converts x to a 128-bit complex number using strconv.ParseComplex.
func Complex128(x string) option.Option[complex128] {
	return orNone(parse.Complex128(x))
}

// Duration converts x to a duration using time.ParseDuration.
func Duration(x string) option.Option[time.Duration] {
	return orNone(parse.Duration(x))
}

// Time converts x to a time using the given Go time layout.
func Time(x string, layout string) option.Option[time.Time] {
	return orNone(parse.Time(x, layout))
}

// TimeRFC3339 converts x to a time using the RFC3339 layout.
func TimeRFC3339(x string) option.Option[time.Time] {
	return orNone(parse.TimeRFC3339(x))
}

func orNone[T any](r result.Result[T, error]) option.Option[T] {
	return result.Map(r, option.Some).DefaultWith(option.None)
}

// Package parse wraps Go's standard parsing functions in Result values.
// Successful parses contain the parsed value; invalid input preserves the
// underlying error.
package parse

import (
	"strconv"
	"time"

	"github.com/eastcitysoftware/danom/result"
)

// Bool parses x as a boolean using strconv.ParseBool.
func Bool(x string) result.Result[bool, error] {
	if b, e := strconv.ParseBool(x); e == nil {
		return result.Ok[bool, error](b)
	} else {
		return result.Error[bool](e)
	}
}

// Float32 parses x as a 32-bit floating-point number using strconv.ParseFloat.
func Float32(x string) result.Result[float32, error] {
	if i, e := strconv.ParseFloat(x, 32); e == nil {
		return result.Ok[float32, error](float32(i))
	} else {
		return result.Error[float32](e)
	}
}

// Float64 parses x as a 64-bit floating-point number using strconv.ParseFloat.
func Float64(x string) result.Result[float64, error] {
	if i, e := strconv.ParseFloat(x, 64); e == nil {
		return result.Ok[float64, error](i)
	} else {
		return result.Error[float64](e)
	}
}

// Int8 parses x as a base-10, 8-bit signed integer.
func Int8(x string) result.Result[int8, error] {
	if i, e := strconv.ParseInt(x, 10, 8); e == nil {
		return result.Ok[int8, error](int8(i))
	} else {
		return result.Error[int8](e)
	}
}

// Int16 parses x as a base-10, 16-bit signed integer.
func Int16(x string) result.Result[int16, error] {
	if i, e := strconv.ParseInt(x, 10, 16); e == nil {
		return result.Ok[int16, error](int16(i))
	} else {
		return result.Error[int16](e)
	}
}

// Int32 parses x as a base-10, 32-bit signed integer.
func Int32(x string) result.Result[int32, error] {
	if i, e := strconv.ParseInt(x, 10, 32); e == nil {
		return result.Ok[int32, error](int32(i))
	} else {
		return result.Error[int32](e)
	}
}

// Int64 parses x as a base-10, 64-bit signed integer.
func Int64(x string) result.Result[int64, error] {
	if i, e := strconv.ParseInt(x, 10, 64); e == nil {
		return result.Ok[int64, error](i)
	} else {
		return result.Error[int64](e)
	}
}

// Uint8 parses x as a base-10, 8-bit unsigned integer.
func Uint8(x string) result.Result[uint8, error] {
	if i, e := strconv.ParseUint(x, 10, 8); e == nil {
		return result.Ok[uint8, error](uint8(i))
	} else {
		return result.Error[uint8](e)
	}
}

// Uint16 parses x as a base-10, 16-bit unsigned integer.
func Uint16(x string) result.Result[uint16, error] {
	if i, e := strconv.ParseUint(x, 10, 16); e == nil {
		return result.Ok[uint16, error](uint16(i))
	} else {
		return result.Error[uint16](e)
	}
}

// Uint32 parses x as a base-10, 32-bit unsigned integer.
func Uint32(x string) result.Result[uint32, error] {
	if i, e := strconv.ParseUint(x, 10, 32); e == nil {
		return result.Ok[uint32, error](uint32(i))
	} else {
		return result.Error[uint32](e)
	}
}

// Uint64 parses x as a base-10, 64-bit unsigned integer.
func Uint64(x string) result.Result[uint64, error] {
	if i, e := strconv.ParseUint(x, 10, 64); e == nil {
		return result.Ok[uint64, error](uint64(i))
	} else {
		return result.Error[uint64](e)
	}
}

// Complex64 parses x as a 64-bit complex number using strconv.ParseComplex.
func Complex64(x string) result.Result[complex64, error] {
	if i, e := strconv.ParseComplex(x, 64); e == nil {
		return result.Ok[complex64, error](complex64(i))
	} else {
		return result.Error[complex64](e)
	}
}

// Complex128 parses x as a 128-bit complex number using strconv.ParseComplex.
func Complex128(x string) result.Result[complex128, error] {
	if i, e := strconv.ParseComplex(x, 128); e == nil {
		return result.Ok[complex128, error](i)
	} else {
		return result.Error[complex128](e)
	}
}

// Duration parses x as a duration using time.ParseDuration.
func Duration(x string) result.Result[time.Duration, error] {
	if d, e := time.ParseDuration(x); e == nil {
		return result.Ok[time.Duration, error](d)
	} else {
		return result.Error[time.Duration](e)
	}
}

// Time parses x using the given Go time layout.
func Time(x string, layout string) result.Result[time.Time, error] {
	if t, e := time.Parse(layout, x); e == nil {
		return result.Ok[time.Time, error](t)
	} else {
		return result.Error[time.Time](e)
	}
}

// TimeRFC3339 parses x using the RFC3339 time layout.
func TimeRFC3339(x string) result.Result[time.Time, error] {
	return Time(x, time.RFC3339)
}

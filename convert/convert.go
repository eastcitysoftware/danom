package convert

import (
	"time"

	"github.com/eastcitysoftware/danom/option"
	"github.com/eastcitysoftware/danom/parse"
	"github.com/eastcitysoftware/danom/result"
)

func Bool(x string) option.Option[bool] {
	return orNone(parse.Bool(x))
}

func Float32(x string) option.Option[float32] {
	return orNone(parse.Float32(x))
}

func Float64(x string) option.Option[float64] {
	return orNone(parse.Float64(x))
}

func Int8(x string) option.Option[int8] {
	return orNone(parse.Int8(x))
}

func Int16(x string) option.Option[int16] {
	return orNone(parse.Int16(x))
}

func Int32(x string) option.Option[int32] {
	return orNone(parse.Int32(x))
}

func Int64(x string) option.Option[int64] {
	return orNone(parse.Int64(x))
}

func Uint8(x string) option.Option[uint8] {
	return orNone(parse.Uint8(x))
}

func Uint16(x string) option.Option[uint16] {
	return orNone(parse.Uint16(x))
}

func Uint32(x string) option.Option[uint32] {
	return orNone(parse.Uint32(x))
}

func Uint64(x string) option.Option[uint64] {
	return orNone(parse.Uint64(x))
}

func Complex64(x string) option.Option[complex64] {
	return orNone(parse.Complex64(x))
}

func Complex128(x string) option.Option[complex128] {
	return orNone(parse.Complex128(x))
}

func Duration(x string) option.Option[time.Duration] {
	return orNone(parse.Duration(x))
}

func Time(x string, layout string) option.Option[time.Time] {
	return orNone(parse.Time(x, layout))
}

func TimeRFC3339(x string) option.Option[time.Time] {
	return orNone(parse.TimeRFC3339(x))
}

func orNone[T any](r result.Result[T, error]) option.Option[T] {
	return result.Map(r, option.Some).DefaultWith(option.None)
}

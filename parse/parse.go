package parse

import (
	"strconv"
	"time"

	"github.com/eastcitysoftware/danom/result"
)

func Bool(x string) result.Result[bool, error] {
	if b, e := strconv.ParseBool(x); e == nil {
		return result.Ok[bool, error](b)
	} else {
		return result.Error[bool](e)
	}
}

func Float32(x string) result.Result[float32, error] {
	if i, e := strconv.ParseFloat(x, 32); e == nil {
		return result.Ok[float32, error](float32(i))
	} else {
		return result.Error[float32](e)
	}
}

func Float64(x string) result.Result[float64, error] {
	if i, e := strconv.ParseFloat(x, 64); e == nil {
		return result.Ok[float64, error](i)
	} else {
		return result.Error[float64](e)
	}
}

func Int8(x string) result.Result[int8, error] {
	if i, e := strconv.ParseInt(x, 10, 8); e == nil {
		return result.Ok[int8, error](int8(i))
	} else {
		return result.Error[int8](e)
	}
}

func Int16(x string) result.Result[int16, error] {
	if i, e := strconv.ParseInt(x, 10, 16); e == nil {
		return result.Ok[int16, error](int16(i))
	} else {
		return result.Error[int16](e)
	}
}

func Int32(x string) result.Result[int32, error] {
	if i, e := strconv.ParseInt(x, 10, 32); e == nil {
		return result.Ok[int32, error](int32(i))
	} else {
		return result.Error[int32](e)
	}
}

func Int64(x string) result.Result[int64, error] {
	if i, e := strconv.ParseInt(x, 10, 64); e == nil {
		return result.Ok[int64, error](i)
	} else {
		return result.Error[int64](e)
	}
}

func Uint8(x string) result.Result[uint8, error] {
	if i, e := strconv.ParseUint(x, 10, 8); e == nil {
		return result.Ok[uint8, error](uint8(i))
	} else {
		return result.Error[uint8](e)
	}
}

func Uint16(x string) result.Result[uint16, error] {
	if i, e := strconv.ParseUint(x, 10, 16); e == nil {
		return result.Ok[uint16, error](uint16(i))
	} else {
		return result.Error[uint16](e)
	}
}

func Uint32(x string) result.Result[uint32, error] {
	if i, e := strconv.ParseUint(x, 10, 32); e == nil {
		return result.Ok[uint32, error](uint32(i))
	} else {
		return result.Error[uint32](e)
	}
}

func Uint64(x string) result.Result[uint64, error] {
	if i, e := strconv.ParseUint(x, 10, 64); e == nil {
		return result.Ok[uint64, error](uint64(i))
	} else {
		return result.Error[uint64](e)
	}
}

func Complex64(x string) result.Result[complex64, error] {
	if i, e := strconv.ParseComplex(x, 64); e == nil {
		return result.Ok[complex64, error](complex64(i))
	} else {
		return result.Error[complex64](e)
	}
}

func Complex128(x string) result.Result[complex128, error] {
	if i, e := strconv.ParseComplex(x, 128); e == nil {
		return result.Ok[complex128, error](i)
	} else {
		return result.Error[complex128](e)
	}
}

func Duration(x string) result.Result[time.Duration, error] {
	if d, e := time.ParseDuration(x); e == nil {
		return result.Ok[time.Duration, error](d)
	} else {
		return result.Error[time.Duration](e)
	}
}

func Time(x string, layout string) result.Result[time.Time, error] {
	if t, e := time.Parse(layout, x); e == nil {
		return result.Ok[time.Time, error](t)
	} else {
		return result.Error[time.Time](e)
	}
}

func TimeRFC3339(x string) result.Result[time.Time, error] {
	return Time(x, time.RFC3339)
}

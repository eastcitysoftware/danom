package convert

import (
	"testing"
	"time"

	"github.com/eastcitysoftware/danom/option"
)

func TestBool(t *testing.T) {
	assertConverted(t, Bool("true"), true)
	assertConverted(t, Bool("FALSE"), false)
	assertNone(t, Bool("truthy"))
}

func TestFloat32(t *testing.T) {
	assertConverted(t, Float32("1.25"), float32(1.25))
	assertNone(t, Float32("1e100"))
}

func TestFloat64(t *testing.T) {
	assertConverted(t, Float64("-1.25"), -1.25)
	assertNone(t, Float64("1.2.3"))
}

func TestInt8(t *testing.T) {
	assertConverted(t, Int8("-128"), int8(-128))
	assertConverted(t, Int8("127"), int8(127))
	assertNone(t, Int8("128"))
}

func TestInt16(t *testing.T) {
	assertConverted(t, Int16("-32768"), int16(-32768))
	assertConverted(t, Int16("32767"), int16(32767))
	assertNone(t, Int16("32768"))
}

func TestInt32(t *testing.T) {
	assertConverted(t, Int32("-2147483648"), int32(-2147483648))
	assertConverted(t, Int32("2147483647"), int32(2147483647))
	assertNone(t, Int32("2147483648"))
}

func TestInt64(t *testing.T) {
	assertConverted(t, Int64("-9223372036854775808"), int64(-9223372036854775808))
	assertConverted(t, Int64("9223372036854775807"), int64(9223372036854775807))
	assertNone(t, Int64("9223372036854775808"))
}

func TestUint8(t *testing.T) {
	assertConverted(t, Uint8("0"), uint8(0))
	assertConverted(t, Uint8("255"), uint8(255))
	assertNone(t, Uint8("-1"))
	assertNone(t, Uint8("256"))
}

func TestUint16(t *testing.T) {
	assertConverted(t, Uint16("0"), uint16(0))
	assertConverted(t, Uint16("65535"), uint16(65535))
	assertNone(t, Uint16("-1"))
	assertNone(t, Uint16("65536"))
}

func TestUint32(t *testing.T) {
	assertConverted(t, Uint32("0"), uint32(0))
	assertConverted(t, Uint32("4294967295"), uint32(4294967295))
	assertNone(t, Uint32("-1"))
	assertNone(t, Uint32("4294967296"))
}

func TestUint64(t *testing.T) {
	assertConverted(t, Uint64("0"), uint64(0))
	assertConverted(t, Uint64("18446744073709551615"), uint64(18446744073709551615))
	assertNone(t, Uint64("-1"))
	assertNone(t, Uint64("18446744073709551616"))
}

func TestComplex64(t *testing.T) {
	assertConverted(t, Complex64("1+2i"), complex64(1+2i))
	assertNone(t, Complex64("1+2"))
}

func TestComplex128(t *testing.T) {
	assertConverted(t, Complex128("-3.5+2.25i"), complex(-3.5, 2.25))
	assertNone(t, Complex128("not complex"))
}

func TestDuration(t *testing.T) {
	assertConverted(t, Duration("1h2m3.5s"), time.Hour+2*time.Minute+3500*time.Millisecond)
	assertNone(t, Duration("1 hour"))
}

func TestTime(t *testing.T) {
	const layout = "2006-01-02"
	assertConvertedTime(t, Time("2024-03-10", layout), time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC))
	assertNone(t, Time("not a date", layout))
}

func TestTimeRFC3339(t *testing.T) {
	assertConvertedTime(
		t,
		TimeRFC3339("2024-03-10T12:30:45Z"),
		time.Date(2024, 3, 10, 12, 30, 45, 0, time.UTC),
	)
	assertNone(t, TimeRFC3339("2024-03-10 12:30:45"))
}

func assertConverted[T comparable](t *testing.T, got option.Option[T], want T) {
	t.Helper()

	value, ok := got.Value()
	if !ok {
		t.Fatal("expected a converted value")
	}
	if value != want {
		t.Errorf("converted value = %v, want %v", value, want)
	}
}

func assertConvertedTime(t *testing.T, got option.Option[time.Time], want time.Time) {
	t.Helper()

	value, ok := got.Value()
	if !ok {
		t.Fatal("expected a converted time")
	}
	if !value.Equal(want) {
		t.Errorf("converted time = %v, want %v", value, want)
	}
}

func assertNone[T any](t *testing.T, got option.Option[T]) {
	t.Helper()

	if _, ok := got.Value(); ok {
		t.Fatal("expected no value for invalid input")
	}
}

package parse

import (
	"testing"
	"time"

	"github.com/eastcitysoftware/danom/result"
)

func TestBool(t *testing.T) {
	assertParsed(t, Bool("true"), true)
	assertParsed(t, Bool("FALSE"), false)
	assertParsed(t, Bool("1"), true)
	assertParseError(t, Bool("truthy"))
}

func TestFloat32(t *testing.T) {
	assertParsed(t, Float32("1.25"), float32(1.25))
	assertParseError(t, Float32("1e100"))
}

func TestFloat64(t *testing.T) {
	assertParsed(t, Float64("-1.25"), -1.25)
	assertParseError(t, Float64("1.2.3"))
}

func TestInt8(t *testing.T) {
	assertParsed(t, Int8("-128"), int8(-128))
	assertParsed(t, Int8("127"), int8(127))
	assertParseError(t, Int8("128"))
}

func TestInt16(t *testing.T) {
	assertParsed(t, Int16("-32768"), int16(-32768))
	assertParsed(t, Int16("32767"), int16(32767))
	assertParseError(t, Int16("32768"))
}

func TestInt32(t *testing.T) {
	assertParsed(t, Int32("-2147483648"), int32(-2147483648))
	assertParsed(t, Int32("2147483647"), int32(2147483647))
	assertParseError(t, Int32("2147483648"))
}

func TestInt64(t *testing.T) {
	assertParsed(t, Int64("-9223372036854775808"), int64(-9223372036854775808))
	assertParsed(t, Int64("9223372036854775807"), int64(9223372036854775807))
	assertParseError(t, Int64("9223372036854775808"))
}

func TestUint8(t *testing.T) {
	assertParsed(t, Uint8("0"), uint8(0))
	assertParsed(t, Uint8("255"), uint8(255))
	assertParseError(t, Uint8("-1"))
	assertParseError(t, Uint8("256"))
}

func TestUint16(t *testing.T) {
	assertParsed(t, Uint16("0"), uint16(0))
	assertParsed(t, Uint16("65535"), uint16(65535))
	assertParseError(t, Uint16("-1"))
	assertParseError(t, Uint16("65536"))
}

func TestUint32(t *testing.T) {
	assertParsed(t, Uint32("0"), uint32(0))
	assertParsed(t, Uint32("4294967295"), uint32(4294967295))
	assertParseError(t, Uint32("-1"))
	assertParseError(t, Uint32("4294967296"))
}

func TestUint64(t *testing.T) {
	assertParsed(t, Uint64("0"), uint64(0))
	assertParsed(t, Uint64("18446744073709551615"), uint64(18446744073709551615))
	assertParseError(t, Uint64("-1"))
	assertParseError(t, Uint64("18446744073709551616"))
}

func TestComplex64(t *testing.T) {
	assertParsed(t, Complex64("1+2i"), complex64(1+2i))
	assertParseError(t, Complex64("1+2"))
}

func TestComplex128(t *testing.T) {
	assertParsed(t, Complex128("-3.5+2.25i"), complex(-3.5, 2.25))
	assertParseError(t, Complex128("not complex"))
}

func TestDuration(t *testing.T) {
	assertParsed(t, Duration("1h2m3.5s"), time.Hour+2*time.Minute+3500*time.Millisecond)
	assertParseError(t, Duration("1 hour"))
}

func TestTime(t *testing.T) {
	const layout = "2006-01-02"
	assertParsed(t, Time("2024-03-10", layout), time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC))
	assertParseError(t, Time("not a date", layout))
}

func TestTimeRFC3339(t *testing.T) {
	got, ok := TimeRFC3339("2024-03-10T12:30:45Z").Value()
	if !ok {
		t.Fatal("expected a parsed time")
	}

	want := time.Date(2024, 3, 10, 12, 30, 45, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("TimeRFC3339() = %v, want %v", got, want)
	}

	assertParseError(t, TimeRFC3339("2024-03-10 12:30:45"))
}

func assertParsed[T comparable](t *testing.T, got result.Result[T, error], want T) {
	t.Helper()

	value, ok := got.Value()
	if !ok {
		t.Fatal("expected a successful parse")
	}
	if value != want {
		t.Errorf("parsed value = %v, want %v", value, want)
	}
}

func assertParseError[T any](t *testing.T, got result.Result[T, error]) {
	t.Helper()

	err, ok := got.ErrorValue()
	if !ok || err == nil {
		t.Fatal("expected a parse error")
	}
}

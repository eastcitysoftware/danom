package check

import (
	"cmp"
	"net"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/eastcitysoftware/danom/option"
)

var (
	alphaPattern    = regexp.MustCompile("^[a-zA-Z]+$")
	alphaNumPattern = regexp.MustCompile("^[a-zA-Z0-9]+$")
)

// When runs facts when cond succeeds and returns their messages.
func When[T any](cond Fact[T], facts ...Fact[T]) Fact[T] {
	return func(x T) []Message {
		var m []Message
		if len(cond(x)) == 0 {
			for _, f := range facts {
				fm := f(x)
				if len(fm) > 0 {
					m = append(m, fm...)
				}
			}
		}

		if len(m) > 0 {
			return m
		}

		return nil
	}
}

// NoneOf succeeds when every fact fails.
func NoneOf[T any](facts ...Fact[T]) Fact[T] {
	return func(x T) []Message {
		for _, f := range facts {
			if len(f(x)) == 0 {
				return []Message{NewMessage("'%s' must not satisfy any of the conditions")}
			}
		}
		return nil
	}
}

// AnyOf succeeds when at least one fact succeeds.
func AnyOf[T any](facts ...Fact[T]) Fact[T] {
	return func(x T) []Message {
		for _, f := range facts {
			if len(f(x)) == 0 {
				return nil
			}
		}
		return []Message{NewMessage("'%s' must satisfy at least one condition")}
	}
}

// Equals succeeds when a value equals e.
func Equals[C comparable](e C) Fact[C] {
	return NewFact(
		func(x C) bool { return x == e },
		NewMessage("'%s' must be equal to %v", e))
}

// NotEquals succeeds when a value does not equal e.
func NotEquals[C comparable](e C) Fact[C] {
	return NewFact(
		func(x C) bool { return x != e },
		NewMessage("'%s' must not be equal to %v", e))
}

// True succeeds when a value is true.
func True() Fact[bool] {
	return NewFact(
		func(x bool) bool { return x == true },
		NewMessage("'%s' must be true"))
}

// False succeeds when a value is false.
func False() Fact[bool] {
	return NewFact(
		func(x bool) bool { return x != true },
		NewMessage("'%s' must be false"))
}

// Number includes the built-in integer, unsigned integer, and floating-point types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Zero succeeds when a numeric value is zero.
func Zero[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x == 0 },
		NewMessage("'%s' must be equal to 0"))
}

// Positive succeeds when a numeric value is greater than zero.
func Positive[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x > 0 },
		NewMessage("'%s' must be positive"))
}

// Negative succeeds when a numeric value is less than zero.
func Negative[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x < 0 },
		NewMessage("'%s' must be negative"))
}

// NonZero succeeds when a numeric value is not zero.
func NonZero[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x != 0 },
		NewMessage("'%s' must not be zero"))
}

// NonPositive succeeds when a numeric value is less than or equal to zero.
func NonPositive[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x <= 0 },
		NewMessage("'%s' must be non-positive"))
}

// NonNegative succeeds when a numeric value is greater than or equal to zero.
func NonNegative[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x >= 0 },
		NewMessage("'%s' must be non-negative"))
}

// Between succeeds when an ordered value is between min and max, inclusive.
func Between[N cmp.Ordered](min, max N) Fact[N] {
	return NewFact(
		func(x N) bool { return x >= min && x <= max },
		NewMessage("'%s' must be between %v and %v", min, max))
}

// NotBetween succeeds when an ordered value is outside the inclusive range [min, max].
func NotBetween[N cmp.Ordered](min, max N) Fact[N] {
	return NewFact(
		func(x N) bool { return x < min || x > max },
		NewMessage("'%s' must not be between %v and %v", min, max))
}

// GreaterThan succeeds when an ordered value is greater than min.
func GreaterThan[N cmp.Ordered](min N) Fact[N] {
	return NewFact(
		func(x N) bool { return x > min },
		NewMessage("'%s' must be greater than %v", min))
}

// GreaterThanOrEqual succeeds when an ordered value is at least minimum.
func GreaterThanOrEqual[N cmp.Ordered](minimum N) Fact[N] {
	return NewFact(
		func(x N) bool { return x >= minimum },
		NewMessage("'%s' must be greater than or equal to %v", minimum))
}

// LessThan succeeds when an ordered value is less than min.
func LessThan[N cmp.Ordered](min N) Fact[N] {
	return NewFact(
		func(x N) bool { return x < min },
		NewMessage("'%s' must be less than %v", min))
}

// LessThanOrEqual succeeds when an ordered value is at most minimum.
func LessThanOrEqual[N cmp.Ordered](minimum N) Fact[N] {
	return NewFact(
		func(x N) bool { return x <= minimum },
		NewMessage("'%s' must be less than or equal to %v", minimum))
}

// StrEmpty succeeds when a string is empty or contains only whitespace.
func StrEmpty() Fact[string] {
	return NewFact(
		func(x string) bool {
			if x == "" {
				return true
			}

			if strings.TrimSpace(x) == "" {
				return true
			}

			return false
		},
		NewMessage("'%s' must be empty"))
}

// StrNotEmpty succeeds when a string contains a non-whitespace character.
func StrNotEmpty() Fact[string] {
	return NewFact(
		func(x string) bool {
			if x == "" {
				return false
			}

			if strings.TrimSpace(x) == "" {
				return false
			}

			return true
		},
		NewMessage("'%s' must not be empty"))
}

// StrEqualsFold succeeds when a string equals s, ignoring case.
func StrEqualsFold(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return strings.EqualFold(x, s) },
		NewMessage("'%s' must equal '%s'", s))
}

// StrStartsWith succeeds when a string starts with s.
func StrStartsWith(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return strings.HasPrefix(x, s) },
		NewMessage("'%s' must start with '%s'", s))
}

// StrEndsWith succeeds when a string ends with s.
func StrEndsWith(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return strings.HasSuffix(x, s) },
		NewMessage("'%s' must end with '%s'", s))
}

// StrContains succeeds when a string contains s.
func StrContains(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return strings.Contains(x, s) },
		NewMessage("'%s' must contain '%s'", s))
}

// StrNotContains succeeds when a string does not contain s.
func StrNotContains(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return !strings.Contains(x, s) },
		NewMessage("'%s' must not contain '%s'", s))
}

// StrLength succeeds when a string contains exactly n Unicode code points.
func StrLength(n int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) == n },
		NewMessage("'%s' must be %d characters", n))
}

// StrLengthBetween succeeds when a string contains between min and max Unicode code points, inclusive.
func StrLengthBetween(min, max int) Fact[string] {
	return NewFact(
		func(x string) bool {
			l := utf8.RuneCountInString(x)
			return l >= min && l <= max
		},
		NewMessage("'%s' must be between %d and %d characters", min, max))
}

// StrGreaterThan succeeds when a string contains more than min Unicode code points.
func StrGreaterThan(min int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) > min },
		NewMessage("'%s' must be greater than %d characters", min))
}

// StrGreaterThanOrEqual succeeds when a string contains at least min Unicode code points.
func StrGreaterThanOrEqual(min int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) >= min },
		NewMessage("'%s' must be greater than or equal to %d characters", min))
}

// StrLessThan succeeds when a string contains fewer than max Unicode code points.
func StrLessThan(max int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) < max },
		NewMessage("'%s' must be less than %d characters", max))
}

// StrLessThanOrEqual succeeds when a string contains at most max Unicode code points.
func StrLessThanOrEqual(max int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) <= max },
		NewMessage("'%s' must be less than or equal to %d characters", max))
}

// StrMatch succeeds when a string matches r.
func StrMatch(r *regexp.Regexp) Fact[string] {
	return NewFact(r.MatchString,
		NewMessage("'%s' must match the pattern '%s'", r.String()))
}

// StrAlpha succeeds when a string contains only ASCII letters.
func StrAlpha() Fact[string] {
	return NewFact(
		func(x string) bool { return alphaPattern.MatchString(x) },
		NewMessage("'%s' must be alphabetic"))
}

// StrAlphanumeric succeeds when a string contains only ASCII letters and digits.
func StrAlphanumeric() Fact[string] {
	return NewFact(
		func(x string) bool { return alphaNumPattern.MatchString(x) },
		NewMessage("'%s' must be alphanumeric"))
}

// StrIsUrl succeeds when a string is an HTTP or HTTPS URL with a host.
// This is a structural check only; it does not validate that the host is reachable or valid.
func StrIsUrl() Fact[string] {
	return NewFact(
		func(x string) bool {
			if u, e := url.ParseRequestURI(x); e == nil {
				scheme := strings.ToLower(u.Scheme)
				return (scheme == "http" || scheme == "https") && u.Host != ""
			}

			return false
		},
		NewMessage("'%s' is not a valid URL"))
}

// StrIsIPV4 succeeds when a string is an IPv4 address.
func StrIsIPV4() Fact[string] {
	return NewFact(
		func(x string) bool {
			ip := net.ParseIP(x)
			return ip != nil && !strings.Contains(x, ":")
		},
		NewMessage("'%s' is not a valid IPv4 address"))
}

// StrIsIPV6 succeeds when a string is an IPv6 address.
func StrIsIPV6() Fact[string] {
	return NewFact(
		func(x string) bool {
			ip := net.ParseIP(x)
			return ip != nil && strings.Contains(x, ":")
		},
		NewMessage("'%s' is not a valid IPv6 address"))
}

// Empty succeeds when a slice has no elements.
func Empty[T any]() Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) == 0 },
		NewMessage("'%s' must be empty"))
}

// NotEmpty succeeds when a slice has at least one element.
func NotEmpty[T any]() Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) > 0 },
		NewMessage("'%s' must not be empty"))
}

// Length succeeds when a slice contains exactly n elements.
func Length[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) == n },
		NewMessage("'%s' must contain %v items", n))
}

// LengthBetween succeeds when a slice contains between min and max elements, inclusive.
func LengthBetween[T any](min, max int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) >= min && len(items) <= max },
		NewMessage("'%s' must contain between %v and %v items", min, max))
}

// Min succeeds when a slice contains at least n elements.
func Min[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) >= n },
		NewMessage("'%s' must contain at least %v items", n))
}

// Max succeeds when a slice contains no more than n elements.
func Max[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) <= n },
		NewMessage("'%s' must contain no more than %v items", n))
}

// All succeeds when pred returns true for every slice element.
func All[T any](pred func(T) bool) Fact[[]T] {
	return func(items []T) []Message {
		for i, item := range items {
			if !pred(item) {
				return []Message{
					NewMessage("'%s' contains an invalid item at index %v", i+1),
				}
			}
		}
		return nil
	}
}

// Any succeeds when pred returns true for at least one slice element.
func Any[T any](pred func(T) bool) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			return slices.ContainsFunc(items, pred)
		},
		NewMessage("'%s' must contain one valid item"))
}

// UniqueBy succeeds when selector returns a distinct key for every slice element.
func UniqueBy[T any, U comparable](selector func(T) U) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			seen := make(map[U]struct{}, len(items))

			for _, item := range items {
				key := selector(item)
				if _, exists := seen[key]; exists {
					return false
				}
				seen[key] = struct{}{}
			}

			return true
		},
		NewMessage("'%s' must contain unique items"))
}

// Unique succeeds when every element in a slice is distinct.
func Unique[S comparable]() Fact[[]S] {
	return UniqueBy(func(x S) S { return x })
}

// Contains succeeds when a slice contains item.
func Contains[T comparable](item T) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			for _, i := range items {
				if i == item {
					return true
				}
			}
			return false
		},
		NewMessage("'%s' is missing the item: %v", item))
}

// NotContains succeeds when a slice does not contain item.
func NotContains[T comparable](item T) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			for _, i := range items {
				if i == item {
					return false
				}
			}
			return true
		},
		NewMessage("'%s' must not contain the item: %v", item))
}

// MapEmpty succeeds when a map has no entries.
func MapEmpty[K comparable, V any]() Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) == 0 },
		NewMessage("'%s' must be empty"))
}

// MapNotEmpty succeeds when a map has at least one entry.
func MapNotEmpty[K comparable, V any]() Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) > 0 },
		NewMessage("'%s' must not be empty"))
}

// MapLength succeeds when a map has exactly n entries.
func MapLength[K comparable, V any](n int) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) == n },
		NewMessage("'%s' must have %v keys", n))
}

// MapLengthBetween succeeds when a map has between min and max entries, inclusive.
func MapLengthBetween[K comparable, V any](min, max int) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) >= min && len(m) <= max },
		NewMessage("'%s' must have between %v and %v keys", min, max))
}

// MapContains succeeds when a map contains key.
func MapContains[K comparable, V any](key K) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { _, ok := m[key]; return ok },
		NewMessage("'%s' is missing key: %v", key))
}

// MapNotContains succeeds when a map does not contain key.
func MapNotContains[K comparable, V any](key K) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { _, ok := m[key]; return !ok },
		NewMessage("'%s' must not have key: %v", key))
}

// Some succeeds when an Option contains a value.
func Some[T any]() Fact[option.Option[T]] {
	return NewFact[option.Option[T]](
		func(o option.Option[T]) bool { return o.IsSome() },
		NewMessage("'%s' should have a value"))
}

// None succeeds when an Option has no value.
func None[T any]() Fact[option.Option[T]] {
	return NewFact[option.Option[T]](
		func(o option.Option[T]) bool { return o.IsNone() },
		NewMessage("'%s' should have no value"))
}

// Required succeeds when an Option has a value and that value passes every fact.
func Required[T any](facts ...Fact[T]) Fact[option.Option[T]] {
	return func(x option.Option[T]) []Message {
		return option.Match(x,
			func(some T) []Message {
				var m []Message
				for _, fact := range facts {
					for _, message := range fact(some) {
						m = append(m, message)
					}
				}
				return m
			},
			func() []Message { return []Message{NewMessage("'%s' is required")} })
	}
}

// Optional validates an Option's value when present and succeeds when it is None.
func Optional[T any](facts ...Fact[T]) Fact[option.Option[T]] {
	return func(x option.Option[T]) []Message {
		return option.Match(x,
			func(some T) []Message {
				var m []Message
				for _, fact := range facts {
					for _, message := range fact(some) {
						m = append(m, message)
					}
				}
				return m
			},
			func() []Message { return []Message{} })
	}
}

// TimeBetween succeeds when a time is between s and e, inclusive.
func TimeBetween(s, e time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return (x.Before(e) || x.Equal(e)) && (x.After(s) || x.Equal(s)) },
		NewMessage("'%s' must be on or after %s and on or before %s",
			s.Format(time.RFC3339Nano), e.Format(time.RFC3339Nano)))
}

// TimeLessThan succeeds when a time is before t.
func TimeLessThan(t time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.Before(t) },
		NewMessage("'%s' must be before %v", t))
}

// TimeGreaterThan succeeds when a time is after t.
func TimeGreaterThan(t time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.After(t) },
		NewMessage("'%s' must be after %v", t))
}

// Future succeeds when a time is after from.
func Future(from time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.After(from) },
		NewMessage("'%s' must be in the future"))
}

// Past succeeds when a time is before from.
func Past(from time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.Before(from) },
		NewMessage("'%s' must be in the past"))
}

// DurationBetween succeeds when a duration is between s and e, inclusive.
func DurationBetween(s, e time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x >= s && x <= e },
		NewMessage("'%s' must be between %v and %v", s, e))
}

// DurationLessThan succeeds when a duration is less than t.
func DurationLessThan(t time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x < t },
		NewMessage("'%s' must be less than %v", t))
}

// DurationGreaterThan succeeds when a duration is greater than t.
func DurationGreaterThan(t time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x > t },
		NewMessage("'%s' must be greater than %v", t))
}

// DurationPositive succeeds when a duration is greater than zero.
func DurationPositive() Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x > 0 },
		NewMessage("'%s' must be positive"))
}

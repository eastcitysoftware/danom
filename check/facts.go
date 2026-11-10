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
	e164Pattern     = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	emailPattern    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// Combinators
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

// Equality
func Equals[C comparable](e C) Fact[C] {
	return NewFact(
		func(x C) bool { return x == e },
		NewMessage("'%s' must be equal to %v", e))
}

func NotEquals[C comparable](e C) Fact[C] {
	return NewFact(
		func(x C) bool { return x != e },
		NewMessage("'%s' must not be equal to %v", e))
}

// Booleans
func True() Fact[bool] {
	return NewFact(
		func(x bool) bool { return x == true },
		NewMessage("'%s' must be true"))
}

func False() Fact[bool] {
	return NewFact(
		func(x bool) bool { return x != true },
		NewMessage("'%s' must be false"))
}

// Numerics
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func Zero[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x == 0 },
		NewMessage("'%s' must be equal to 0"))
}

func Positive[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x > 0 },
		NewMessage("'%s' must be positive"))
}

func Negative[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x < 0 },
		NewMessage("'%s' must be negative"))
}

func NonZero[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x != 0 },
		NewMessage("'%s' must not be zero"))
}

func NonPositive[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x <= 0 },
		NewMessage("'%s' must be non-positive"))
}

func NonNegative[N Number]() Fact[N] {
	return NewFact(
		func(x N) bool { return x >= 0 },
		NewMessage("'%s' must be non-negative"))
}

// Ordered
func Between[N cmp.Ordered](min, max N) Fact[N] {
	return NewFact(
		func(x N) bool { return x >= min && x <= max },
		NewMessage("'%s' must be between %v and %v", min, max))
}

func NotBetween[N cmp.Ordered](min, max N) Fact[N] {
	return NewFact(
		func(x N) bool { return x < min || x > max },
		NewMessage("'%s' must not be between %v and %v", min, max))
}

func GreaterThan[N cmp.Ordered](min N) Fact[N] {
	return NewFact(
		func(x N) bool { return x > min },
		NewMessage("'%s' must be greater than %v", min))
}

func GreaterThanOrEqual[N cmp.Ordered](minimum N) Fact[N] {
	return NewFact(
		func(x N) bool { return x >= minimum },
		NewMessage("'%s' must be greater than or equal to %v", minimum))
}

func LessThan[N cmp.Ordered](min N) Fact[N] {
	return NewFact(
		func(x N) bool { return x < min },
		NewMessage("'%s' must be less than %v", min))
}

func LessThanOrEqual[N cmp.Ordered](minimum N) Fact[N] {
	return NewFact(
		func(x N) bool { return x <= minimum },
		NewMessage("'%s' must be less than or equal to %v", minimum))
}

// Strings
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

func StrEqualsFold(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return strings.EqualFold(x, s) },
		NewMessage("'%s' must equal '%s'", s))
}

func StrStartsWith(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return strings.HasPrefix(x, s) },
		NewMessage("'%s' must start with '%s'", s))
}

func StrEndsWith(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return strings.HasSuffix(x, s) },
		NewMessage("'%s' must end with '%s'", s))
}

func StrContains(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return strings.Contains(x, s) },
		NewMessage("'%s' must contain '%s'", s))
}

func StrNotContains(s string) Fact[string] {
	return NewFact(
		func(x string) bool { return !strings.Contains(x, s) },
		NewMessage("'%s' must not contain '%s'", s))
}

func StrLength(n int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) == n },
		NewMessage("'%s' must be %d characters", n))
}

func StrLengthBetween(min, max int) Fact[string] {
	return NewFact(
		func(x string) bool {
			l := utf8.RuneCountInString(x)
			return l >= min && l <= max
		},
		NewMessage("'%s' must be between %d and %d characters", min, max))
}

func StrGreaterThan(min int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) > min },
		NewMessage("'%s' must be greater than %d characters", min))
}

func StrGreaterThanOrEqual(min int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) >= min },
		NewMessage("'%s' must be greater than or equal to %d characters", min))
}

func StrLessThan(max int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) < max },
		NewMessage("'%s' must be less than %d characters", max))
}

func StrLessThanOrEqual(max int) Fact[string] {
	return NewFact(
		func(x string) bool { return utf8.RuneCountInString(x) <= max },
		NewMessage("'%s' must be less than or equal to %d characters", max))
}

func StrMatch(r *regexp.Regexp) Fact[string] {
	return NewFact(r.MatchString,
		NewMessage("'%s' must match the pattern '%s'", r.String()))
}

func StrAlpha() Fact[string] {
	return NewFact(
		func(x string) bool { return alphaPattern.MatchString(x) },
		NewMessage("'%s' must be alphabetic"))
}

func StrAlphanumeric() Fact[string] {
	return NewFact(
		func(x string) bool { return alphaNumPattern.MatchString(x) },
		NewMessage("'%s' must be alphanumeric"))
}

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

func StrIsEmail() Fact[string] {
	return NewFact(emailPattern.MatchString,
		NewMessage("'%s' is not a valid email address"))
}

func StrIsE164() Fact[string] {
	return NewFact(e164Pattern.MatchString,
		NewMessage("'%s' is not a valid E.164 phone number"))
}

func StrIsIPV4() Fact[string] {
	return NewFact(
		func(x string) bool {
			ip := net.ParseIP(x)
			return ip != nil && !strings.Contains(x, ":")
		},
		NewMessage("'%s' is not a valid IPv4 address"))
}

func StrIsIPV6() Fact[string] {
	return NewFact(
		func(x string) bool {
			ip := net.ParseIP(x)
			return ip != nil && strings.Contains(x, ":")
		},
		NewMessage("'%s' is not a valid IPv6 address"))
}

// Slice
func Empty[T any]() Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) == 0 },
		NewMessage("'%s' must be empty"))
}

func NotEmpty[T any]() Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) > 0 },
		NewMessage("'%s' must not be empty"))
}

func Length[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) == n },
		NewMessage("'%s' must contain %v items", n))
}

func LengthBetween[T any](min, max int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) >= min && len(items) <= max },
		NewMessage("'%s' must contain between %v and %v items", min, max))
}

func Min[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) >= n },
		NewMessage("'%s' must contain at least %v items", n))
}

func Max[T any](n int) Fact[[]T] {
	return NewFact(
		func(items []T) bool { return len(items) <= n },
		NewMessage("'%s' must contain no more than %v items", n))
}

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

func Any[T any](pred func(T) bool) Fact[[]T] {
	return NewFact(
		func(items []T) bool {
			return slices.ContainsFunc(items, pred)
		},
		NewMessage("'%s' must contain one valid item"))
}

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

func Unique[S comparable]() Fact[[]S] {
	return UniqueBy(func(x S) S { return x })
}

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

// Maps
func MapEmpty[K comparable, V any]() Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) == 0 },
		NewMessage("'%s' must be empty"))
}

func MapNotEmpty[K comparable, V any]() Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) > 0 },
		NewMessage("'%s' must not be empty"))
}

func MapLength[K comparable, V any](n int) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) == n },
		NewMessage("'%s' must have %v keys", n))
}

func MapLengthBetween[K comparable, V any](min, max int) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { return len(m) >= min && len(m) <= max },
		NewMessage("'%s' must have between %v and %v keys", min, max))
}

func MapContains[K comparable, V any](key K) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { _, ok := m[key]; return ok },
		NewMessage("'%s' is missing key: %v", key))
}

func MapNotContains[K comparable, V any](key K) Fact[map[K]V] {
	return NewFact(
		func(m map[K]V) bool { _, ok := m[key]; return !ok },
		NewMessage("'%s' must not have key: %v", key))
}

// Options
func Some[T any]() Fact[option.Option[T]] {
	return NewFact[option.Option[T]](
		option.IsSome,
		NewMessage("'%s' should have a value"))
}

func None[T any]() Fact[option.Option[T]] {
	return NewFact[option.Option[T]](
		option.IsNone,
		NewMessage("'%s' should have no value"))
}

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

// Time
func TimeBetween(s, e time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return (x.Before(e) || x.Equal(e)) && (x.After(s) || x.Equal(s)) },
		NewMessage("'%s' must be on or after %s and on or before %s",
			s.Format(time.RFC3339Nano), e.Format(time.RFC3339Nano)))
}

func TimeLessThan(t time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.Before(t) },
		NewMessage("'%s' must be before %v", t))
}

func TimeGreaterThan(t time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.After(t) },
		NewMessage("'%s' must be after %v", t))
}

func Future(from time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.After(from) },
		NewMessage("'%s' must be in the future"))
}

func Past(from time.Time) Fact[time.Time] {
	return NewFact(
		func(x time.Time) bool { return x.Before(from) },
		NewMessage("'%s' must be in the past"))
}

func DurationBetween(s, e time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x >= s && x <= e },
		NewMessage("'%s' must be between %v and %v", s, e))
}

func DurationLessThan(t time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x < t },
		NewMessage("'%s' must be less than %v", t))
}

func DurationGreaterThan(t time.Duration) Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x > t },
		NewMessage("'%s' must be greater than %v", t))
}

func DurationPositive() Fact[time.Duration] {
	return NewFact(
		func(x time.Duration) bool { return x > 0 },
		NewMessage("'%s' must be positive"))
}

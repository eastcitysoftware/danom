package check

import (
	"net"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	alphaPattern    = regexp.MustCompile("^[a-zA-Z]+$")
	alphaNumPattern = regexp.MustCompile("^[a-zA-Z0-9]+$")
)

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

package check

import (
	"regexp"
	"testing"
	"time"

	"github.com/eastcitysoftware/danom/option"
)

func TestWhen(t *testing.T) {
	fact := When(StrStartsWith("a"), StrLengthBetween(3, 5), StrNotContains("z"))

	runFact(t, fact, []factCase[string]{
		{"condition not met skips facts", "b", true, ""},
		{"condition met and facts pass", "abc", true, ""},
		{"condition met and one fact fails", "a", false, "'f' must be between 3 and 5 characters"},
		{"condition met and both facts fail", "az", false, ""},
	})

	if got := len(fact("az")); got != 2 {
		t.Errorf("expected 2 messages, got %d", got)
	}
}

func TestNoneOf(t *testing.T) {
	fact := NoneOf(StrEqualsFold("here"), Equals("there"))

	runFact(t, fact, []factCase[string]{
		{"matches none", "elsewhere", true, ""},
		{"matches first", "HERE", false, "'f' must not satisfy any of the conditions"},
		{"matches second", "there", false, ""},
	})

	runFact(t, NoneOf[string](), []factCase[string]{
		{"no facts", "anything", true, ""},
	})
}

func TestAnyOf(t *testing.T) {
	fact := AnyOf(Equals(1), Equals(2))

	runFact(t, fact, []factCase[int]{
		{"matches first", 1, true, ""},
		{"matches second", 2, true, ""},
		{"matches none", 3, false, "'f' must satisfy at least one condition"},
	})

	runFact(t, AnyOf[int](), []factCase[int]{
		{"no facts", 1, false, ""},
	})
}

func TestEquals(t *testing.T) {
	runFact(t, Equals("a"), []factCase[string]{
		{"equal", "a", true, ""},
		{"not equal", "b", false, "'f' must be equal to a"},
	})

	runFact(t, Equals("100%"), []factCase[string]{
		{"percent in argument", "x", false, "'f' must be equal to 100%"},
	})
}

func TestNotEquals(t *testing.T) {
	runFact(t, NotEquals(1), []factCase[int]{
		{"not equal", 2, true, ""},
		{"equal", 1, false, "'f' must not be equal to 1"},
	})
}

func TestTrue(t *testing.T) {
	runFact(t, True(), []factCase[bool]{
		{"true", true, true, ""},
		{"false", false, false, "'f' must be true"},
	})
}

func TestFalse(t *testing.T) {
	runFact(t, False(), []factCase[bool]{
		{"false", false, true, ""},
		{"true", true, false, "'f' must be false"},
	})
}

func TestZero(t *testing.T) {
	runFact(t, Zero[int](), []factCase[int]{
		{"zero", 0, true, ""},
		{"positive", 1, false, "'f' must be equal to 0"},
		{"negative", -1, false, ""},
	})
}

func TestPositive(t *testing.T) {
	runFact(t, Positive[float64](), []factCase[float64]{
		{"positive", 0.1, true, ""},
		{"zero", 0, false, "'f' must be positive"},
		{"negative", -1, false, ""},
	})
}

func TestNegative(t *testing.T) {
	runFact(t, Negative[int](), []factCase[int]{
		{"negative", -1, true, ""},
		{"zero", 0, false, "'f' must be negative"},
		{"positive", 1, false, ""},
	})
}

func TestNonZero(t *testing.T) {
	runFact(t, NonZero[int](), []factCase[int]{
		{"positive", 1, true, ""},
		{"negative", -1, true, ""},
		{"zero", 0, false, "'f' must not be zero"},
	})
}

func TestNonPositive(t *testing.T) {
	runFact(t, NonPositive[int](), []factCase[int]{
		{"negative", -1, true, ""},
		{"zero", 0, true, ""},
		{"positive", 1, false, "'f' must be non-positive"},
	})
}

func TestNonNegative(t *testing.T) {
	runFact(t, NonNegative[int](), []factCase[int]{
		{"positive", 1, true, ""},
		{"zero", 0, true, ""},
		{"negative", -1, false, "'f' must be non-negative"},
	})
}

func TestBetween(t *testing.T) {
	runFact(t, Between(1, 3), []factCase[int]{
		{"below", 0, false, "'f' must be between 1 and 3"},
		{"lower bound", 1, true, ""},
		{"inside", 2, true, ""},
		{"upper bound", 3, true, ""},
		{"above", 4, false, ""},
	})

	runFact(t, Between("b", "d"), []factCase[string]{
		{"string inside", "c", true, ""},
		{"string outside", "e", false, ""},
	})
}

func TestNotBetween(t *testing.T) {
	runFact(t, NotBetween(1, 3), []factCase[int]{
		{"below", 0, true, ""},
		{"lower bound", 1, false, "'f' must not be between 1 and 3"},
		{"inside", 2, false, ""},
		{"upper bound", 3, false, ""},
		{"above", 4, true, ""},
	})
}

func TestGreaterThan(t *testing.T) {
	runFact(t, GreaterThan(5), []factCase[int]{
		{"greater", 6, true, ""},
		{"equal", 5, false, "'f' must be greater than 5"},
		{"less", 4, false, ""},
	})
}

func TestGreaterThanOrEqual(t *testing.T) {
	runFact(t, GreaterThanOrEqual(5), []factCase[int]{
		{"greater", 6, true, ""},
		{"equal", 5, true, ""},
		{"less", 4, false, "'f' must be greater than or equal to 5"},
	})
}

func TestLessThan(t *testing.T) {
	runFact(t, LessThan(5), []factCase[int]{
		{"less", 4, true, ""},
		{"equal", 5, false, "'f' must be less than 5"},
		{"greater", 6, false, ""},
	})
}

func TestLessThanOrEqual(t *testing.T) {
	runFact(t, LessThanOrEqual(5), []factCase[int]{
		{"less", 4, true, ""},
		{"equal", 5, true, ""},
		{"greater", 6, false, "'f' must be less than or equal to 5"},
	})
}

func TestStrEmpty(t *testing.T) {
	runFact(t, StrEmpty(), []factCase[string]{
		{"empty", "", true, ""},
		{"whitespace", "  \t", true, ""},
		{"text", "a", false, "'f' must be empty"},
	})
}

func TestStrNotEmpty(t *testing.T) {
	runFact(t, StrNotEmpty(), []factCase[string]{
		{"text", "a", true, ""},
		{"empty", "", false, "'f' must not be empty"},
		{"whitespace", "  ", false, ""},
	})
}

func TestStrEqualsFold(t *testing.T) {
	runFact(t, StrEqualsFold("Danom"), []factCase[string]{
		{"same case", "Danom", true, ""},
		{"different case", "dANOM", true, ""},
		{"different", "other", false, "'f' must equal 'Danom'"},
	})
}

func TestStrStartsWith(t *testing.T) {
	runFact(t, StrStartsWith("da"), []factCase[string]{
		{"prefix", "danom", true, ""},
		{"no prefix", "nom", false, "'f' must start with 'da'"},
	})
}

func TestStrEndsWith(t *testing.T) {
	runFact(t, StrEndsWith("om"), []factCase[string]{
		{"suffix", "danom", true, ""},
		{"no suffix", "dan", false, "'f' must end with 'om'"},
	})
}

func TestStrContains(t *testing.T) {
	runFact(t, StrContains("an"), []factCase[string]{
		{"contains", "danom", true, ""},
		{"missing", "dom", false, "'f' must contain 'an'"},
	})

	runFact(t, StrContains("100%"), []factCase[string]{
		{"percent in argument", "x", false, "'f' must contain '100%'"},
	})
}

func TestStrNotContains(t *testing.T) {
	runFact(t, StrNotContains("an"), []factCase[string]{
		{"missing", "dom", true, ""},
		{"contains", "danom", false, "'f' must not contain 'an'"},
	})
}

func TestStrLength(t *testing.T) {
	runFact(t, StrLength(3), []factCase[string]{
		{"exact", "abc", true, ""},
		{"multi-byte runes", "日本語", true, ""},
		{"short", "ab", false, "'f' must be 3 characters"},
		{"long", "abcd", false, ""},
	})
}

func TestStrLengthBetween(t *testing.T) {
	runFact(t, StrLengthBetween(2, 4), []factCase[string]{
		{"short", "a", false, "'f' must be between 2 and 4 characters"},
		{"lower bound", "ab", true, ""},
		{"upper bound", "abcd", true, ""},
		{"long", "abcde", false, ""},
		{"multi-byte runes", "日本", true, ""},
	})
}

func TestStrGreaterThan(t *testing.T) {
	runFact(t, StrGreaterThan(2), []factCase[string]{
		{"longer", "abc", true, ""},
		{"equal", "ab", false, "'f' must be greater than 2 characters"},
		{"shorter", "a", false, ""},
	})
}

func TestStrGreaterThanOrEqual(t *testing.T) {
	runFact(t, StrGreaterThanOrEqual(2), []factCase[string]{
		{"longer", "abc", true, ""},
		{"equal", "ab", true, ""},
		{"shorter", "a", false, "'f' must be greater than or equal to 2 characters"},
	})
}

func TestStrLessThan(t *testing.T) {
	runFact(t, StrLessThan(2), []factCase[string]{
		{"shorter", "a", true, ""},
		{"equal", "ab", false, "'f' must be less than 2 characters"},
		{"longer", "abc", false, ""},
	})
}

func TestStrLessThanOrEqual(t *testing.T) {
	runFact(t, StrLessThanOrEqual(2), []factCase[string]{
		{"shorter", "a", true, ""},
		{"equal", "ab", true, ""},
		{"longer", "abc", false, "'f' must be less than or equal to 2 characters"},
	})
}

func TestStrMatch(t *testing.T) {
	runFact(t, StrMatch(regexp.MustCompile(`^\d+$`)), []factCase[string]{
		{"match", "123", true, ""},
		{"no match", "12a", false, `'f' must match the pattern '^\d+$'`},
	})
}

func TestStrAlpha(t *testing.T) {
	runFact(t, StrAlpha(), []factCase[string]{
		{"letters", "abcXYZ", true, ""},
		{"digits", "abc1", false, "'f' must be alphabetic"},
		{"empty", "", false, ""},
	})
}

func TestStrAlphanumeric(t *testing.T) {
	runFact(t, StrAlphanumeric(), []factCase[string]{
		{"letters and digits", "abc123", true, ""},
		{"symbol", "abc-1", false, "'f' must be alphanumeric"},
		{"empty", "", false, ""},
	})
}

func TestStrIsUrl(t *testing.T) {
	runFact(t, StrIsUrl(), []factCase[string]{
		{"http", "http://example.com", true, ""},
		{"https with path", "https://example.com/a?b=c", true, ""},
		{"uppercase scheme", "HTTPS://example.com", true, ""},
		{"empty", "", false, "'f' is not a valid URL"},
		{"no scheme", "example.com", false, ""},
		{"relative path", "/a/b", false, ""},
		{"unsupported scheme", "ftp://example.com", false, ""},
		{"no host", "http://", false, ""},
	})
}

func TestStrIsIPV4(t *testing.T) {
	runFact(t, StrIsIPV4(), []factCase[string]{
		{"valid", "192.168.0.1", true, ""},
		{"out of range", "256.1.1.1", false, "'f' is not a valid IPv4 address"},
		{"ipv6", "2001:db8::1", false, ""},
		{"empty", "", false, ""},
	})
}

func TestStrIsIPV6(t *testing.T) {
	runFact(t, StrIsIPV6(), []factCase[string]{
		{"full", "2001:0db8:0000:0000:0000:0000:0000:0001", true, ""},
		{"compressed", "2001:db8::1", true, ""},
		{"loopback", "::1", true, ""},
		{"ipv4", "192.168.0.1", false, "'f' is not a valid IPv6 address"},
		{"invalid", "2001:db8:::1", false, ""},
	})
}

func TestEmpty(t *testing.T) {
	runFact(t, Empty[int](), []factCase[[]int]{
		{"nil", nil, true, ""},
		{"empty", []int{}, true, ""},
		{"items", []int{1}, false, "'f' must be empty"},
	})
}

func TestNotEmpty(t *testing.T) {
	runFact(t, NotEmpty[int](), []factCase[[]int]{
		{"items", []int{1}, true, ""},
		{"empty", []int{}, false, "'f' must not be empty"},
		{"nil", nil, false, ""},
	})
}

func TestLength(t *testing.T) {
	runFact(t, Length[int](2), []factCase[[]int]{
		{"exact", []int{1, 2}, true, ""},
		{"short", []int{1}, false, "'f' must contain 2 items"},
		{"long", []int{1, 2, 3}, false, ""},
	})
}

func TestLengthBetween(t *testing.T) {
	runFact(t, LengthBetween[int](1, 2), []factCase[[]int]{
		{"empty", nil, false, "'f' must contain between 1 and 2 items"},
		{"lower bound", []int{1}, true, ""},
		{"upper bound", []int{1, 2}, true, ""},
		{"long", []int{1, 2, 3}, false, ""},
	})
}

func TestMin(t *testing.T) {
	runFact(t, Min[int](2), []factCase[[]int]{
		{"more", []int{1, 2, 3}, true, ""},
		{"equal", []int{1, 2}, true, ""},
		{"fewer", []int{1}, false, "'f' must contain at least 2 items"},
	})
}

func TestMax(t *testing.T) {
	runFact(t, Max[int](2), []factCase[[]int]{
		{"fewer", []int{1}, true, ""},
		{"nil", nil, true, ""},
		{"equal", []int{1, 2}, true, ""},
		{"more", []int{1, 2, 3}, false, "'f' must contain no more than 2 items"},
	})
}

func TestAll(t *testing.T) {
	runFact(t, All(func(x int) bool { return x > 0 }), []factCase[[]int]{
		{"all valid", []int{1, 2}, true, ""},
		{"empty", nil, true, ""},
		{"first invalid reports position", []int{1, 0, -1}, false, "'f' contains an invalid item at index 2"},
	})
}

func TestAny(t *testing.T) {
	runFact(t, Any(func(x int) bool { return x > 1 }), []factCase[[]int]{
		{"one valid", []int{0, 2}, true, ""},
		{"none valid", []int{0, 1}, false, "'f' must contain one valid item"},
		{"empty", nil, false, ""},
	})
}

func TestUniqueBy(t *testing.T) {
	runFact(t, UniqueBy(func(p mockPerson) string { return p.name }), []factCase[[]mockPerson]{
		{"unique keys", []mockPerson{{"a", 1}, {"b", 1}}, true, ""},
		{"duplicate keys", []mockPerson{{"a", 1}, {"a", 2}}, false, "'f' must contain unique items"},
		{"empty", nil, true, ""},
	})
}

func TestUnique(t *testing.T) {
	runFact(t, Unique[int](), []factCase[[]int]{
		{"unique", []int{1, 2, 3}, true, ""},
		{"duplicate", []int{1, 2, 1}, false, "'f' must contain unique items"},
		{"empty", nil, true, ""},
	})
}

func TestContains(t *testing.T) {
	runFact(t, Contains("b"), []factCase[[]string]{
		{"present", []string{"a", "b"}, true, ""},
		{"missing", []string{"a"}, false, "'f' is missing the item: b"},
		{"empty", nil, false, ""},
	})
}

func TestNotContains(t *testing.T) {
	runFact(t, NotContains("b"), []factCase[[]string]{
		{"missing", []string{"a"}, true, ""},
		{"empty", nil, true, ""},
		{"present", []string{"a", "b"}, false, "'f' must not contain the item: b"},
	})
}

func TestMapEmpty(t *testing.T) {
	runFact(t, MapEmpty[string, int](), []factCase[map[string]int]{
		{"nil", nil, true, ""},
		{"empty", map[string]int{}, true, ""},
		{"items", map[string]int{"a": 1}, false, "'f' must be empty"},
	})
}

func TestMapNotEmpty(t *testing.T) {
	runFact(t, MapNotEmpty[string, int](), []factCase[map[string]int]{
		{"items", map[string]int{"a": 1}, true, ""},
		{"empty", map[string]int{}, false, "'f' must not be empty"},
		{"nil", nil, false, ""},
	})
}

func TestMapLength(t *testing.T) {
	runFact(t, MapLength[string, int](1), []factCase[map[string]int]{
		{"exact", map[string]int{"a": 1}, true, ""},
		{"empty", nil, false, "'f' must have 1 keys"},
		{"more", map[string]int{"a": 1, "b": 2}, false, ""},
	})
}

func TestMapLengthBetween(t *testing.T) {
	runFact(t, MapLengthBetween[string, int](1, 2), []factCase[map[string]int]{
		{"empty", nil, false, "'f' must have between 1 and 2 keys"},
		{"lower bound", map[string]int{"a": 1}, true, ""},
		{"upper bound", map[string]int{"a": 1, "b": 2}, true, ""},
		{"more", map[string]int{"a": 1, "b": 2, "c": 3}, false, ""},
	})
}

func TestMapContains(t *testing.T) {
	runFact(t, MapContains[string, int]("a"), []factCase[map[string]int]{
		{"present", map[string]int{"a": 1}, true, ""},
		{"missing", map[string]int{"b": 1}, false, "'f' is missing key: a"},
		{"nil", nil, false, ""},
	})
}

func TestMapNotContains(t *testing.T) {
	runFact(t, MapNotContains[string, int]("a"), []factCase[map[string]int]{
		{"missing", map[string]int{"b": 1}, true, ""},
		{"nil", nil, true, ""},
		{"present", map[string]int{"a": 1}, false, "'f' must not have key: a"},
	})
}

func TestMapAll(t *testing.T) {
	runFact(t, MapAll[string, int](func(value int) bool { return value > 0 }), []factCase[map[string]int]{
		{"all values pass", map[string]int{"a": 1, "b": 2}, true, ""},
		{"one value fails", map[string]int{"a": 1, "b": 0}, false, "'f' must contain only valid values"},
		{"empty", map[string]int{}, true, ""},
		{"nil", nil, true, ""},
	})
}

func TestMapAny(t *testing.T) {
	runFact(t, MapAny[string, int](func(value int) bool { return value > 0 }), []factCase[map[string]int]{
		{"one value passes", map[string]int{"a": -1, "b": 2}, true, ""},
		{"no values pass", map[string]int{"a": -1, "b": 0}, false, "'f' must contain at least one valid value"},
		{"empty", map[string]int{}, false, "'f' must contain at least one valid value"},
		{"nil", nil, false, ""},
	})
}

func TestSome(t *testing.T) {
	runFact(t, Some[int](), []factCase[option.Option[int]]{
		{"some", option.Some(1), true, ""},
		{"none", option.None[int](), false, "'f' should have a value"},
	})
}

func TestNone(t *testing.T) {
	runFact(t, None[int](), []factCase[option.Option[int]]{
		{"none", option.None[int](), true, ""},
		{"some", option.Some(1), false, "'f' should have no value"},
	})
}

func TestRequired(t *testing.T) {
	runFact(t, Required(GreaterThan(0), LessThan(10)), []factCase[option.Option[int]]{
		{"some and valid", option.Some(5), true, ""},
		{"none", option.None[int](), false, "'f' is required"},
		{"some and invalid", option.Some(-1), false, "'f' must be greater than 0"},
	})

	runFact(t, Required[int](), []factCase[option.Option[int]]{
		{"some without facts", option.Some(1), true, ""},
	})

	if got := len(Required(GreaterThan(0), GreaterThan(1))(option.Some(0))); got != 2 {
		t.Errorf("expected 2 messages, got %d", got)
	}
}

func TestOptional(t *testing.T) {
	runFact(t, Optional(GreaterThan(0), LessThan(10)), []factCase[option.Option[int]]{
		{"none", option.None[int](), true, ""},
		{"some and valid", option.Some(5), true, ""},
		{"some and invalid", option.Some(-1), false, "'f' must be greater than 0"},
	})

	if got := len(Optional(GreaterThan(0), GreaterThan(1))(option.Some(0))); got != 2 {
		t.Errorf("expected 2 messages, got %d", got)
	}
}

func TestTimeBetween(t *testing.T) {
	s := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	runFact(t, TimeBetween(s, e), []factCase[time.Time]{
		{"before", s.Add(-time.Nanosecond), false,
			"'f' must be on or after 2024-01-01T00:00:00Z and on or before 2024-01-31T00:00:00Z"},
		{"start boundary", s, true, ""},
		{"inside", s.Add(24 * time.Hour), true, ""},
		{"end boundary", e, true, ""},
		{"after", e.Add(time.Nanosecond), false, ""},
	})
}

func TestTimeLessThan(t *testing.T) {
	b := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	runFact(t, TimeLessThan(b), []factCase[time.Time]{
		{"before", b.Add(-time.Second), true, ""},
		{"equal", b, false, "'f' must be before " + b.String()},
		{"after", b.Add(time.Second), false, ""},
	})
}

func TestTimeGreaterThan(t *testing.T) {
	b := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	runFact(t, TimeGreaterThan(b), []factCase[time.Time]{
		{"after", b.Add(time.Second), true, ""},
		{"equal", b, false, "'f' must be after " + b.String()},
		{"before", b.Add(-time.Second), false, ""},
	})
}

func TestFuture(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	runFact(t, Future(from), []factCase[time.Time]{
		{"after", from.Add(time.Second), true, ""},
		{"equal", from, false, "'f' must be in the future"},
		{"before", from.Add(-time.Second), false, ""},
	})
}

func TestPast(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	runFact(t, Past(from), []factCase[time.Time]{
		{"before", from.Add(-time.Second), true, ""},
		{"equal", from, false, "'f' must be in the past"},
		{"after", from.Add(time.Second), false, ""},
	})
}

func TestDurationBetween(t *testing.T) {
	runFact(t, DurationBetween(time.Second, time.Minute), []factCase[time.Duration]{
		{"below", time.Millisecond, false, "'f' must be between 1s and 1m0s"},
		{"lower bound", time.Second, true, ""},
		{"inside", 30 * time.Second, true, ""},
		{"upper bound", time.Minute, true, ""},
		{"above", 2 * time.Minute, false, ""},
	})
}

func TestDurationLessThan(t *testing.T) {
	runFact(t, DurationLessThan(time.Second), []factCase[time.Duration]{
		{"less", time.Millisecond, true, ""},
		{"equal", time.Second, false, "'f' must be less than 1s"},
		{"greater", time.Minute, false, ""},
	})
}

func TestDurationGreaterThan(t *testing.T) {
	runFact(t, DurationGreaterThan(time.Second), []factCase[time.Duration]{
		{"greater", time.Minute, true, ""},
		{"equal", time.Second, false, "'f' must be greater than 1s"},
		{"less", time.Millisecond, false, ""},
	})
}

func TestDurationPositive(t *testing.T) {
	runFact(t, DurationPositive(), []factCase[time.Duration]{
		{"positive", time.Nanosecond, true, ""},
		{"zero", 0, false, "'f' must be positive"},
		{"negative", -time.Second, false, ""},
	})
}

// factCase describes one input for a fact. When valid is false and msg is
// not empty, msg is compared to the first message produced for the field "f".
type factCase[T any] struct {
	name  string
	input T
	valid bool
	msg   string
}

type mockPerson struct {
	name string
	age  int
}

func runFact[T any](t *testing.T, fact Fact[T], cases []factCase[T]) {
	t.Helper()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			messages := fact(c.input)

			if c.valid {
				if len(messages) > 0 {
					t.Fatalf("expected input %v to be valid, got message '%s'", c.input, messages[0]("f"))
				}
				return
			}

			if len(messages) == 0 {
				t.Fatalf("expected input %v to be invalid, got no messages", c.input)
			}

			if got := messages[0]("f"); c.msg != "" && got != c.msg {
				t.Errorf("unexpected message, got '%s' expected '%s'", got, c.msg)
			}
		})
	}
}

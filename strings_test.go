package spvalidator

import "testing"

func TestStringValidators(t *testing.T) {
	cases := []struct {
		name string
		tag  string
		ok   string
		bad  string
		fn   func(string) error
	}{
		{"alpha", "alpha", "Go", "Go1", Alpha},
		{"alpha_space", "alphaspace", "Go Team", "Go1", AlphaSpace},
		{"alphanum", "alphanum", "Go123", "Go-123", Alphanum},
		{"alphanum_space", "alphanumspace", "Go 123", "Go!", AlphanumSpace},
		{"alphanum_unicode", "alphanumunicode", "東京123", "Go!", AlphanumUnicode},
		{"alpha_unicode", "alphaunicode", "東京", "東京1", AlphaUnicode},
		{"ascii", "ascii", "hello", "héllo", ASCII},
		{"boolean", "boolean", "true", "maybe", Boolean},
		{"lowercase", "lowercase", "hello", "Hello", Lowercase},
		{"multibyte", "multibyte", "hello✓", "hello", Multibyte},
		{"number", "number", "123", "12.3", Number},
		{"numeric", "numeric", "12.3", " 12.3 ", Numeric},
		{"printascii", "printascii", "Hello!", "Hello\n", PrintASCII},
		{"uppercase", "uppercase", "HELLO", "Hello", Uppercase},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expectNoError(t, tc.fn(tc.ok))
			validateErr(t, tc.fn(tc.bad), tc.tag)
		})
	}
}

func TestStringPairValidators(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		expectNoError(t, Contains("hello", "ell"))
		validateErr(t, Contains("hello", "xyz"), "contains")
	})
	t.Run("contains_any", func(t *testing.T) {
		expectNoError(t, ContainsAny("hello", "xyzol"))
		validateErr(t, ContainsAny("hello", "xyz"), "containsany")
	})
	t.Run("contains_rune", func(t *testing.T) {
		expectNoError(t, ContainsRune("hello", 'e'))
		validateErr(t, ContainsRune("hello", 'z'), "containsrune")
	})
	t.Run("ends_with", func(t *testing.T) {
		expectNoError(t, EndsWith("hello", "lo"))
		validateErr(t, EndsWith("hello", "he"), "endswith")
	})
	t.Run("ends_not_with", func(t *testing.T) {
		expectNoError(t, EndsNotWith("hello", "he"))
		validateErr(t, EndsNotWith("hello", "lo"), "endsnotwith")
	})
	t.Run("excludes", func(t *testing.T) {
		expectNoError(t, Excludes("hello", "xyz"))
		validateErr(t, Excludes("hello", "ell"), "excludes")
	})
	t.Run("excludes_all", func(t *testing.T) {
		expectNoError(t, ExcludesAll("hello", "xyz"))
		validateErr(t, ExcludesAll("hello", "el"), "excludesall")
	})
	t.Run("excludes_rune", func(t *testing.T) {
		expectNoError(t, ExcludesRune("hello", 'z'))
		validateErr(t, ExcludesRune("hello", 'e'), "excludesrune")
	})
	t.Run("starts_with", func(t *testing.T) {
		expectNoError(t, StartsWith("hello", "he"))
		validateErr(t, StartsWith("hello", "lo"), "startswith")
	})
	t.Run("starts_not_with", func(t *testing.T) {
		expectNoError(t, StartsNotWith("hello", "lo"))
		validateErr(t, StartsNotWith("hello", "he"), "startsnotwith")
	})
}

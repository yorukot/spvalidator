package spvalidator

import (
	"strconv"
	"strings"
	"unicode"
)

func Alpha(value string) error {
	if value == "" {
		return fail("alpha", value, nil, "must contain ASCII letters only")
	}
	for _, r := range value {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return fail("alpha", value, nil, "must contain ASCII letters only")
		}
	}
	return nil
}

func AlphaSpace(value string) error {
	if value == "" {
		return fail("alphaspace", value, nil, "must contain ASCII letters and spaces only")
	}
	for _, r := range value {
		if r != ' ' && ((r < 'A' || r > 'Z') && (r < 'a' || r > 'z')) {
			return fail("alphaspace", value, nil, "must contain ASCII letters and spaces only")
		}
	}
	return nil
}

func Alphanum(value string) error {
	if value == "" {
		return fail("alphanum", value, nil, "must contain ASCII letters and numbers only")
	}
	for _, r := range value {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return fail("alphanum", value, nil, "must contain ASCII letters and numbers only")
		}
	}
	return nil
}

func AlphanumSpace(value string) error {
	if value == "" {
		return fail("alphanumspace", value, nil, "must contain ASCII letters, numbers, and spaces only")
	}
	for _, r := range value {
		if r != ' ' && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return fail("alphanumspace", value, nil, "must contain ASCII letters, numbers, and spaces only")
		}
	}
	return nil
}

func AlphanumUnicode(value string) error {
	if value == "" {
		return fail("alphanumunicode", value, nil, "must contain letters and numbers only")
	}
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return fail("alphanumunicode", value, nil, "must contain letters and numbers only")
		}
	}
	return nil
}

func AlphaUnicode(value string) error {
	if value == "" {
		return fail("alphaunicode", value, nil, "must contain letters only")
	}
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return fail("alphaunicode", value, nil, "must contain letters only")
		}
	}
	return nil
}

func ASCII(value string) error {
	for _, r := range value {
		if r > 127 {
			return fail("ascii", value, nil, "must contain ASCII characters only")
		}
	}
	return nil
}

func Boolean(value string) error {
	if _, err := strconv.ParseBool(value); err == nil {
		return nil
	}
	switch strings.ToLower(value) {
	case "yes", "no", "y", "n", "0", "1":
		return nil
	default:
		return fail("boolean", value, nil, "must be a boolean string")
	}
}

func Contains(value string, substr string) error {
	if strings.Contains(value, substr) {
		return nil
	}
	return failf("contains", value, substr, "must contain %q", substr)
}

func ContainsAny(value string, chars string) error {
	if strings.ContainsAny(value, chars) {
		return nil
	}
	return failf("containsany", value, chars, "must contain any of %q", chars)
}

func ContainsRune(value string, r rune) error {
	if strings.ContainsRune(value, r) {
		return nil
	}
	return failf("containsrune", value, string(r), "must contain %q", string(r))
}

func EndsNotWith(value string, suffix string) error {
	if !strings.HasSuffix(value, suffix) {
		return nil
	}
	return failf("endsnotwith", value, suffix, "must not end with %q", suffix)
}

func EndsWith(value string, suffix string) error {
	if strings.HasSuffix(value, suffix) {
		return nil
	}
	return failf("endswith", value, suffix, "must end with %q", suffix)
}

func Excludes(value string, substr string) error {
	if !strings.Contains(value, substr) {
		return nil
	}
	return failf("excludes", value, substr, "must not contain %q", substr)
}

func ExcludesAll(value string, chars string) error {
	if !strings.ContainsAny(value, chars) {
		return nil
	}
	return failf("excludesall", value, chars, "must not contain any of %q", chars)
}

func ExcludesRune(value string, r rune) error {
	if !strings.ContainsRune(value, r) {
		return nil
	}
	return failf("excludesrune", value, string(r), "must not contain %q", string(r))
}

func Lowercase(value string) error {
	if value == strings.ToLower(value) {
		return nil
	}
	return fail("lowercase", value, nil, "must be lowercase")
}

func Multibyte(value string) error {
	for _, r := range value {
		if r > 127 {
			return nil
		}
	}
	return fail("multibyte", value, nil, "must contain at least one multibyte character")
}

func Number(value string) error {
	if value == "" {
		return fail("number", value, nil, "must be a number")
	}
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return nil
	}
	if _, err := strconv.ParseUint(value, 10, 64); err == nil {
		return nil
	}
	return fail("number", value, nil, "must be an integer number")
}

func Numeric(value string) error {
	if _, err := strconv.ParseFloat(value, 64); err == nil && strings.TrimSpace(value) == value && value != "" {
		return nil
	}
	return fail("numeric", value, nil, "must be numeric")
}

func PrintASCII(value string) error {
	for _, r := range value {
		if r < 32 || r > 126 {
			return fail("printascii", value, nil, "must contain printable ASCII characters only")
		}
	}
	return nil
}

func StartsNotWith(value string, prefix string) error {
	if !strings.HasPrefix(value, prefix) {
		return nil
	}
	return failf("startsnotwith", value, prefix, "must not start with %q", prefix)
}

func StartsWith(value string, prefix string) error {
	if strings.HasPrefix(value, prefix) {
		return nil
	}
	return failf("startswith", value, prefix, "must start with %q", prefix)
}

func Uppercase(value string) error {
	if value == strings.ToUpper(value) {
		return nil
	}
	return fail("uppercase", value, nil, "must be uppercase")
}

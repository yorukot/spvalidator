package spvalidator

import (
	"strconv"
	"strings"
	"unicode"
)

func Alpha(value string) error {
	if value == "" {
		return fail("alpha", value, nil, "value must contain ASCII letters only")
	}
	for _, r := range value {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return fail("alpha", value, nil, "value must contain ASCII letters only")
		}
	}
	return nil
}

func AlphaSpace(value string) error {
	if value == "" {
		return fail("alphaspace", value, nil, "value must contain ASCII letters and spaces only")
	}
	for _, r := range value {
		if r != ' ' && ((r < 'A' || r > 'Z') && (r < 'a' || r > 'z')) {
			return fail("alphaspace", value, nil, "value must contain ASCII letters and spaces only")
		}
	}
	return nil
}

func Alphanum(value string) error {
	if value == "" {
		return fail("alphanum", value, nil, "value must contain ASCII letters and numbers only")
	}
	for _, r := range value {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return fail("alphanum", value, nil, "value must contain ASCII letters and numbers only")
		}
	}
	return nil
}

func AlphanumSpace(value string) error {
	if value == "" {
		return fail("alphanumspace", value, nil, "value must contain ASCII letters, numbers, and spaces only")
	}
	for _, r := range value {
		if r != ' ' && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return fail("alphanumspace", value, nil, "value must contain ASCII letters, numbers, and spaces only")
		}
	}
	return nil
}

func AlphanumUnicode(value string) error {
	if value == "" {
		return fail("alphanumunicode", value, nil, "value must contain letters and numbers only")
	}
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return fail("alphanumunicode", value, nil, "value must contain letters and numbers only")
		}
	}
	return nil
}

func AlphaUnicode(value string) error {
	if value == "" {
		return fail("alphaunicode", value, nil, "value must contain letters only")
	}
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return fail("alphaunicode", value, nil, "value must contain letters only")
		}
	}
	return nil
}

func ASCII(value string) error {
	for _, r := range value {
		if r > 127 {
			return fail("ascii", value, nil, "value must contain ASCII characters only")
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
		return fail("boolean", value, nil, "value must be a boolean string")
	}
}

func Contains(value string, substr string) error {
	if strings.Contains(value, substr) {
		return nil
	}
	return fail("contains", value, substr, "value must contain substring")
}

func ContainsAny(value string, chars string) error {
	if strings.ContainsAny(value, chars) {
		return nil
	}
	return fail("containsany", value, chars, "value must contain at least one requested character")
}

func ContainsRune(value string, r rune) error {
	if strings.ContainsRune(value, r) {
		return nil
	}
	return fail("containsrune", value, string(r), "value must contain rune")
}

func EndsNotWith(value string, suffix string) error {
	if !strings.HasSuffix(value, suffix) {
		return nil
	}
	return fail("endsnotwith", value, suffix, "value must not end with suffix")
}

func EndsWith(value string, suffix string) error {
	if strings.HasSuffix(value, suffix) {
		return nil
	}
	return fail("endswith", value, suffix, "value must end with suffix")
}

func Excludes(value string, substr string) error {
	if !strings.Contains(value, substr) {
		return nil
	}
	return fail("excludes", value, substr, "value must not contain substring")
}

func ExcludesAll(value string, chars string) error {
	if !strings.ContainsAny(value, chars) {
		return nil
	}
	return fail("excludesall", value, chars, "value must not contain any requested character")
}

func ExcludesRune(value string, r rune) error {
	if !strings.ContainsRune(value, r) {
		return nil
	}
	return fail("excludesrune", value, string(r), "value must not contain rune")
}

func Lowercase(value string) error {
	if value == strings.ToLower(value) {
		return nil
	}
	return fail("lowercase", value, nil, "value must be lowercase")
}

func Multibyte(value string) error {
	for _, r := range value {
		if r > 127 {
			return nil
		}
	}
	return fail("multibyte", value, nil, "value must contain at least one multibyte character")
}

func Number(value string) error {
	if value == "" {
		return fail("number", value, nil, "value must be a number")
	}
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return nil
	}
	if _, err := strconv.ParseUint(value, 10, 64); err == nil {
		return nil
	}
	return fail("number", value, nil, "value must be an integer number")
}

func Numeric(value string) error {
	if _, err := strconv.ParseFloat(value, 64); err == nil && strings.TrimSpace(value) == value && value != "" {
		return nil
	}
	return fail("numeric", value, nil, "value must be numeric")
}

func PrintASCII(value string) error {
	for _, r := range value {
		if r < 32 || r > 126 {
			return fail("printascii", value, nil, "value must contain printable ASCII characters only")
		}
	}
	return nil
}

func StartsNotWith(value string, prefix string) error {
	if !strings.HasPrefix(value, prefix) {
		return nil
	}
	return fail("startsnotwith", value, prefix, "value must not start with prefix")
}

func StartsWith(value string, prefix string) error {
	if strings.HasPrefix(value, prefix) {
		return nil
	}
	return fail("startswith", value, prefix, "value must start with prefix")
}

func Uppercase(value string) error {
	if value == strings.ToUpper(value) {
		return nil
	}
	return fail("uppercase", value, nil, "value must be uppercase")
}

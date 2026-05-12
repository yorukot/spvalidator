package spvalidator

import "strings"

// Eq validates that value equals expected.
func Eq(value any, expected any) error {
	if equalAny(value, expected) {
		return nil
	}
	return failf("eq", value, expected, "must equal %v", expected)
}

// EqIgnoreCase validates that value equals expected case-insensitively.
func EqIgnoreCase(value string, expected string) error {
	if strings.EqualFold(value, expected) {
		return nil
	}
	return failf("eq_ignore_case", value, expected, "must equal %v ignoring case", expected)
}

// Ne validates that value does not equal disallowed.
func Ne(value any, disallowed any) error {
	if !equalAny(value, disallowed) {
		return nil
	}
	return failf("ne", value, disallowed, "must not equal %v", disallowed)
}

// NE is an alias for Ne.
func NE(value any, disallowed any) error { return Ne(value, disallowed) }

// NeIgnoreCase validates that value does not equal disallowed case-insensitively.
func NeIgnoreCase(value string, disallowed string) error {
	if !strings.EqualFold(value, disallowed) {
		return nil
	}
	return failf("ne_ignore_case", value, disallowed, "must not equal %v ignoring case", disallowed)
}

// Gt validates that value is greater than threshold.
func Gt(value any, threshold any) error {
	if cmp, ok := compareOrder(value, threshold); ok && cmp > 0 {
		return nil
	}
	return failf("gt", value, threshold, "must be greater than %v", threshold)
}

// GT is an alias for Gt.
func GT(value any, threshold any) error { return Gt(value, threshold) }

// Gte validates that value is greater than or equal to threshold.
func Gte(value any, threshold any) error {
	if cmp, ok := compareOrder(value, threshold); ok && cmp >= 0 {
		return nil
	}
	return failf("gte", value, threshold, "must be %v or greater", threshold)
}

// GTE is an alias for Gte.
func GTE(value any, threshold any) error { return Gte(value, threshold) }

// Lt validates that value is less than threshold.
func Lt(value any, threshold any) error {
	if cmp, ok := compareOrder(value, threshold); ok && cmp < 0 {
		return nil
	}
	return failf("lt", value, threshold, "must be less than %v", threshold)
}

// LT is an alias for Lt.
func LT(value any, threshold any) error { return Lt(value, threshold) }

// Lte validates that value is less than or equal to threshold.
func Lte(value any, threshold any) error {
	if cmp, ok := compareOrder(value, threshold); ok && cmp <= 0 {
		return nil
	}
	return failf("lte", value, threshold, "must be %v or less", threshold)
}

// LTE is an alias for Lte.
func LTE(value any, threshold any) error { return Lte(value, threshold) }

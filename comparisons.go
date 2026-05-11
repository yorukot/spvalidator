package spvalidator

import "strings"

// Eq validates that value equals expected.
func Eq(value any, expected any) error {
	if equalAny(value, expected) {
		return nil
	}
	return fail("eq", value, expected, "value must equal expected value")
}

// EqIgnoreCase validates that value equals expected case-insensitively.
func EqIgnoreCase(value string, expected string) error {
	if strings.EqualFold(value, expected) {
		return nil
	}
	return fail("eq_ignore_case", value, expected, "value must equal expected value ignoring case")
}

// Ne validates that value does not equal disallowed.
func Ne(value any, disallowed any) error {
	if !equalAny(value, disallowed) {
		return nil
	}
	return fail("ne", value, disallowed, "value must not equal disallowed value")
}

// NE is an alias for Ne.
func NE(value any, disallowed any) error { return Ne(value, disallowed) }

// NeIgnoreCase validates that value does not equal disallowed case-insensitively.
func NeIgnoreCase(value string, disallowed string) error {
	if !strings.EqualFold(value, disallowed) {
		return nil
	}
	return fail("ne_ignore_case", value, disallowed, "value must not equal disallowed value ignoring case")
}

// Gt validates that value is greater than threshold.
func Gt(value any, threshold any) error {
	if cmp, ok := compareOrder(value, threshold); ok && cmp > 0 {
		return nil
	}
	return fail("gt", value, threshold, "value must be greater than threshold")
}

// GT is an alias for Gt.
func GT(value any, threshold any) error { return Gt(value, threshold) }

// Gte validates that value is greater than or equal to threshold.
func Gte(value any, threshold any) error {
	if cmp, ok := compareOrder(value, threshold); ok && cmp >= 0 {
		return nil
	}
	return fail("gte", value, threshold, "value must be greater than or equal to threshold")
}

// GTE is an alias for Gte.
func GTE(value any, threshold any) error { return Gte(value, threshold) }

// Lt validates that value is less than threshold.
func Lt(value any, threshold any) error {
	if cmp, ok := compareOrder(value, threshold); ok && cmp < 0 {
		return nil
	}
	return fail("lt", value, threshold, "value must be less than threshold")
}

// LT is an alias for Lt.
func LT(value any, threshold any) error { return Lt(value, threshold) }

// Lte validates that value is less than or equal to threshold.
func Lte(value any, threshold any) error {
	if cmp, ok := compareOrder(value, threshold); ok && cmp <= 0 {
		return nil
	}
	return fail("lte", value, threshold, "value must be less than or equal to threshold")
}

// LTE is an alias for Lte.
func LTE(value any, threshold any) error { return Lte(value, threshold) }

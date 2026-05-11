package spvalidator

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestValidationError(t *testing.T) {
	t.Run("message", func(t *testing.T) {
		err := (&ValidationError{Tag: "eq", Value: "x", Param: "y", Message: "custom"}).Error()
		if err != "custom" {
			t.Fatalf("expected custom message, got %q", err)
		}
	})

	t.Run("parameter", func(t *testing.T) {
		err := (&ValidationError{Tag: "eq", Value: "x", Param: "y"}).Error()
		if err == "" {
			t.Fatal("expected non-empty message")
		}
	})

	t.Run("nil receiver", func(t *testing.T) {
		var err *ValidationError
		if got := err.Error(); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})
}

func TestIsValidationError(t *testing.T) {
	base := &ValidationError{Tag: "required", Value: "", Message: "value is required"}
	if !IsValidationError(base) {
		t.Fatal("expected ValidationError to be detected")
	}
	wrapped := fmt.Errorf("wrapped: %w", base)
	if !IsValidationError(wrapped) {
		t.Fatal("expected wrapped ValidationError to be detected")
	}
	if IsValidationError(errors.New("other")) {
		t.Fatal("expected non-validation error to be rejected")
	}
}

func TestEqualityValidators(t *testing.T) {
	t.Run("exact numeric equality", func(t *testing.T) {
		large := uint64(9007199254740992)
		if err := Eq(large, large); err != nil {
			t.Fatalf("expected exact numeric equality, got %v", err)
		}
		if err := Eq(large, large+1); err == nil {
			t.Fatal("expected distinct large integers to differ")
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		expectNoError(t, EqIgnoreCase("Hello", "hello"))
		expectNoError(t, NeIgnoreCase("Hello", "world"))
		if err := NeIgnoreCase("Hello", "HELLO"); err == nil {
			t.Fatal("expected equal values to fail NeIgnoreCase")
		}
	})

	t.Run("aliases", func(t *testing.T) {
		expectNoError(t, NE(1, 2))
		expectNoError(t, GT(3, 2))
		expectNoError(t, GTE(3, 3))
		expectNoError(t, LT(2, 3))
		expectNoError(t, LTE(3, 3))
	})
}

func TestOrderedValidators(t *testing.T) {
	now := time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC)
	later := now.Add(time.Hour)
	earlier := now.Add(-time.Hour)

	cases := []struct {
		name string
		err  error
		ok   bool
	}{
		{name: "gt number", err: Gt(3, 2), ok: true},
		{name: "gte number", err: Gte(3, 3), ok: true},
		{name: "lt number", err: Lt(2, 3), ok: true},
		{name: "lte number", err: Lte(3, 3), ok: true},
		{name: "gt time", err: Gt(later, now), ok: true},
		{name: "lt time", err: Lt(earlier, now), ok: true},
		{name: "gt string", err: Gt("b", "a"), ok: true},
		{name: "lt string", err: Lt("a", "b"), ok: true},
		{name: "gt bool", err: Gt(true, false), ok: true},
		{name: "lt bool", err: Lt(false, true), ok: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.ok {
				expectNoError(t, tc.err)
			} else {
				expectError(t, tc.err)
			}
		})
	}

	if err := Gt(struct{}{}, 1); err == nil {
		t.Fatal("expected incomparable types to fail")
	}
}

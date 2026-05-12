package spvalidator

import "testing"

func TestStringValidatorChain(t *testing.T) {
	const uuid = "123e4567-e89b-12d3-a456-426614174000"

	got, err := String("  " + uuid + "  ").
		TrimSpace().
		Required().
		Max(36).
		UUID().
		Value()

	expectNoError(t, err)
	if got != uuid {
		t.Fatalf("expected trimmed UUID %q, got %q", uuid, got)
	}
}

func TestStringValidatorErr(t *testing.T) {
	err := String("  ").
		TrimSpace().
		Required().
		UUID().
		Err()

	validateErr(t, err, "required")
}

func TestStringValidatorKeepsFirstError(t *testing.T) {
	err := String("123e4567-e89b-12d3-a456-426614174000").
		Max(10).
		UUID().
		Err()

	validateErr(t, err, "max")
}

func TestStringValidatorCheck(t *testing.T) {
	t.Run("custom validator", func(t *testing.T) {
		err := String("hello").
			Check(func(value string) error {
				return Contains(value, "ell")
			}).
			Err()

		expectNoError(t, err)
	})

	t.Run("nil validator", func(t *testing.T) {
		err := String("hello").Check(nil).Err()

		validateErr(t, err, "check")
	})
}

func TestStringValidatorMethodGroups(t *testing.T) {
	dir := t.TempDir()
	file := writeTempFile(t, dir, "sample.txt", []byte("hello"))

	cases := []struct {
		name string
		err  error
	}{
		{"string", String("Hello").Alpha().Err()},
		{"string_param", String("hello").Contains("ell").Err()},
		{"comparison", String("b").Gt("a").Lte("b").Err()},
		{"format", String("user@example.com").Email().Err()},
		{"format_param", String("2026/05/13").DateTime("2006/01/02").Err()},
		{"network", String("https://example.com").HTTPSURL().Err()},
		{"path", String(file).File().Err()},
		{"alias", String("#fff").IsColor().Err()},
		{"any_backed", String("4111111111111111").LuhnChecksum().Err()},
		{"postcode", String("12345").PostcodeISO3166Alpha2("US").Err()},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expectNoError(t, tc.err)
		})
	}
}

func TestValueValidatorChain(t *testing.T) {
	got, err := Any([]int{1, 2, 3}).
		Required().
		Len(3).
		Unique().
		Value()

	expectNoError(t, err)
	values, ok := got.([]int)
	if !ok || len(values) != 3 {
		t.Fatalf("expected original slice, got %#v", got)
	}

	expectNoError(t, Any(21).Required().Gte(18).Lte(120).Err())
	expectNoError(t, Any("8080").Port().Err())
	expectNoError(t, Any("25.0330").Latitude().Err())
	expectNoError(t, Any(methodValidator{ok: true}).ValidateFn().Err())
}

func TestValueValidatorKeepsFirstError(t *testing.T) {
	err := Any(10).
		Lt(5).
		Gt(20).
		Err()

	validateErr(t, err, "lt")
}

func TestValueValidatorCheck(t *testing.T) {
	t.Run("custom validator", func(t *testing.T) {
		err := Any(10).
			Check(func(value any) error {
				return Gte(value, 10)
			}).
			Err()

		expectNoError(t, err)
	})

	t.Run("nil validator", func(t *testing.T) {
		err := Any(10).Check(nil).Err()

		validateErr(t, err, "check")
	})
}

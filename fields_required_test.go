package spvalidator

import "testing"

type fieldTarget struct {
	Password string
	Confirm  string
	Enabled  bool
	Age      int
	Profile  struct {
		Name string
	}
	Meta map[string]string
}

func TestRequiredValidators(t *testing.T) {
	target := fieldTarget{
		Password: "secret",
		Confirm:  "",
		Enabled:  true,
		Age:      21,
		Meta:     map[string]string{"role": "admin", "limit": "18"},
	}
	target.Profile.Name = "alice"

	expectNoError(t, Required("value"))
	validateErr(t, Required(""), "required")

	t.Run("conditional required", func(t *testing.T) {
		validateErr(t, RequiredIf(target, "Confirm", Condition("Enabled", true)), "required_if")
		validateErr(t, RequiredUnless(target, "Confirm", Condition("Enabled", false)), "required_unless")
		validateErr(t, RequiredWith(target, "Confirm", "Password"), "required_with")
		validateErr(t, RequiredWithAll(target, "Confirm", "Password", "Enabled"), "required_with_all")
		validateErr(t, RequiredWithout(target, "Confirm", "Missing"), "required_without")
		validateErr(t, RequiredWithoutAll(target, "Confirm", "Missing1", "Missing2"), "required_without_all")
	})

	t.Run("conditional excluded", func(t *testing.T) {
		validateErr(t, ExcludedIf(target, "Password", Condition("Enabled", true)), "excluded_if")
		validateErr(t, ExcludedUnless(target, "Password", Condition("Enabled", false)), "excluded_unless")
		validateErr(t, ExcludedWith(target, "Password", "Enabled"), "excluded_with")
		validateErr(t, ExcludedWithAll(target, "Password", "Password", "Enabled"), "excluded_with_all")
		validateErr(t, ExcludedWithout(target, "Password", "Missing"), "excluded_without")
		validateErr(t, ExcludedWithoutAll(target, "Password", "Missing1", "Missing2"), "excluded_without_all")
	})
}

func TestFieldValidators(t *testing.T) {
	target := fieldTarget{
		Password: "secret",
		Confirm:  "secret",
		Enabled:  true,
		Age:      21,
		Meta:     map[string]string{"role": "admin", "name": "alice", "limit": "18"},
	}
	target.Profile.Name = "alice"

	expectNoError(t, EqField(target, "Profile.Name", "Meta.name"))
	expectNoError(t, NeField(target, "Profile.Name", "Meta.limit"))
	expectNoError(t, GtField(target, "Age", "Meta.limit"))
	expectNoError(t, GteField(target, "Age", "Meta.limit"))
	expectNoError(t, LtField(target, "Meta.limit", "Age"))
	expectNoError(t, LteField(target, "Meta.limit", "Age"))
	expectNoError(t, FieldContains(target, "Profile.Name", "lic"))
	expectNoError(t, FieldExcludes(target, "Profile.Name", "zzz"))

	if err := EqField(target, "Profile.Name", "Password"); err == nil {
		t.Fatal("expected unequal fields to fail")
	}
	if err := FieldContains(target, "Age", "1"); err == nil {
		t.Fatal("expected non-string field to fail FieldContains")
	}
}

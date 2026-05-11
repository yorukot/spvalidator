package spvalidator

import "testing"

func TestAliasValidators(t *testing.T) {
	expectNoError(t, IsColor("#fff"))
	validateErr(t, IsColor("not-a-color"), "iscolor")

	expectNoError(t, CountryCode("US"))
	validateErr(t, CountryCode("Z1"), "country_code")
}

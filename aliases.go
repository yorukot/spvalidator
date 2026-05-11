package spvalidator

func IsColor(value string) error {
	if HexColor(value) == nil || RGB(value) == nil || RGBA(value) == nil || HSL(value) == nil || HSLA(value) == nil || CMYK(value) == nil {
		return nil
	}
	return fail("iscolor", value, nil, "value must be a supported color")
}

func CountryCode(value string) error {
	if ISO3166Alpha2(value) == nil || ISO3166Alpha3(value) == nil || ISO3166AlphaNumeric(value) == nil {
		return nil
	}
	return fail("country_code", value, nil, "value must be an ISO 3166 country code")
}

package spvalidator

import "strings"

// String starts a fluent validation chain for a string value.
func String(value string) *StringValidator {
	return &StringValidator{value: value}
}

// Any starts a fluent validation chain for any value.
func Any(value any) *ValueValidator {
	return &ValueValidator{value: value}
}

// StringValidator validates and transforms a string through a fluent chain.
type StringValidator struct {
	value string
	err   error
}

// ValueValidator validates any value through a fluent chain.
type ValueValidator struct {
	value any
	err   error
}

// Value returns the current string value and the first validation error.
func (v *StringValidator) Value() (string, error) {
	return v.value, v.err
}

// Err returns the first validation error in the chain.
func (v *StringValidator) Err() error {
	return v.err
}

// Check runs fn against the current value when the chain has not failed.
func (v *StringValidator) Check(fn func(string) error) *StringValidator {
	if v.err != nil {
		return v
	}
	if fn == nil {
		v.err = fail("check", v.value, nil, "validator function is nil")
		return v
	}
	v.err = fn(v.value)
	return v
}

func (v *StringValidator) checkAny(fn func(any) error) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = fn(v.value)
	return v
}

// TrimSpace removes leading and trailing white space from the current value.
func (v *StringValidator) TrimSpace() *StringValidator {
	if v.err != nil {
		return v
	}
	v.value = strings.TrimSpace(v.value)
	return v
}

// Required validates that the current value is not empty.
func (v *StringValidator) Required() *StringValidator { return v.checkAny(Required) }

// IsDefault validates that the current value is the zero value.
func (v *StringValidator) IsDefault() *StringValidator { return v.checkAny(IsDefault) }

// Len validates that the current value has exactly length runes.
func (v *StringValidator) Len(length int) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Len(v.value, length)
	return v
}

// Max validates that the current value length is at most max.
func (v *StringValidator) Max(max any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Max(v.value, max)
	return v
}

// Min validates that the current value length is at least min.
func (v *StringValidator) Min(min any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Min(v.value, min)
	return v
}

// OneOf validates that the current value equals one of choices.
func (v *StringValidator) OneOf(choices ...any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = OneOf(v.value, choices...)
	return v
}

// NoneOf validates that the current value does not equal any choice.
func (v *StringValidator) NoneOf(choices ...any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = NoneOf(v.value, choices...)
	return v
}

// Unique validates that the current value contains unique runes.
func (v *StringValidator) Unique() *StringValidator { return v.checkAny(Unique) }

// Eq validates that the current value equals expected.
func (v *StringValidator) Eq(expected any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Eq(v.value, expected)
	return v
}

// EqIgnoreCase validates that the current value equals expected case-insensitively.
func (v *StringValidator) EqIgnoreCase(expected string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = EqIgnoreCase(v.value, expected)
	return v
}

// Ne validates that the current value does not equal disallowed.
func (v *StringValidator) Ne(disallowed any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Ne(v.value, disallowed)
	return v
}

// NeIgnoreCase validates that the current value does not equal disallowed case-insensitively.
func (v *StringValidator) NeIgnoreCase(disallowed string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = NeIgnoreCase(v.value, disallowed)
	return v
}

// Gt validates that the current value is greater than threshold.
func (v *StringValidator) Gt(threshold any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Gt(v.value, threshold)
	return v
}

// Gte validates that the current value is greater than or equal to threshold.
func (v *StringValidator) Gte(threshold any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Gte(v.value, threshold)
	return v
}

// Lt validates that the current value is less than threshold.
func (v *StringValidator) Lt(threshold any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Lt(v.value, threshold)
	return v
}

// Lte validates that the current value is less than or equal to threshold.
func (v *StringValidator) Lte(threshold any) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Lte(v.value, threshold)
	return v
}

// Port validates that the current value is a TCP/UDP port.
func (v *StringValidator) Port() *StringValidator { return v.checkAny(Port) }

// Latitude validates that the current value is a latitude.
func (v *StringValidator) Latitude() *StringValidator { return v.checkAny(Latitude) }

// Longitude validates that the current value is a longitude.
func (v *StringValidator) Longitude() *StringValidator { return v.checkAny(Longitude) }

// LuhnChecksum validates that the current value passes the Luhn checksum.
func (v *StringValidator) LuhnChecksum() *StringValidator { return v.checkAny(LuhnChecksum) }

func (v *StringValidator) Alpha() *StringValidator         { return v.Check(Alpha) }
func (v *StringValidator) AlphaSpace() *StringValidator    { return v.Check(AlphaSpace) }
func (v *StringValidator) Alphanum() *StringValidator      { return v.Check(Alphanum) }
func (v *StringValidator) AlphanumSpace() *StringValidator { return v.Check(AlphanumSpace) }
func (v *StringValidator) AlphanumUnicode() *StringValidator {
	return v.Check(AlphanumUnicode)
}
func (v *StringValidator) AlphaUnicode() *StringValidator { return v.Check(AlphaUnicode) }
func (v *StringValidator) ASCII() *StringValidator        { return v.Check(ASCII) }
func (v *StringValidator) Boolean() *StringValidator      { return v.Check(Boolean) }
func (v *StringValidator) Contains(substr string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Contains(v.value, substr)
	return v
}
func (v *StringValidator) ContainsAny(chars string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = ContainsAny(v.value, chars)
	return v
}
func (v *StringValidator) ContainsRune(r rune) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = ContainsRune(v.value, r)
	return v
}
func (v *StringValidator) EndsNotWith(suffix string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = EndsNotWith(v.value, suffix)
	return v
}
func (v *StringValidator) EndsWith(suffix string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = EndsWith(v.value, suffix)
	return v
}
func (v *StringValidator) Excludes(substr string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = Excludes(v.value, substr)
	return v
}
func (v *StringValidator) ExcludesAll(chars string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = ExcludesAll(v.value, chars)
	return v
}
func (v *StringValidator) ExcludesRune(r rune) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = ExcludesRune(v.value, r)
	return v
}
func (v *StringValidator) Lowercase() *StringValidator  { return v.Check(Lowercase) }
func (v *StringValidator) Multibyte() *StringValidator  { return v.Check(Multibyte) }
func (v *StringValidator) Number() *StringValidator     { return v.Check(Number) }
func (v *StringValidator) Numeric() *StringValidator    { return v.Check(Numeric) }
func (v *StringValidator) PrintASCII() *StringValidator { return v.Check(PrintASCII) }
func (v *StringValidator) StartsNotWith(prefix string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = StartsNotWith(v.value, prefix)
	return v
}
func (v *StringValidator) StartsWith(prefix string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = StartsWith(v.value, prefix)
	return v
}
func (v *StringValidator) Uppercase() *StringValidator { return v.Check(Uppercase) }

func (v *StringValidator) IsColor() *StringValidator     { return v.Check(IsColor) }
func (v *StringValidator) CountryCode() *StringValidator { return v.Check(CountryCode) }

func (v *StringValidator) Dir() *StringValidator      { return v.Check(Dir) }
func (v *StringValidator) DirPath() *StringValidator  { return v.Check(DirPath) }
func (v *StringValidator) File() *StringValidator     { return v.Check(File) }
func (v *StringValidator) FilePath() *StringValidator { return v.Check(FilePath) }
func (v *StringValidator) Image() *StringValidator    { return v.Check(Image) }
func (v *StringValidator) MIMEType() *StringValidator { return v.Check(MIMEType) }

func (v *StringValidator) CIDR() *StringValidator            { return v.Check(CIDR) }
func (v *StringValidator) CIDRv4() *StringValidator          { return v.Check(CIDRv4) }
func (v *StringValidator) CIDRv6() *StringValidator          { return v.Check(CIDRv6) }
func (v *StringValidator) DataURI() *StringValidator         { return v.Check(DataURI) }
func (v *StringValidator) Domain() *StringValidator          { return v.Check(Domain) }
func (v *StringValidator) FQDN() *StringValidator            { return v.Check(FQDN) }
func (v *StringValidator) Hostname() *StringValidator        { return v.Check(Hostname) }
func (v *StringValidator) HostnameRFC1123() *StringValidator { return v.Check(HostnameRFC1123) }
func (v *StringValidator) HostnamePort() *StringValidator    { return v.Check(HostnamePort) }
func (v *StringValidator) IP() *StringValidator              { return v.Check(IP) }
func (v *StringValidator) IPAddr() *StringValidator          { return v.Check(IPAddr) }
func (v *StringValidator) IP4Addr() *StringValidator         { return v.Check(IP4Addr) }
func (v *StringValidator) IP6Addr() *StringValidator         { return v.Check(IP6Addr) }
func (v *StringValidator) IPv4() *StringValidator            { return v.Check(IPv4) }
func (v *StringValidator) IPv6() *StringValidator            { return v.Check(IPv6) }
func (v *StringValidator) MAC() *StringValidator             { return v.Check(MAC) }
func (v *StringValidator) TCPAddr() *StringValidator         { return v.Check(TCPAddr) }
func (v *StringValidator) TCP4Addr() *StringValidator        { return v.Check(TCP4Addr) }
func (v *StringValidator) TCP6Addr() *StringValidator        { return v.Check(TCP6Addr) }
func (v *StringValidator) UDPAddr() *StringValidator         { return v.Check(UDPAddr) }
func (v *StringValidator) UDP4Addr() *StringValidator        { return v.Check(UDP4Addr) }
func (v *StringValidator) UDP6Addr() *StringValidator        { return v.Check(UDP6Addr) }
func (v *StringValidator) UnixAddr() *StringValidator        { return v.Check(UnixAddr) }
func (v *StringValidator) UDSExists() *StringValidator       { return v.Check(UDSExists) }
func (v *StringValidator) URI() *StringValidator             { return v.Check(URI) }
func (v *StringValidator) URL() *StringValidator             { return v.Check(URL) }
func (v *StringValidator) HTTPURL() *StringValidator         { return v.Check(HTTPURL) }
func (v *StringValidator) HTTPSURL() *StringValidator        { return v.Check(HTTPSURL) }
func (v *StringValidator) Origin() *StringValidator          { return v.Check(Origin) }
func (v *StringValidator) URLEncoded() *StringValidator      { return v.Check(URLEncoded) }
func (v *StringValidator) URNRFC2141() *StringValidator      { return v.Check(URNRFC2141) }

func (v *StringValidator) Base64() *StringValidator       { return v.Check(Base64) }
func (v *StringValidator) Base64URL() *StringValidator    { return v.Check(Base64URL) }
func (v *StringValidator) Base64RawURL() *StringValidator { return v.Check(Base64RawURL) }
func (v *StringValidator) BICISO93622014() *StringValidator {
	return v.Check(BICISO93622014)
}
func (v *StringValidator) BIC() *StringValidator              { return v.Check(BIC) }
func (v *StringValidator) BCP47LanguageTag() *StringValidator { return v.Check(BCP47LanguageTag) }
func (v *StringValidator) BCP47StrictLanguageTag() *StringValidator {
	return v.Check(BCP47StrictLanguageTag)
}
func (v *StringValidator) BTCAddr() *StringValidator       { return v.Check(BTCAddr) }
func (v *StringValidator) BTCAddrBech32() *StringValidator { return v.Check(BTCAddrBech32) }
func (v *StringValidator) CreditCard() *StringValidator    { return v.Check(CreditCard) }
func (v *StringValidator) MongoDB() *StringValidator       { return v.Check(MongoDB) }
func (v *StringValidator) MongoDBConnectionString() *StringValidator {
	return v.Check(MongoDBConnectionString)
}
func (v *StringValidator) Cron() *StringValidator                { return v.Check(Cron) }
func (v *StringValidator) SpiceDB() *StringValidator             { return v.Check(SpiceDB) }
func (v *StringValidator) E164() *StringValidator                { return v.Check(E164) }
func (v *StringValidator) EIN() *StringValidator                 { return v.Check(EIN) }
func (v *StringValidator) Email() *StringValidator               { return v.Check(Email) }
func (v *StringValidator) ETHAddr() *StringValidator             { return v.Check(ETHAddr) }
func (v *StringValidator) Hexadecimal() *StringValidator         { return v.Check(Hexadecimal) }
func (v *StringValidator) HexColor() *StringValidator            { return v.Check(HexColor) }
func (v *StringValidator) HSL() *StringValidator                 { return v.Check(HSL) }
func (v *StringValidator) HSLA() *StringValidator                { return v.Check(HSLA) }
func (v *StringValidator) CMYK() *StringValidator                { return v.Check(CMYK) }
func (v *StringValidator) HTML() *StringValidator                { return v.Check(HTML) }
func (v *StringValidator) HTMLEncoded() *StringValidator         { return v.Check(HTMLEncoded) }
func (v *StringValidator) ISBN() *StringValidator                { return v.Check(ISBN) }
func (v *StringValidator) ISBN10() *StringValidator              { return v.Check(ISBN10) }
func (v *StringValidator) ISBN13() *StringValidator              { return v.Check(ISBN13) }
func (v *StringValidator) ISSN() *StringValidator                { return v.Check(ISSN) }
func (v *StringValidator) ISO3166Alpha2() *StringValidator       { return v.Check(ISO3166Alpha2) }
func (v *StringValidator) ISO3166Alpha3() *StringValidator       { return v.Check(ISO3166Alpha3) }
func (v *StringValidator) ISO3166AlphaNumeric() *StringValidator { return v.Check(ISO3166AlphaNumeric) }
func (v *StringValidator) ISO31662() *StringValidator            { return v.Check(ISO31662) }
func (v *StringValidator) ISO4217() *StringValidator             { return v.Check(ISO4217) }
func (v *StringValidator) JSON() *StringValidator                { return v.Check(JSON) }
func (v *StringValidator) JWT() *StringValidator                 { return v.Check(JWT) }
func (v *StringValidator) RGB() *StringValidator                 { return v.Check(RGB) }
func (v *StringValidator) RGBA() *StringValidator                { return v.Check(RGBA) }
func (v *StringValidator) SSN() *StringValidator                 { return v.Check(SSN) }
func (v *StringValidator) Timezone() *StringValidator            { return v.Check(Timezone) }
func (v *StringValidator) UUID() *StringValidator                { return v.Check(UUID) }
func (v *StringValidator) UUID3() *StringValidator               { return v.Check(UUID3) }
func (v *StringValidator) UUID3RFC4122() *StringValidator        { return v.Check(UUID3RFC4122) }
func (v *StringValidator) UUID4() *StringValidator               { return v.Check(UUID4) }
func (v *StringValidator) UUID4RFC4122() *StringValidator        { return v.Check(UUID4RFC4122) }
func (v *StringValidator) UUID5() *StringValidator               { return v.Check(UUID5) }
func (v *StringValidator) UUID5RFC4122() *StringValidator        { return v.Check(UUID5RFC4122) }
func (v *StringValidator) UUIDRFC4122() *StringValidator         { return v.Check(UUIDRFC4122) }
func (v *StringValidator) MD4() *StringValidator                 { return v.Check(MD4) }
func (v *StringValidator) MD5() *StringValidator                 { return v.Check(MD5) }
func (v *StringValidator) SHA256() *StringValidator              { return v.Check(SHA256) }
func (v *StringValidator) SHA384() *StringValidator              { return v.Check(SHA384) }
func (v *StringValidator) SHA512() *StringValidator              { return v.Check(SHA512) }
func (v *StringValidator) RIPEMD128() *StringValidator           { return v.Check(RIPEMD128) }
func (v *StringValidator) RIPEMD160() *StringValidator           { return v.Check(RIPEMD160) }
func (v *StringValidator) TIGER128() *StringValidator            { return v.Check(TIGER128) }
func (v *StringValidator) TIGER160() *StringValidator            { return v.Check(TIGER160) }
func (v *StringValidator) TIGER192() *StringValidator            { return v.Check(TIGER192) }
func (v *StringValidator) SemVer() *StringValidator              { return v.Check(SemVer) }
func (v *StringValidator) ULID() *StringValidator                { return v.Check(ULID) }
func (v *StringValidator) CVE() *StringValidator                 { return v.Check(CVE) }

func (v *StringValidator) DateTime(layouts ...string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = DateTime(v.value, layouts...)
	return v
}

func (v *StringValidator) PostcodeISO3166Alpha2(country string) *StringValidator {
	if v.err != nil {
		return v
	}
	v.err = PostcodeISO3166Alpha2(v.value, country)
	return v
}

// Value returns the current value and the first validation error.
func (v *ValueValidator) Value() (any, error) {
	return v.value, v.err
}

// Err returns the first validation error in the chain.
func (v *ValueValidator) Err() error {
	return v.err
}

// Check runs fn against the current value when the chain has not failed.
func (v *ValueValidator) Check(fn func(any) error) *ValueValidator {
	if v.err != nil {
		return v
	}
	if fn == nil {
		v.err = fail("check", v.value, nil, "validator function is nil")
		return v
	}
	v.err = fn(v.value)
	return v
}

func (v *ValueValidator) Required() *ValueValidator     { return v.Check(Required) }
func (v *ValueValidator) IsDefault() *ValueValidator    { return v.Check(IsDefault) }
func (v *ValueValidator) Unique() *ValueValidator       { return v.Check(Unique) }
func (v *ValueValidator) Latitude() *ValueValidator     { return v.Check(Latitude) }
func (v *ValueValidator) Longitude() *ValueValidator    { return v.Check(Longitude) }
func (v *ValueValidator) LuhnChecksum() *ValueValidator { return v.Check(LuhnChecksum) }
func (v *ValueValidator) Port() *ValueValidator         { return v.Check(Port) }

func (v *ValueValidator) Len(length int) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Len(v.value, length)
	return v
}

func (v *ValueValidator) Max(max any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Max(v.value, max)
	return v
}

func (v *ValueValidator) Min(min any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Min(v.value, min)
	return v
}

func (v *ValueValidator) OneOf(choices ...any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = OneOf(v.value, choices...)
	return v
}

func (v *ValueValidator) NoneOf(choices ...any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = NoneOf(v.value, choices...)
	return v
}

func (v *ValueValidator) ValidateFn(methodName ...string) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = ValidateFn(v.value, methodName...)
	return v
}

func (v *ValueValidator) Eq(expected any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Eq(v.value, expected)
	return v
}

func (v *ValueValidator) Ne(disallowed any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Ne(v.value, disallowed)
	return v
}

func (v *ValueValidator) Gt(threshold any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Gt(v.value, threshold)
	return v
}

func (v *ValueValidator) Gte(threshold any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Gte(v.value, threshold)
	return v
}

func (v *ValueValidator) Lt(threshold any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Lt(v.value, threshold)
	return v
}

func (v *ValueValidator) Lte(threshold any) *ValueValidator {
	if v.err != nil {
		return v
	}
	v.err = Lte(v.value, threshold)
	return v
}

package spvalidator

import (
	"encoding/base64"
	"encoding/json"
	"html"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	bicPattern       = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
	bcp47Pattern     = regexp.MustCompile(`^[A-Za-z]{2,3}(-[A-Za-z]{4})?(-([A-Za-z]{2}|[0-9]{3}))?(-([A-Za-z0-9]{5,8}|[0-9][A-Za-z0-9]{3}))*$`)
	btcLegacyPattern = regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`)
	btcBech32Pattern = regexp.MustCompile(`^(bc1|tb1)[023456789acdefghjklmnpqrstuvwxyz]{11,71}$`)
	mongoIDPattern   = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)
	e164Pattern      = regexp.MustCompile(`^\+[1-9][0-9]{1,14}$`)
	einPattern       = regexp.MustCompile(`^[0-9]{2}-?[0-9]{7}$`)
	ethPattern       = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	hexPattern       = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	hexColorPattern  = regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	rgbPattern       = regexp.MustCompile(`(?i)^rgb\(\s*(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\s*,\s*(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\s*,\s*(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\s*\)$`)
	rgbaPattern      = regexp.MustCompile(`(?i)^rgba\(\s*(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\s*,\s*(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\s*,\s*(25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\s*,\s*(0|1|0?\.[0-9]+)\s*\)$`)
	hslPattern       = regexp.MustCompile(`(?i)^hsl\(\s*(360|3[0-5][0-9]|[12]?[0-9]{1,2})\s*,\s*(100|[0-9]{1,2})%\s*,\s*(100|[0-9]{1,2})%\s*\)$`)
	hslaPattern      = regexp.MustCompile(`(?i)^hsla\(\s*(360|3[0-5][0-9]|[12]?[0-9]{1,2})\s*,\s*(100|[0-9]{1,2})%\s*,\s*(100|[0-9]{1,2})%\s*,\s*(0|1|0?\.[0-9]+)\s*\)$`)
	cmykPattern      = regexp.MustCompile(`(?i)^cmyk\(\s*(100|[0-9]{1,2})%\s*,\s*(100|[0-9]{1,2})%\s*,\s*(100|[0-9]{1,2})%\s*,\s*(100|[0-9]{1,2})%\s*\)$`)
	htmlTagPattern   = regexp.MustCompile(`(?is)<[a-z][^>]*>.*</[a-z][^>]*>|<[a-z][^>]*/>`)
	ssnPattern       = regexp.MustCompile(`^([0-9]{3})-?([0-9]{2})-?([0-9]{4})$`)
	uuidPattern      = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	semverPattern    = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-((?:0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)(?:\.(?:0|[1-9][0-9]*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*))*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)
	ulidPattern      = regexp.MustCompile(`^[0-7][0-9A-HJKMNP-TV-Z]{25}$`)
	cvePattern       = regexp.MustCompile(`^CVE-[0-9]{4}-[0-9]{4,}$`)
)

func Base64(value string) error {
	if value != "" {
		if _, err := base64.StdEncoding.DecodeString(value); err == nil {
			return nil
		}
	}
	return fail("base64", value, nil, "value must be standard base64")
}

func Base64URL(value string) error {
	if value != "" {
		if _, err := base64.URLEncoding.DecodeString(value); err == nil {
			return nil
		}
	}
	return fail("base64url", value, nil, "value must be padded base64 URL encoding")
}

func Base64RawURL(value string) error {
	if value != "" {
		if _, err := base64.RawURLEncoding.DecodeString(value); err == nil {
			return nil
		}
	}
	return fail("base64rawurl", value, nil, "value must be raw base64 URL encoding")
}

func BICISO93622014(value string) error { return bic("bic_iso_9362_2014", value) }

func BIC(value string) error { return bic("bic", value) }

func bic(tag string, value string) error {
	if bicPattern.MatchString(value) {
		return nil
	}
	return fail(tag, value, nil, "value must be a BIC")
}

func BCP47LanguageTag(value string) error {
	if bcp47Pattern.MatchString(value) {
		return nil
	}
	return fail("bcp47_language_tag", value, nil, "value must be a BCP 47 language tag")
}

func BCP47StrictLanguageTag(value string) error {
	if bcp47Pattern.MatchString(value) && !strings.Contains(value, "_") {
		return nil
	}
	return fail("bcp47_strict_language_tag", value, nil, "value must be a strict BCP 47 language tag")
}

func BTCAddr(value string) error {
	if BTCAddrBech32(value) == nil || btcLegacyPattern.MatchString(value) {
		return nil
	}
	return fail("btc_addr", value, nil, "value must be a Bitcoin address")
}

func BTCAddrBech32(value string) error {
	if btcBech32Pattern.MatchString(strings.ToLower(value)) {
		return nil
	}
	return fail("btc_addr_bech32", value, nil, "value must be a Bitcoin Bech32 address")
}

func CreditCard(value string) error {
	digits := digitsOnly(value)
	if len(digits) < 12 || len(digits) > 19 || LuhnChecksum(digits) != nil {
		return fail("credit_card", value, nil, "value must be a valid credit card number")
	}
	return nil
}

func MongoDB(value string) error {
	if mongoIDPattern.MatchString(value) {
		return nil
	}
	return fail("mongodb", value, nil, "value must be a MongoDB ObjectID")
}

func MongoDBConnectionString(value string) error {
	parsed, err := url.Parse(value)
	if err == nil && (parsed.Scheme == "mongodb" || parsed.Scheme == "mongodb+srv") && parsed.Host != "" {
		return nil
	}
	return fail("mongodb_connection_string", value, nil, "value must be a MongoDB connection string")
}

func Cron(value string) error {
	fields := strings.Fields(value)
	if len(fields) != 5 && len(fields) != 6 {
		return fail("cron", value, nil, "cron expression must have 5 or 6 fields")
	}
	for i, field := range fields {
		if !validCronField(field, i, len(fields)) {
			return fail("cron", value, nil, "cron expression contains an invalid field")
		}
	}
	return nil
}

func SpiceDB(value string) error {
	if value == "" {
		return fail("spicedb", value, nil, "value must be a SpiceDB identifier")
	}
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune("_-:/#.", r) {
			return fail("spicedb", value, nil, "value must be a SpiceDB identifier")
		}
	}
	return nil
}

func DateTime(value string, layouts ...string) error {
	if len(layouts) == 0 {
		layouts = []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05", "2006-01-02"}
	}
	for _, layout := range layouts {
		if _, err := time.Parse(layout, value); err == nil {
			return nil
		}
	}
	return fail("datetime", value, layouts, "value must be a datetime")
}

func E164(value string) error {
	if e164Pattern.MatchString(value) {
		return nil
	}
	return fail("e164", value, nil, "value must be an E.164 phone number")
}

func EIN(value string) error {
	if einPattern.MatchString(value) {
		return nil
	}
	return fail("ein", value, nil, "value must be a US EIN")
}

func Email(value string) error {
	address, err := mail.ParseAddress(value)
	if err == nil && address.Address == value && strings.Contains(address.Address, "@") {
		return nil
	}
	return fail("email", value, nil, "value must be an email address")
}

func ETHAddr(value string) error {
	if ethPattern.MatchString(value) {
		return nil
	}
	return fail("eth_addr", value, nil, "value must be an Ethereum address")
}

func Hexadecimal(value string) error {
	if value != "" && hexPattern.MatchString(value) {
		return nil
	}
	return fail("hexadecimal", value, nil, "value must be hexadecimal")
}

func HexColor(value string) error {
	if hexColorPattern.MatchString(value) {
		return nil
	}
	return fail("hexcolor", value, nil, "value must be a hex color")
}

func HSL(value string) error {
	if hslPattern.MatchString(value) {
		return nil
	}
	return fail("hsl", value, nil, "value must be an HSL color")
}

func HSLA(value string) error {
	if hslaPattern.MatchString(value) {
		return nil
	}
	return fail("hsla", value, nil, "value must be an HSLA color")
}

func CMYK(value string) error {
	if cmykPattern.MatchString(value) {
		return nil
	}
	return fail("cmyk", value, nil, "value must be a CMYK color")
}

func HTML(value string) error {
	if htmlTagPattern.MatchString(value) {
		return nil
	}
	return fail("html", value, nil, "value must contain HTML tags")
}

func HTMLEncoded(value string) error {
	if html.UnescapeString(value) != value {
		return nil
	}
	return fail("html_encoded", value, nil, "value must contain HTML entities")
}

func ISBN(value string) error {
	if ISBN10(value) == nil || ISBN13(value) == nil {
		return nil
	}
	return fail("isbn", value, nil, "value must be an ISBN")
}

func ISBN10(value string) error {
	code := strings.ToUpper(strings.NewReplacer("-", "", " ", "").Replace(value))
	if len(code) != 10 {
		return fail("isbn10", value, nil, "value must be an ISBN-10")
	}
	sum := 0
	for i, r := range code {
		var digit int
		if i == 9 && r == 'X' {
			digit = 10
		} else if r >= '0' && r <= '9' {
			digit = int(r - '0')
		} else {
			return fail("isbn10", value, nil, "value must be an ISBN-10")
		}
		sum += digit * (10 - i)
	}
	if sum%11 == 0 {
		return nil
	}
	return fail("isbn10", value, nil, "value must be an ISBN-10")
}

func ISBN13(value string) error {
	code := strings.NewReplacer("-", "", " ", "").Replace(value)
	if len(code) != 13 || !allDigits(code) {
		return fail("isbn13", value, nil, "value must be an ISBN-13")
	}
	sum := 0
	for i, r := range code[:12] {
		digit := int(r - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	check := (10 - (sum % 10)) % 10
	if check == int(code[12]-'0') {
		return nil
	}
	return fail("isbn13", value, nil, "value must be an ISBN-13")
}

func ISSN(value string) error {
	code := strings.ToUpper(strings.ReplaceAll(value, "-", ""))
	if len(code) != 8 {
		return fail("issn", value, nil, "value must be an ISSN")
	}
	sum := 0
	for i, r := range code {
		var digit int
		if i == 7 && r == 'X' {
			digit = 10
		} else if r >= '0' && r <= '9' {
			digit = int(r - '0')
		} else {
			return fail("issn", value, nil, "value must be an ISSN")
		}
		sum += digit * (8 - i)
	}
	if sum%11 == 0 {
		return nil
	}
	return fail("issn", value, nil, "value must be an ISSN")
}

func ISO3166Alpha2(value string) error {
	if regexp.MustCompile(`^[A-Z]{2}$`).MatchString(value) {
		return nil
	}
	return fail("iso3166_1_alpha2", value, nil, "value must be an ISO 3166-1 alpha-2 code")
}

func ISO3166Alpha3(value string) error {
	if regexp.MustCompile(`^[A-Z]{3}$`).MatchString(value) {
		return nil
	}
	return fail("iso3166_1_alpha3", value, nil, "value must be an ISO 3166-1 alpha-3 code")
}

func ISO3166AlphaNumeric(value string) error {
	if regexp.MustCompile(`^[0-9]{3}$`).MatchString(value) {
		return nil
	}
	return fail("iso3166_1_alpha_numeric", value, nil, "value must be an ISO 3166-1 numeric code")
}

func ISO31662(value string) error {
	if regexp.MustCompile(`^[A-Z]{2}-[A-Z0-9]{1,6}$`).MatchString(value) {
		return nil
	}
	return fail("iso3166_2", value, nil, "value must be an ISO 3166-2 subdivision code")
}

func ISO4217(value string) error {
	if regexp.MustCompile(`^[A-Z]{3}$`).MatchString(value) {
		return nil
	}
	return fail("iso4217", value, nil, "value must be an ISO 4217 currency code")
}

func JSON(value string) error {
	if json.Valid([]byte(value)) {
		return nil
	}
	return fail("json", value, nil, "value must be JSON")
}

func JWT(value string) error {
	parts := strings.Split(value, ".")
	if len(parts) != 3 || parts[2] == "" {
		return fail("jwt", value, nil, "value must be a JWT")
	}
	for _, part := range parts[:2] {
		decoded, err := base64.RawURLEncoding.DecodeString(part)
		if err != nil || !json.Valid(decoded) {
			return fail("jwt", value, nil, "value must be a JWT")
		}
	}
	return nil
}

func Latitude(value any) error {
	if n, ok := numberValue(value); ok && n >= -90 && n <= 90 {
		return nil
	}
	return fail("latitude", value, nil, "value must be a latitude")
}

func Longitude(value any) error {
	if n, ok := numberValue(value); ok && n >= -180 && n <= 180 {
		return nil
	}
	return fail("longitude", value, nil, "value must be a longitude")
}

func LuhnChecksum(value any) error {
	digits := digitsOnlyValue(value)
	if digits == "" || !luhnValid(digits) {
		return fail("luhn_checksum", value, nil, "value must pass the Luhn checksum")
	}
	return nil
}

func PostcodeISO3166Alpha2(value string, country string) error {
	country = strings.ToUpper(country)
	var pattern *regexp.Regexp
	switch country {
	case "US":
		pattern = regexp.MustCompile(`^[0-9]{5}(-[0-9]{4})?$`)
	case "CA":
		pattern = regexp.MustCompile(`^[A-Z][0-9][A-Z][ -]?[0-9][A-Z][0-9]$`)
	case "GB":
		pattern = regexp.MustCompile(`^[A-Z]{1,2}[0-9][A-Z0-9]?\s?[0-9][A-Z]{2}$`)
	case "JP":
		pattern = regexp.MustCompile(`^[0-9]{3}-?[0-9]{4}$`)
	case "DE", "FR", "IT", "ES":
		pattern = regexp.MustCompile(`^[0-9]{5}$`)
	default:
		pattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 -]{1,15}$`)
	}
	if pattern.MatchString(strings.ToUpper(value)) {
		return nil
	}
	return fail("postcode_iso3166_alpha2", value, country, "value must be a postcode for the country")
}

func PostcodeISO3166Alpha2Field(target any, fieldPath string, countryFieldPath string) error {
	value, err := fieldInterface(target, fieldPath)
	if err != nil {
		return fail("postcode_iso3166_alpha2_field", target, fieldPath, err.Error())
	}
	country, err := fieldInterface(target, countryFieldPath)
	if err != nil {
		return fail("postcode_iso3166_alpha2_field", target, countryFieldPath, err.Error())
	}
	vs, ok := stringValue(value)
	if !ok {
		return fail("postcode_iso3166_alpha2_field", value, fieldPath, "postcode field must be string-like")
	}
	cs, ok := stringValue(country)
	if !ok {
		return fail("postcode_iso3166_alpha2_field", country, countryFieldPath, "country field must be string-like")
	}
	if PostcodeISO3166Alpha2(vs, cs) == nil {
		return nil
	}
	return fail("postcode_iso3166_alpha2_field", vs, cs, "postcode field must match country field")
}

func RGB(value string) error {
	if rgbPattern.MatchString(value) {
		return nil
	}
	return fail("rgb", value, nil, "value must be an RGB color")
}

func RGBA(value string) error {
	if rgbaPattern.MatchString(value) {
		return nil
	}
	return fail("rgba", value, nil, "value must be an RGBA color")
}

func SSN(value string) error {
	match := ssnPattern.FindStringSubmatch(value)
	if match == nil {
		return fail("ssn", value, nil, "value must be a US SSN")
	}
	area, _ := strconv.Atoi(match[1])
	if area == 0 || area == 666 || area >= 900 || match[2] == "00" || match[3] == "0000" {
		return fail("ssn", value, nil, "value must be a US SSN")
	}
	return nil
}

func Timezone(value string) error {
	if _, err := time.LoadLocation(value); err == nil {
		return nil
	}
	return fail("timezone", value, nil, "value must be an IANA timezone")
}

func UUID(value string) error {
	if uuidPattern.MatchString(value) {
		return nil
	}
	return fail("uuid", value, nil, "value must be a UUID")
}

func UUID3(value string) error { return uuidVersion("uuid3", value, '3', false) }

func UUID3RFC4122(value string) error { return uuidVersion("uuid3_rfc4122", value, '3', true) }

func UUID4(value string) error { return uuidVersion("uuid4", value, '4', false) }

func UUID4RFC4122(value string) error { return uuidVersion("uuid4_rfc4122", value, '4', true) }

func UUID5(value string) error { return uuidVersion("uuid5", value, '5', false) }

func UUID5RFC4122(value string) error { return uuidVersion("uuid5_rfc4122", value, '5', true) }

func UUIDRFC4122(value string) error {
	if UUID(value) == nil && isRFC4122UUID(value) {
		return nil
	}
	return fail("uuid_rfc4122", value, nil, "value must be an RFC 4122 UUID")
}

func MD4(value string) error       { return hexLength("md4", value, 32) }
func MD5(value string) error       { return hexLength("md5", value, 32) }
func SHA256(value string) error    { return hexLength("sha256", value, 64) }
func SHA384(value string) error    { return hexLength("sha384", value, 96) }
func SHA512(value string) error    { return hexLength("sha512", value, 128) }
func RIPEMD128(value string) error { return hexLength("ripemd128", value, 32) }
func RIPEMD160(value string) error { return hexLength("ripemd160", value, 40) }
func TIGER128(value string) error  { return hexLength("tiger128", value, 32) }
func TIGER160(value string) error  { return hexLength("tiger160", value, 40) }
func TIGER192(value string) error  { return hexLength("tiger192", value, 48) }

func SemVer(value string) error {
	if semverPattern.MatchString(value) {
		return nil
	}
	return fail("semver", value, nil, "value must be semantic version 2.0.0")
}

func ULID(value string) error {
	if ulidPattern.MatchString(strings.ToUpper(value)) {
		return nil
	}
	return fail("ulid", value, nil, "value must be a ULID")
}

func CVE(value string) error {
	if cvePattern.MatchString(value) {
		return nil
	}
	return fail("cve", value, nil, "value must be a CVE identifier")
}

func uuidVersion(tag string, value string, version byte, rfc4122 bool) error {
	if UUID(value) == nil && len(value) > 14 && value[14] == version && (!rfc4122 || isRFC4122UUID(value)) {
		return nil
	}
	return fail(tag, value, nil, "value must be the requested UUID version")
}

func isRFC4122UUID(value string) bool {
	if len(value) < 20 {
		return false
	}
	switch value[19] {
	case '8', '9', 'a', 'A', 'b', 'B':
		return true
	default:
		return false
	}
}

func hexLength(tag string, value string, length int) error {
	if len(value) == length && hexPattern.MatchString(value) {
		return nil
	}
	return fail(tag, value, length, "value must be a hex hash of the required length")
}

func validCronField(field string, index int, total int) bool {
	if field == "" {
		return false
	}
	if allDigits(field) {
		return cronFieldInRange(field, index, total)
	}
	for _, r := range field {
		if !unicode.IsDigit(r) && !strings.ContainsRune("*,/-?", r) {
			return false
		}
	}
	return true
}

func cronFieldInRange(field string, index int, total int) bool {
	min, max := cronFieldBounds(index, total)
	n, err := strconv.Atoi(field)
	if err != nil {
		return false
	}
	return n >= min && n <= max
}

func cronFieldBounds(index int, total int) (int, int) {
	switch {
	case total == 6 && index == 0:
		return 0, 59
	case total == 6 && index == 1:
		return 0, 59
	case total == 6 && index == 2:
		return 0, 23
	case total == 6 && index == 3:
		return 1, 31
	case total == 6 && index == 4:
		return 1, 12
	case total == 6 && index == 5:
		return 0, 7
	case total == 5 && index == 0:
		return 0, 59
	case total == 5 && index == 1:
		return 0, 23
	case total == 5 && index == 2:
		return 1, 31
	case total == 5 && index == 3:
		return 1, 12
	case total == 5 && index == 4:
		return 0, 7
	default:
		return 0, 0
	}
}

func digitsOnly(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if r != ' ' && r != '-' {
			return ""
		}
	}
	return b.String()
}

func digitsOnlyValue(value any) string {
	switch v := value.(type) {
	case string:
		return digitsOnly(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		s := strings.TrimSpace(toDecimalString(v))
		if strings.HasPrefix(s, "-") {
			return ""
		}
		return strings.TrimPrefix(s, "+")
	default:
		return ""
	}
}

func toDecimalString(value any) string {
	switch v := value.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uintptr:
		return strconv.FormatUint(uint64(v), 10)
	default:
		return ""
	}
}

func allDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func luhnValid(digits string) bool {
	sum := 0
	double := false
	for i := len(digits) - 1; i >= 0; i-- {
		r := digits[i]
		if r < '0' || r > '9' {
			return false
		}
		n := int(r - '0')
		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
	}
	return sum%10 == 0
}

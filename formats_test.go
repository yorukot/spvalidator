package spvalidator

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"
)

func TestFormatValidators(t *testing.T) {
	t.Run("encodings", func(t *testing.T) {
		expectNoError(t, Base64("SGVsbG8="))
		validateErr(t, Base64("not base64"), "base64")
		expectNoError(t, Base64URL("SGVsbG8="))
		validateErr(t, Base64URL("not base64"), "base64url")
		expectNoError(t, Base64RawURL("SGVsbG8"))
		validateErr(t, Base64RawURL("not base64"), "base64rawurl")
	})

	t.Run("business_ids", func(t *testing.T) {
		expectNoError(t, BIC("DEUTDEFF"))
		validateErr(t, BIC("DEUTDE"), "bic")
		expectNoError(t, BICISO93622014("DEUTDEFF"))
		expectNoError(t, BCP47LanguageTag("en-US"))
		validateErr(t, BCP47LanguageTag("en_US"), "bcp47_language_tag")
		expectNoError(t, BCP47StrictLanguageTag("en-US"))
		validateErr(t, BCP47StrictLanguageTag("en_US"), "bcp47_strict_language_tag")
		expectNoError(t, BTCAddr("1BoatSLRHtKNngkdXEeobR76b53LETtpyT"))
		validateErr(t, BTCAddr("invalid"), "btc_addr")
		expectNoError(t, BTCAddrBech32("bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kygt080"))
		validateErr(t, BTCAddrBech32("bc1bad"), "btc_addr_bech32")
		expectNoError(t, CreditCard("4111 1111 1111 1111"))
		validateErr(t, CreditCard("4111 1111 1111 1112"), "credit_card")
		expectNoError(t, MongoDB("507f1f77bcf86cd799439011"))
		validateErr(t, MongoDB("not-an-object-id"), "mongodb")
		expectNoError(t, MongoDBConnectionString("mongodb://localhost:27017/test"))
		validateErr(t, MongoDBConnectionString("http://example.com"), "mongodb_connection_string")
		expectNoError(t, SpiceDB("document#viewer"))
		validateErr(t, SpiceDB("bad value"), "spicedb")
		expectNoError(t, CVE("CVE-2024-1234"))
		validateErr(t, CVE("CVE-24-1234"), "cve")
		expectNoError(t, ULID("01ARZ3NDEKTSV4RRFFQ69G5FAV"))
		validateErr(t, ULID("not-a-ulid"), "ulid")
		expectNoError(t, SemVer("1.2.3-alpha.1+build.5"))
		validateErr(t, SemVer("1.0.0-01"), "semver")
	})

	t.Run("hashes", func(t *testing.T) {
		cases := []struct {
			name string
			fn   func(string) error
			size int
		}{
			{"md4", MD4, 32},
			{"md5", MD5, 32},
			{"sha256", SHA256, 64},
			{"sha384", SHA384, 96},
			{"sha512", SHA512, 128},
			{"ripemd128", RIPEMD128, 32},
			{"ripemd160", RIPEMD160, 40},
			{"tiger128", TIGER128, 32},
			{"tiger160", TIGER160, 40},
			{"tiger192", TIGER192, 48},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				expectNoError(t, tc.fn(strings.Repeat("0", tc.size)))
				validateErr(t, tc.fn(strings.Repeat("0", tc.size-1)), tc.name)
			})
		}
	})

	t.Run("contact_and_geo", func(t *testing.T) {
		expectNoError(t, E164("+14155552671"))
		validateErr(t, E164("12345"), "e164")
		expectNoError(t, EIN("12-3456789"))
		validateErr(t, EIN("123"), "ein")
		expectNoError(t, Email("user@example.com"))
		validateErr(t, Email("user@"), "email")
		expectNoError(t, ETHAddr("0x0123456789abcdef0123456789abcdef01234567"))
		validateErr(t, ETHAddr("0x123"), "eth_addr")
		expectNoError(t, Hexadecimal("deadBEEF"))
		validateErr(t, Hexadecimal(""), "hexadecimal")
		expectNoError(t, HexColor("#fff"))
		validateErr(t, HexColor("not-a-color"), "hexcolor")
		expectNoError(t, HSL("hsl(120, 50%, 50%)"))
		validateErr(t, HSL("hsl(999, 50%, 50%)"), "hsl")
		expectNoError(t, HSLA("hsla(120, 50%, 50%, 0.5)"))
		validateErr(t, HSLA("hsla(999, 50%, 50%, 0.5)"), "hsla")
		expectNoError(t, RGB("rgb(1, 2, 3)"))
		validateErr(t, RGB("rgb(256, 0, 0)"), "rgb")
		expectNoError(t, RGBA("rgba(1, 2, 3, 0.5)"))
		validateErr(t, RGBA("rgba(1, 2, 3, 2)"), "rgba")
		expectNoError(t, CMYK("cmyk(0%, 0%, 0%, 100%)"))
		validateErr(t, CMYK("cmyk(101%, 0%, 0%, 0%)"), "cmyk")
		expectNoError(t, HTML("<p>hello</p>"))
		validateErr(t, HTML("plain text"), "html")
		expectNoError(t, HTMLEncoded("&lt;p&gt;hello&lt;/p&gt;"))
		validateErr(t, HTMLEncoded("plain text"), "html_encoded")
		expectNoError(t, Latitude(90.0))
		validateErr(t, Latitude(90.1), "latitude")
		expectNoError(t, Longitude(-180.0))
		validateErr(t, Longitude(-180.1), "longitude")
		expectNoError(t, LuhnChecksum("79927398713"))
		validateErr(t, LuhnChecksum("79927398714"), "luhn_checksum")
		validateErr(t, LuhnChecksum(int64(-79927398713)), "luhn_checksum")
		expectNoError(t, ISBN10("0306406152"))
		validateErr(t, ISBN10("0306406153"), "isbn10")
		expectNoError(t, ISBN13("9780306406157"))
		validateErr(t, ISBN13("9780306406158"), "isbn13")
		expectNoError(t, ISBN("9780306406157"))
		validateErr(t, ISBN("not-an-isbn"), "isbn")
		expectNoError(t, ISSN("0317-8471"))
		validateErr(t, ISSN("0317-8472"), "issn")
		expectNoError(t, ISO3166Alpha2("US"))
		validateErr(t, ISO3166Alpha2("Us"), "iso3166_1_alpha2")
		expectNoError(t, ISO3166Alpha3("USA"))
		validateErr(t, ISO3166Alpha3("US"), "iso3166_1_alpha3")
		expectNoError(t, ISO3166AlphaNumeric("840"))
		validateErr(t, ISO3166AlphaNumeric("84"), "iso3166_1_alpha_numeric")
		expectNoError(t, ISO31662("US-CA"))
		validateErr(t, ISO31662("USA"), "iso3166_2")
		expectNoError(t, ISO4217("USD"))
		validateErr(t, ISO4217("usd"), "iso4217")
		expectNoError(t, JSON(`{"a":1}`))
		validateErr(t, JSON(`not json`), "json")
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
		payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"123"}`))
		expectNoError(t, JWT(header+"."+payload+".signature"))
		validateErr(t, JWT("bad.jwt"), "jwt")
		expectNoError(t, PostcodeISO3166Alpha2("12345", "US"))
		validateErr(t, PostcodeISO3166Alpha2("ABCDE", "US"), "postcode_iso3166_alpha2")
		expectNoError(t, PostcodeISO3166Alpha2("H2B 1X9", "CA"))
		expectNoError(t, PostcodeISO3166Alpha2("SW1A 1AA", "GB"))
	})

	t.Run("datetime_and_location", func(t *testing.T) {
		now := time.Now().UTC().Format(time.RFC3339)
		expectNoError(t, DateTime(now))
		validateErr(t, DateTime("not-a-time"), "datetime")
		expectNoError(t, Timezone("UTC"))
		validateErr(t, Timezone("Not/AZone"), "timezone")
	})

	t.Run("uuid_variants", func(t *testing.T) {
		expectNoError(t, UUID("123e4567-e89b-12d3-a456-426614174000"))
		validateErr(t, UUID("not-a-uuid"), "uuid")
		expectNoError(t, UUID3("123e4567-e89b-32d3-a456-426614174000"))
		validateErr(t, UUID3("123e4567-e89b-12d3-a456-426614174000"), "uuid3")
		expectNoError(t, UUID3RFC4122("123e4567-e89b-32d3-a456-426614174000"))
		validateErr(t, UUID3RFC4122("123e4567-e89b-32d3-7456-426614174000"), "uuid3_rfc4122")
		expectNoError(t, UUID4("123e4567-e89b-42d3-a456-426614174000"))
		validateErr(t, UUID4("123e4567-e89b-12d3-a456-426614174000"), "uuid4")
		expectNoError(t, UUID4RFC4122("123e4567-e89b-42d3-a456-426614174000"))
		validateErr(t, UUID4RFC4122("123e4567-e89b-42d3-7456-426614174000"), "uuid4_rfc4122")
		expectNoError(t, UUID5("123e4567-e89b-52d3-a456-426614174000"))
		validateErr(t, UUID5("123e4567-e89b-12d3-a456-426614174000"), "uuid5")
		expectNoError(t, UUID5RFC4122("123e4567-e89b-52d3-a456-426614174000"))
		validateErr(t, UUID5RFC4122("123e4567-e89b-52d3-7456-426614174000"), "uuid5_rfc4122")
		expectNoError(t, UUIDRFC4122("123e4567-e89b-42d3-a456-426614174000"))
		validateErr(t, UUIDRFC4122("123e4567-e89b-42d3-7456-426614174000"), "uuid_rfc4122")
	})
}

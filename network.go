package spvalidator

import (
	"encoding/base64"
	"net"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var urnRFC2141Pattern = regexp.MustCompile(`(?i)^urn:[a-z0-9][a-z0-9-]{0,31}:.+$`)

func CIDR(value string) error {
	if _, _, err := net.ParseCIDR(value); err == nil {
		return nil
	}
	return fail("cidr", value, nil, "must be a CIDR block")
}

func CIDRv4(value string) error {
	ip, network, err := net.ParseCIDR(value)
	if err == nil && ip.To4() != nil {
		ones, bits := network.Mask.Size()
		if bits == 32 && ones >= 0 {
			return nil
		}
	}
	return fail("cidrv4", value, nil, "must be an IPv4 CIDR block")
}

func CIDRv6(value string) error {
	ip, network, err := net.ParseCIDR(value)
	if err == nil && ip.To4() == nil && ip.To16() != nil {
		ones, bits := network.Mask.Size()
		if bits == 128 && ones >= 0 {
			return nil
		}
	}
	return fail("cidrv6", value, nil, "must be an IPv6 CIDR block")
}

func DataURI(value string) error {
	if !strings.HasPrefix(strings.ToLower(value), "data:") {
		return fail("datauri", value, nil, "must be a data URI")
	}
	comma := strings.IndexByte(value, ',')
	if comma < 0 {
		return fail("datauri", value, nil, "data URI must contain a comma")
	}
	meta, data := value[5:comma], value[comma+1:]
	if strings.HasSuffix(strings.ToLower(meta), ";base64") {
		if _, err := base64.StdEncoding.DecodeString(data); err != nil {
			return fail("datauri", value, nil, "base64 data URI payload is invalid")
		}
		return nil
	}
	if _, err := url.PathUnescape(data); err == nil {
		return nil
	}
	return fail("datauri", value, nil, "data URI payload is invalid")
}

func Domain(value string) error {
	host := strings.TrimSuffix(value, ".")
	if strings.Count(host, ".") >= 1 && net.ParseIP(host) == nil && validHostname(host, true) && topLevelDomainHasLetter(host) {
		return nil
	}
	return fail("domain", value, nil, "must be a domain name")
}

func FQDN(value string) error {
	host := strings.TrimSuffix(value, ".")
	if strings.Count(host, ".") < 1 || !validHostname(host, true) {
		return fail("fqdn", value, nil, "must be a fully qualified domain name")
	}
	return nil
}

func Hostname(value string) error {
	if validHostname(value, false) && startsWithLetter(value) {
		return nil
	}
	return fail("hostname", value, nil, "must be an RFC 952 hostname")
}

func HostnameRFC1123(value string) error {
	if validHostname(value, false) {
		return nil
	}
	return fail("hostname_rfc1123", value, nil, "must be an RFC 1123 hostname")
}

func HostnamePort(value string) error {
	host, port, err := net.SplitHostPort(value)
	if err == nil && validHostForAddress(host) && Port(port) == nil {
		return nil
	}
	return fail("hostname_port", value, nil, "must be host:port")
}

func Port(value any) error {
	var port int64
	switch v := value.(type) {
	case string:
		if v == "" {
			return fail("port", value, nil, "port must be between 1 and 65535")
		}
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return fail("port", value, nil, "port must be numeric")
		}
		port = n
	default:
		n, ok := numberValue(value)
		if !ok {
			return fail("port", value, nil, "port must be numeric")
		}
		port = int64(n)
		if float64(port) != n {
			return fail("port", value, nil, "port must be an integer")
		}
	}
	if port >= 1 && port <= 65535 {
		return nil
	}
	return fail("port", value, nil, "port must be between 1 and 65535")
}

func IP(value string) error {
	if net.ParseIP(value) != nil {
		return nil
	}
	return fail("ip", value, nil, "must be an IP address")
}

func IPAddr(value string) error { return IP(value) }

func IP4Addr(value string) error { return IPv4(value) }

func IP6Addr(value string) error { return IPv6(value) }

func IPv4(value string) error {
	ip := net.ParseIP(value)
	if ip != nil && ip.To4() != nil {
		return nil
	}
	return fail("ipv4", value, nil, "must be an IPv4 address")
}

func IPv6(value string) error {
	ip := net.ParseIP(value)
	if ip != nil && ip.To4() == nil && ip.To16() != nil {
		return nil
	}
	return fail("ipv6", value, nil, "must be an IPv6 address")
}

func MAC(value string) error {
	if _, err := net.ParseMAC(value); err == nil {
		return nil
	}
	return fail("mac", value, nil, "must be a MAC address")
}

func TCPAddr(value string) error {
	return networkAddr("tcp_addr", value, false, false)
}

func TCP4Addr(value string) error {
	return networkAddr("tcp4_addr", value, true, false)
}

func TCP6Addr(value string) error {
	return networkAddr("tcp6_addr", value, false, true)
}

func UDPAddr(value string) error {
	return networkAddr("udp_addr", value, false, false)
}

func UDP4Addr(value string) error {
	return networkAddr("udp4_addr", value, true, false)
}

func UDP6Addr(value string) error {
	return networkAddr("udp6_addr", value, false, true)
}

func UnixAddr(value string) error {
	if strings.TrimSpace(value) != "" && !strings.ContainsRune(value, 0) {
		return nil
	}
	return fail("unix_addr", value, nil, "must be a Unix socket path")
}

func UDSExists(value string) error {
	if runtime.GOOS == "linux" && (strings.HasPrefix(value, "@") || strings.HasPrefix(value, "\x00")) {
		return nil
	}
	info, err := os.Stat(value)
	if err == nil && info.Mode()&os.ModeSocket != 0 {
		return nil
	}
	return fail("uds_exists", value, nil, "Unix domain socket must exist")
}

func URI(value string) error {
	parsed, err := url.Parse(value)
	if err == nil && parsed.Scheme != "" {
		return nil
	}
	return fail("uri", value, nil, "must be a URI")
}

func URL(value string) error {
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		parsed, err = url.Parse(value)
	}
	if err == nil && parsed.Scheme != "" && parsed.Host != "" {
		return nil
	}
	return fail("url", value, nil, "must be an absolute URL")
}

func HTTPURL(value string) error {
	parsed, err := url.Parse(value)
	if err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != "" {
		return nil
	}
	return fail("http_url", value, nil, "must be an HTTP or HTTPS URL")
}

func HTTPSURL(value string) error {
	parsed, err := url.Parse(value)
	if err == nil && parsed.Scheme == "https" && parsed.Host != "" {
		return nil
	}
	return fail("https_url", value, nil, "must be an HTTPS URL")
}

func Origin(value string) error {
	parsed, err := url.Parse(value)
	if err == nil &&
		(parsed.Scheme == "http" || parsed.Scheme == "https") &&
		parsed.Host != "" &&
		(parsed.Path == "" || parsed.Path == "/") &&
		parsed.RawQuery == "" &&
		parsed.Fragment == "" &&
		parsed.User == nil {
		return nil
	}
	return fail("origin", value, nil, "must be a web origin")
}

func URLEncoded(value string) error {
	if value == "" {
		return fail("url_encoded", value, nil, "must be URL encoded")
	}
	decoded, err := url.QueryUnescape(value)
	if err == nil && decoded != value {
		return nil
	}
	return fail("url_encoded", value, nil, "must be URL encoded")
}

func URNRFC2141(value string) error {
	if urnRFC2141Pattern.MatchString(value) {
		return nil
	}
	return fail("urn_rfc2141", value, nil, "must be an RFC 2141 URN")
}

func networkAddr(tag string, value string, requireV4 bool, requireV6 bool) error {
	host, port, err := net.SplitHostPort(value)
	if err != nil || Port(port) != nil {
		return fail(tag, value, nil, "must be host:port")
	}
	if requireV4 {
		if IPv4(host) == nil {
			return nil
		}
		return fail(tag, value, nil, "must contain an IPv4 host")
	}
	if requireV6 {
		if IPv6(host) == nil {
			return nil
		}
		return fail(tag, value, nil, "must contain an IPv6 host")
	}
	if validHostForAddress(host) {
		return nil
	}
	return fail(tag, value, nil, "must contain a valid host")
}

func validHostForAddress(host string) bool {
	host = strings.Trim(host, "[]")
	return net.ParseIP(host) != nil || validHostname(host, true)
}

func validHostname(value string, allowFQDN bool) bool {
	if value == "" || len(value) > 253 || strings.Contains(value, "_") {
		return false
	}
	value = strings.TrimSuffix(value, ".")
	if strings.Contains(value, "..") {
		return false
	}
	labelCount := strings.Count(value, ".") + 1
	if labelCount > 1 && !allowFQDN {
		return false
	}
	for label := range strings.SplitSeq(value, ".") {
		if label == "" || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, r := range label {
			if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return false
			}
		}
	}
	return true
}

func topLevelDomainHasLetter(value string) bool {
	tld := value
	if lastDot := strings.LastIndexByte(value, '.'); lastDot >= 0 {
		tld = value[lastDot+1:]
	}
	for i := 0; i < len(tld); i++ {
		c := tld[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			return true
		}
	}
	return false
}

func startsWithLetter(value string) bool {
	value = strings.TrimSuffix(value, ".")
	if value == "" {
		return false
	}
	first := value[0]
	return (first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z')
}

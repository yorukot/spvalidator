package spvalidator

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNetworkValidators(t *testing.T) {
	t.Run("cidr", func(t *testing.T) {
		expectNoError(t, CIDR("192.168.0.0/24"))
		validateErr(t, CIDR("not-a-cidr"), "cidr")
		expectNoError(t, CIDRv4("192.168.0.0/24"))
		validateErr(t, CIDRv4("2001:db8::/32"), "cidrv4")
		expectNoError(t, CIDRv6("2001:db8::/32"))
		validateErr(t, CIDRv6("192.168.0.0/24"), "cidrv6")
	})

	t.Run("data_uri", func(t *testing.T) {
		expectNoError(t, DataURI("data:text/plain;base64,SGVsbG8="))
		expectNoError(t, DataURI("data:text/plain,hello%20world"))
		validateErr(t, DataURI("http:text/plain,hello"), "datauri")
	})

	t.Run("hostnames", func(t *testing.T) {
		expectNoError(t, FQDN("example.com"))
		validateErr(t, FQDN("example"), "fqdn")
		expectNoError(t, Hostname("example"))
		validateErr(t, Hostname("1example"), "hostname")
		expectNoError(t, HostnameRFC1123("1example"))
		validateErr(t, HostnameRFC1123("bad_host"), "hostname_rfc1123")
		expectNoError(t, HostnamePort("example.com:443"))
		validateErr(t, HostnamePort("example.com:notaport"), "hostname_port")
	})

	t.Run("ip_and_port", func(t *testing.T) {
		expectNoError(t, IP("127.0.0.1"))
		validateErr(t, IP("999.0.0.1"), "ip")
		expectNoError(t, IPAddr("127.0.0.1"))
		expectNoError(t, IP4Addr("127.0.0.1"))
		expectNoError(t, IP6Addr("::1"))
		expectNoError(t, IPv4("127.0.0.1"))
		validateErr(t, IPv4("::1"), "ipv4")
		expectNoError(t, IPv6("::1"))
		validateErr(t, IPv6("127.0.0.1"), "ipv6")
		expectNoError(t, MAC("aa:bb:cc:dd:ee:ff"))
		validateErr(t, MAC("invalid"), "mac")
		expectNoError(t, Port(443))
		expectNoError(t, Port("443"))
		validateErr(t, Port(0), "port")
		validateErr(t, Port("70000"), "port")
	})

	t.Run("socket_and_urls", func(t *testing.T) {
		expectNoError(t, URI("mailto:user@example.com"))
		validateErr(t, URI("://broken"), "uri")
		expectNoError(t, URL("https://example.com"))
		validateErr(t, URL("/relative"), "url")
		expectNoError(t, HTTPURL("http://example.com"))
		expectNoError(t, HTTPURL("https://example.com"))
		validateErr(t, HTTPURL("ftp://example.com"), "http_url")
		expectNoError(t, HTTPSURL("https://example.com"))
		validateErr(t, HTTPSURL("http://example.com"), "https_url")
		expectNoError(t, Origin("https://example.com"))
		validateErr(t, Origin("https://example.com/path"), "origin")
		expectNoError(t, URLEncoded("hello%20world"))
		validateErr(t, URLEncoded("hello world"), "url_encoded")
		expectNoError(t, URNRFC2141("urn:example:animal:ferret:nose"))
		validateErr(t, URNRFC2141("urn::bad"), "urn_rfc2141")
	})

	t.Run("tcp_udp", func(t *testing.T) {
		expectNoError(t, TCPAddr("127.0.0.1:80"))
		expectNoError(t, TCP4Addr("127.0.0.1:80"))
		validateErr(t, TCP4Addr("[::1]:80"), "tcp4_addr")
		expectNoError(t, TCP6Addr("[::1]:80"))
		validateErr(t, TCP6Addr("127.0.0.1:80"), "tcp6_addr")
		expectNoError(t, UDPAddr("127.0.0.1:80"))
		expectNoError(t, UDP4Addr("127.0.0.1:80"))
		validateErr(t, UDP4Addr("[::1]:80"), "udp4_addr")
		expectNoError(t, UDP6Addr("[::1]:80"))
		validateErr(t, UDP6Addr("127.0.0.1:80"), "udp6_addr")
		expectNoError(t, UnixAddr("/tmp/socket"))
		validateErr(t, UnixAddr(""), "unix_addr")
	})

	t.Run("uds_exists", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("unix sockets are not supported on windows")
		}
		socketPath := filepath.Join("/tmp", "spvalidator-test.sock")
		_ = os.Remove(socketPath)
		ln, err := net.Listen("unix", socketPath)
		if err != nil {
			t.Fatalf("listen unix socket: %v", err)
		}
		defer ln.Close()
		defer os.Remove(socketPath)
		expectNoError(t, UDSExists(socketPath))
		validateErr(t, UDSExists(filepath.Join("/tmp", "spvalidator-missing.sock")), "uds_exists")
	})
}

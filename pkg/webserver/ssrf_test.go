package webserver

import (
	"net"
	"testing"
)

func TestIsPrivateOrReserved(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"127.0.0.1", true},             // loopback
		{"10.0.0.5", true},              // RFC1918
		{"192.168.1.1", true},           // RFC1918
		{"172.16.0.1", true},            // RFC1918
		{"169.254.169.254", true},       // link-local (cloud metadata)
		{"100.64.0.1", true},            // CGNAT RFC6598
		{"::1", true},                   // IPv6 loopback
		{"fe80::1", true},               // IPv6 link-local
		{"fc00::1", true},               // IPv6 ULA
		{"0.0.0.0", true},               // unspecified
		{"8.8.8.8", false},              // public
		{"1.1.1.1", false},              // public
		{"2606:4700:4700::1111", false}, // public IPv6
	}
	for _, tt := range tests {
		ip := net.ParseIP(tt.ip)
		if ip == nil {
			t.Fatalf("bad test IP %q", tt.ip)
		}
		if got := isPrivateOrReserved(ip); got != tt.want {
			t.Errorf("isPrivateOrReserved(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestIpBlocked_DefaultBlocksPrivate(t *testing.T) {
	s := &Server{} // allowPrivateTargets defaults to false
	if !s.ipBlocked(net.ParseIP("127.0.0.1")) {
		t.Error("expected 127.0.0.1 blocked by default")
	}
	if s.ipBlocked(net.ParseIP("8.8.8.8")) {
		t.Error("expected 8.8.8.8 allowed")
	}

	s.allowPrivateTargets = true
	if s.ipBlocked(net.ParseIP("127.0.0.1")) {
		t.Error("expected 127.0.0.1 allowed when allowPrivateTargets=true")
	}
}

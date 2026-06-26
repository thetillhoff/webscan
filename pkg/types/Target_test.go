package types

import "testing"

func TestNewTarget(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantErr  bool
		wantHost string
		wantType TargetType
	}{
		{"plain domain", "Example.COM", false, "example.com", Domain},
		{"https url", "https://example.com/Path", false, "example.com", Domain},
		{"http with port", "http://example.com:8080", false, "example.com", Domain},
		{"ipv4", "192.168.0.1", false, "192.168.0.1", Ipv4},
		{"ipv6 bracketed with port", "[::1]:443", false, "::1", Ipv6},
		{"unsupported scheme", "ftp://example.com", true, "", None},
		{"empty", "   ", true, "", None},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := NewTarget(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewTarget(%q) expected error, got nil", tt.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewTarget(%q) unexpected error: %v", tt.in, err)
			}
			if target.Hostname() != tt.wantHost {
				t.Errorf("Hostname() = %q, want %q", target.Hostname(), tt.wantHost)
			}
			if target.TargetType() != tt.wantType {
				t.Errorf("TargetType() = %v, want %v", target.TargetType(), tt.wantType)
			}
		})
	}
}

// Path case must survive parsing (only the host is case-insensitive).
func TestNewTargetPreservesPathCase(t *testing.T) {
	target, err := NewTarget("https://example.com/MixedCase/Path")
	if err != nil {
		t.Fatal(err)
	}
	if got := target.Path(); got != "/MixedCase/Path" {
		t.Errorf("Path() = %q, want %q", got, "/MixedCase/Path")
	}
}

package dnsScan

import "testing"

func TestCheckSPF(t *testing.T) {
	tests := []struct {
		name    string
		records []string
		want    string
	}{
		{"no record", []string{"some other txt"}, "Hint: No SPF record detected."},
		{"valid ip4/ip6/all", []string{"v=spf1 ip4:192.0.2.0/24 ip6:2001:db8::/32 -all"}, ""},
		{"valid single ip4", []string{"v=spf1 ip4:192.0.2.1 -all"}, ""},
		{"multiple records", []string{"v=spf1 -all", "v=spf1 ip4:192.0.2.1 -all"}, "Hint: Multiple SPF records detected."},
		{"ptr discouraged", []string{"v=spf1 ptr -all"}, "Hint: PTR records should not be used inSPF records as they are slow and inefficient."},
		{"invalid ip4", []string{"v=spf1 ip4:not-an-ip -all"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSPF(tt.records); got != tt.want {
				t.Errorf("CheckSPF(%v) = %q, want %q", tt.records, got, tt.want)
			}
		})
	}
}

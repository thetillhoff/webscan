package portScan

import "testing"

func TestFullPortRange(t *testing.T) {
	ports := fullPortRange()

	if len(ports) != 65535 {
		t.Fatalf("expected 65535 ports, got %d", len(ports))
	}
	if ports[0] != 1 {
		t.Fatalf("expected first port to be 1, got %d", ports[0])
	}
	if ports[len(ports)-1] != 65535 {
		t.Fatalf("expected last port to be 65535, got %d", ports[len(ports)-1])
	}
}

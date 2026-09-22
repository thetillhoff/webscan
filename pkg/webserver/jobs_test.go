package webserver

import (
	"strings"
	"testing"
)

func TestSynchronizedBuffer_UnboundedByDefault(t *testing.T) {
	b := &synchronizedBuffer{}
	for range 1000 {
		_, _ = b.Write([]byte("0123456789"))
	}
	if got := len(b.String()); got != 10000 {
		t.Fatalf("expected 10000 bytes retained, got %d", got)
	}
}

func TestSynchronizedBuffer_TrimsToMaxSize(t *testing.T) {
	b := &synchronizedBuffer{maxSize: 50}
	for i := range 1000 {
		_, _ = b.Write([]byte{byte('a' + i%26)})
	}
	got := b.String()
	if len(got) != 50 {
		t.Fatalf("expected buffer trimmed to 50 bytes, got %d", len(got))
	}
	// The most recent byte written must be the last character kept.
	if !strings.HasSuffix(got, string(byte('a'+999%26))) {
		t.Fatalf("expected trimmed buffer to keep the most recent bytes, got %q", got)
	}
}

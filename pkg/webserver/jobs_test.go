package webserver

import (
	"context"
	"strings"
	"testing"
	"time"
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

func TestWatchForStale_CancelsAfterNoProgress(t *testing.T) {
	s := &Server{staleTimeout: 50 * time.Millisecond}
	buf := &synchronizedBuffer{lastWrite: time.Now()}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	defer close(done)
	staleFired := make(chan struct{})

	go s.watchForStale(ctx, cancel, buf, done, staleFired)

	select {
	case <-ctx.Done():
		// expected: no writes ever happen, so the watchdog cancels.
		select {
		case <-staleFired:
		default:
			t.Fatal("expected staleFired to be closed when cancelling due to staleness")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("expected context to be cancelled after no status progress")
	}
}

func TestWatchForStale_DoesNotCancelWhileProgressing(t *testing.T) {
	s := &Server{staleTimeout: 100 * time.Millisecond}
	buf := &synchronizedBuffer{lastWrite: time.Now()}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	staleFired := make(chan struct{})

	go s.watchForStale(ctx, cancel, buf, done, staleFired)

	deadline := time.Now().Add(700 * time.Millisecond)
	for time.Now().Before(deadline) {
		_, _ = buf.Write([]byte("x"))
		time.Sleep(30 * time.Millisecond)
	}
	close(done)

	select {
	case <-ctx.Done():
		t.Fatal("context should not be cancelled while status keeps updating")
	default:
	}
}

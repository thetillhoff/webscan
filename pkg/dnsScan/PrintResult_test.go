package dnsScan

import (
	"strings"
	"testing"
)

func TestPrintResultEmpty(t *testing.T) {
	var out strings.Builder
	PrintResult(Result{}, &out)

	got := out.String()
	if !strings.Contains(got, "No DNS records found") {
		t.Errorf("empty result should report no records, got:\n%s", got)
	}
	if strings.Contains(got, "DNS records:") {
		t.Errorf("empty result should not print the records header, got:\n%s", got)
	}
}

func TestPrintResultWithRecords(t *testing.T) {
	var out strings.Builder
	PrintResult(Result{ARecords: []string{"192.0.2.1"}}, &out)

	got := out.String()
	if !strings.Contains(got, "DNS records:") {
		t.Errorf("result with records should print the header, got:\n%s", got)
	}
	if strings.Contains(got, "No DNS records found") {
		t.Errorf("result with records should not report empty, got:\n%s", got)
	}
	if !strings.Contains(got, "192.0.2.1") {
		t.Errorf("A record should be printed, got:\n%s", got)
	}
}

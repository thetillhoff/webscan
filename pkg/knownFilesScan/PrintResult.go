package knownFilesScan

import (
	"fmt"
	"io"
	"log/slog"
	"sort"
)

func PrintResult(result Result, schemaLabel string, out io.Writer) {
	var expectedPresent []FileResult
	var expectedMissing []FileResult
	var sensitiveExposed []FileResult

	sort.Slice(result.files, func(i, j int) bool {
		return result.files[i].Path < result.files[j].Path
	})

	for _, f := range result.files {
		switch {
		case f.Category == categorySensitive && f.Found:
			sensitiveExposed = append(sensitiveExposed, f)
		case f.Category == categoryExpected && f.Found:
			expectedPresent = append(expectedPresent, f)
		case f.Category == categoryExpected && !f.Found:
			expectedMissing = append(expectedMissing, f)
		}
	}

	if len(expectedPresent) == 0 && len(expectedMissing) == 0 && len(sensitiveExposed) == 0 {
		return
	}

	if _, err := fmt.Fprintf(out, "\n\n## Well-known files scan results (%s)\n\n", schemaLabel); err != nil {
		slog.Debug("knownFilesScan: Error writing to output", "error", err)
	}

	if len(sensitiveExposed) > 0 {
		if _, err := fmt.Fprintf(out, "WARNING: Sensitive files publicly accessible:\n"); err != nil {
			slog.Debug("knownFilesScan: Error writing to output", "error", err)
		}
		for _, f := range sensitiveExposed {
			if _, err := fmt.Fprintf(out, "  - %s (%s) — should not be publicly accessible\n", f.Path, f.Label); err != nil {
				slog.Debug("knownFilesScan: Error writing to output", "error", err)
			}
		}
		if _, err := fmt.Fprintln(out); err != nil {
			slog.Debug("knownFilesScan: Error writing to output", "error", err)
		}
	}

	for _, f := range expectedPresent {
		if _, err := fmt.Fprintf(out, "- %s found\n", f.Label); err != nil {
			slog.Debug("knownFilesScan: Error writing to output", "error", err)
		}
		printFileDetails(f, out)
	}

	for _, f := range expectedMissing {
		if _, err := fmt.Fprintf(out, "- %s not found (consider adding %s)\n", f.Label, f.Path); err != nil {
			slog.Debug("knownFilesScan: Error writing to output", "error", err)
		}
	}
}

func printFileDetails(f FileResult, out io.Writer) {
	hasObs := len(f.Observations) > 0
	hasRec := len(f.Recommendations) > 0

	if !hasObs && !hasRec {
		return
	}

	if hasObs && hasRec {
		if _, err := fmt.Fprintf(out, "  Observations:\n"); err != nil {
			slog.Debug("knownFilesScan: Error writing to output", "error", err)
		}
		for _, obs := range f.Observations {
			if _, err := fmt.Fprintf(out, "  - %s\n", obs); err != nil {
				slog.Debug("knownFilesScan: Error writing to output", "error", err)
			}
		}
		if _, err := fmt.Fprintf(out, "  Recommendations:\n"); err != nil {
			slog.Debug("knownFilesScan: Error writing to output", "error", err)
		}
		for _, rec := range f.Recommendations {
			if _, err := fmt.Fprintf(out, "  - %s\n", rec); err != nil {
				slog.Debug("knownFilesScan: Error writing to output", "error", err)
			}
		}
	} else if hasObs {
		for _, obs := range f.Observations {
			if _, err := fmt.Fprintf(out, "  - %s\n", obs); err != nil {
				slog.Debug("knownFilesScan: Error writing to output", "error", err)
			}
		}
	} else {
		if _, err := fmt.Fprintf(out, "  Recommendations:\n"); err != nil {
			slog.Debug("knownFilesScan: Error writing to output", "error", err)
		}
		for _, rec := range f.Recommendations {
			if _, err := fmt.Fprintf(out, "  - %s\n", rec); err != nil {
				slog.Debug("knownFilesScan: Error writing to output", "error", err)
			}
		}
	}
}

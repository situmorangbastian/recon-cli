package service

import (
	"testing"
	"time"
)

func parseDate(t *testing.T, dateStr string) time.Time {
	t.Helper()
	tm, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t.Fatalf("failed to parse date %q: %v", dateStr, err)
	}
	return tm
}

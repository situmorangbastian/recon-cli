package service

import (
	"testing"
	"time"

	"github.com/situmorangbastian/recon-cli/internal/reader"
)

func TestFetchBankStatements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
		startDate time.Time
		endDate   time.Time
	}{
		{
			name:      "success_with_empty_filter",
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "success_with_filter",
			wantCount: 3,
			wantErr:   false,
			startDate: parseDate(t, "2025-08-01"),
			endDate:   parseDate(t, "2025-08-31"),
		},
		{
			name:      "success_with_filter_with_empty_result",
			wantCount: 0,
			wantErr:   false,
			startDate: parseDate(t, "2025-07-01"),
			endDate:   parseDate(t, "2025-07-31"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := reader.NewReader(nil, []string{"2006-01-02"}, []string{"../../testdata/bank_a_statement.csv"}, "")
			service := NewService(reader)
			txns, err := service.FetchBankStatements(tt.startDate, tt.endDate)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FetchBankStatements() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && len(txns) != tt.wantCount {
				t.Errorf("FetchBankStatements() got %d transactions, want %d", len(txns), tt.wantCount)
			}
		})
	}
}

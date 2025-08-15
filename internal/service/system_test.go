package service

import (
	"testing"
	"time"

	"github.com/situmorangbastian/recon-cli/internal/reader"
)

func TestFetchSystemTransactions(t *testing.T) {
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
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "success_with_filter",
			wantCount: 5,
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

			reader := reader.NewReader([]string{"2006-01-02 15:04:05"}, nil, nil, "../../testdata/system_transactions.csv")
			service := NewService(reader)
			txns, err := service.FetchSystemTransactions(tt.startDate, tt.endDate)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FetchSystemTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && len(txns) != tt.wantCount {
				t.Errorf("FetchSystemTransactions() got %d transactions, want %d", len(txns), tt.wantCount)
			}
		})
	}
}

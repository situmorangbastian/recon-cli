package reconcile

import (
	"testing"

	"github.com/situmorangbastian/recon-cli/internal/reader"
	"github.com/situmorangbastian/recon-cli/internal/service"
)

func TestReconcile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		wantResult ReconcileSummary
		wantErr    bool
		startDate  string
		endDate    string
	}{
		{
			name: "success_reconcile",
			wantResult: ReconcileSummary{
				TotalTransactionsProcessed: 10,
			},
			wantErr:   false,
			startDate: "2025-08-01",
			endDate:   "2025-08-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := reader.NewReader(
				[]string{"2006-01-02 15:04:05"}, []string{"2006-01-02"},
				[]string{
					"../../testdata/bank_a_statement.csv",
					"../../testdata/bank_b_statement.csv",
				},
				"../../testdata/system_transactions.csv")
			svc := service.NewService(reader)
			reconcile := New(svc)
			result, err := reconcile.Reconcile(tt.startDate, tt.endDate)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if result.TotalTransactionsProcessed != tt.wantResult.TotalTransactionsProcessed {
				t.Errorf("Reconcile() got %d transactions, want %d", result.TotalTransactionsProcessed, tt.wantResult.TotalTransactionsProcessed)
			}
		})
	}
}

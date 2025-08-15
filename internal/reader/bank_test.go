package reader

import (
	"testing"
)

func TestReadBankStmtsCSV(t *testing.T) {
	t.Parallel()

	bankStmtDateLayout := []string{"2006-01-02"}
	reader := NewReader(nil, bankStmtDateLayout)

	tests := []struct {
		name      string
		filePaths []string
		wantCount int
		wantErr   bool
	}{
		{
			name: "success_valid_file",
			filePaths: []string{
				"../../testdata/bank_a_statement.csv",
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "error_missing_file",
			filePaths: []string{"missingfile.csv"},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			txns, err := reader.ReadBankStmtsCSV(tt.filePaths)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadBankStmtsCSV() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && len(txns) != tt.wantCount {
				t.Errorf("ReadBankStmtsCSV() got %d transactions, want %d", len(txns), tt.wantCount)
			}
		})
	}
}

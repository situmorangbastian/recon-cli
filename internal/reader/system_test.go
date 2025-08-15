package reader

import (
	"testing"
)

func TestReadSysTxnsCSV(t *testing.T) {
	t.Parallel()

	sysTxnTimeLayout := []string{"2006-01-02 15:04:05"}

	tests := []struct {
		name      string
		filePath  string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "success_valid_file",
			filePath:  "../../testdata/system_transactions.csv",
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "error_missing_file",
			filePath:  "missingfile.csv",
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := NewReader(sysTxnTimeLayout, nil, nil, tt.filePath)
			txns, err := reader.ReadSysTxnsCSV()
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadSysTxnsCSV() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && len(txns) != tt.wantCount {
				t.Errorf("ReadSysTxnsCSV() got %d transactions, want %d", len(txns), tt.wantCount)
			}
		})
	}
}

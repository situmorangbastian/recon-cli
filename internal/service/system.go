package service

import (
	"fmt"
	"time"

	"github.com/situmorangbastian/recon-cli/internal"
)

func (s *Service) FetchSystemTransactions(startDate, endDate time.Time) ([]internal.Transaction, error) {
	transactions, err := s.reader.ReadSysTxnsCSV()
	if err != nil {
		return nil, fmt.Errorf("FetchSystemTransactions: failed get transactions: %w", err)
	}
	if startDate.IsZero() || endDate.IsZero() {
		return transactions, nil
	}

	var filteredTxns []internal.Transaction
	for _, tx := range transactions {
		txDate := tx.TransactionTime.Truncate(24 * time.Hour)
		if (txDate.Equal(startDate) || txDate.After(startDate)) &&
			(txDate.Equal(endDate) || txDate.Before(endDate)) {
			filteredTxns = append(filteredTxns, tx)
		}
	}

	return filteredTxns, nil
}

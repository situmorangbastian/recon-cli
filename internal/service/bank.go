package service

import (
	"fmt"
	"time"

	"github.com/situmorangbastian/recon-cli/internal"
)

func (s *Service) FetchBankStatements(startDate, endDate time.Time) ([]internal.BankStatement, error) {
	bankStmts, err := s.reader.ReadBankStmtsCSV()
	if err != nil {
		return nil, fmt.Errorf("FetchBankStatements: failed get bank statements: %w", err)
	}

	if startDate.IsZero() || endDate.IsZero() {
		return bankStmts, nil
	}

	var filteredBankStmts []internal.BankStatement
	for _, stmt := range bankStmts {
		stmtDate := stmt.Date.Truncate(24 * time.Hour)
		if (stmtDate.Equal(startDate) || stmtDate.After(startDate)) &&
			(stmtDate.Equal(endDate) || stmtDate.Before(endDate)) {
			filteredBankStmts = append(filteredBankStmts, stmt)
		}
	}

	return filteredBankStmts, nil
}

package reconcile

import (
	"fmt"
	"math"
	"time"

	"github.com/situmorangbastian/recon-cli/internal"
	"github.com/situmorangbastian/recon-cli/internal/service"
)

type Reconcile struct {
	svc *service.Service
}

func New(svc *service.Service) *Reconcile {
	return &Reconcile{
		svc: svc,
	}
}

func (r *Reconcile) Reconcile(startDateStr, endDateStr string) (*ReconcileSummary, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("Reconcile: failed parse start date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("Reconcile: failed parse end date: %w", err)
	}

	systemTxns, err := r.svc.FetchSystemTransactions(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("Reconcile: failed get system transactions: %w", err)
	}

	bankStmts, err := r.svc.FetchBankStatements(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("Reconcile: failed get bank statments: %w", err)
	}

	return r.performMatching(systemTxns, bankStmts), nil
}

func (r *Reconcile) performMatching(systemTxns []internal.Transaction, bankStatements []internal.BankStatement) *ReconcileSummary {
	summary := &ReconcileSummary{
		UnmatchedBankTx: make(map[string][]internal.BankStatement),
	}

	unmatchedSystemTxns := make([]internal.Transaction, len(systemTxns))
	copy(unmatchedSystemTxns, systemTxns)

	unmatchedBankStmts := make([]internal.BankStatement, len(bankStatements))
	copy(unmatchedBankStmts, bankStatements)

	var discrepantTxs []DiscrepantTransaction

	for i := len(unmatchedSystemTxns) - 1; i >= 0; i-- {
		systemTxn := unmatchedSystemTxns[i]

		for j := len(unmatchedBankStmts) - 1; j >= 0; j-- {
			bankStmt := unmatchedBankStmts[j]

			if r.isMatch(systemTxn, bankStmt) {
				summary.TotalMatchedTransactions++

				expectedAmount := r.getExpectedBankAmount(systemTxn)
				if math.Abs(expectedAmount-bankStmt.Amount) > 1 {
					difference := math.Abs(expectedAmount - bankStmt.Amount)
					summary.TotalDiscrepancies += difference

					discrepantTxs = append(discrepantTxs, DiscrepantTransaction{
						SystemTxn:     systemTxn,
						BankStatement: &bankStmt,
						Difference:    difference,
					})
				}

				unmatchedSystemTxns = append(unmatchedSystemTxns[:i], unmatchedSystemTxns[i+1:]...)
				unmatchedBankStmts = append(unmatchedBankStmts[:j], unmatchedBankStmts[j+1:]...)
				break
			}
		}
	}

	for _, stmt := range unmatchedBankStmts {
		summary.UnmatchedBankTx[stmt.File] = append(summary.UnmatchedBankTx[stmt.File], stmt)
	}

	summary.TotalTransactionsProcessed = len(systemTxns) + len(bankStatements)
	summary.TotalUnmatchedTransactions = len(unmatchedSystemTxns) + len(unmatchedBankStmts)
	summary.UnmatchedSystemTxn = unmatchedSystemTxns
	summary.DiscrepantTransactions = discrepantTxs
	return summary
}

func (r *Reconcile) isMatch(systemTxn internal.Transaction, bankStmt internal.BankStatement) bool {
	systemDate := systemTxn.TransactionTime.Truncate(24 * time.Hour)
	bankDate := bankStmt.Date.Truncate(24 * time.Hour)

	if !systemDate.Equal(bankDate) {
		return false
	}

	expectedAmount := r.getExpectedBankAmount(systemTxn)
	return math.Abs(expectedAmount-bankStmt.Amount) <= 1000
}

func (r *Reconcile) getExpectedBankAmount(systemTxn internal.Transaction) float64 {
	if systemTxn.Type == internal.DEBIT {
		return -systemTxn.Amount
	}
	return systemTxn.Amount
}

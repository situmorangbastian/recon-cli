package reconcile

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
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

	reconciled := r.performMatching(systemTxns, bankStmts)
	reconciled.StartDate = startDateStr
	reconciled.EndDate = endDateStr
	return reconciled, nil
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

			if systemTxn.BankRefNo == bankStmt.UniqueIdentifier {
				summary.TotalMatchedTransactions++
				expectedAmount := r.getExpectedBankAmount(systemTxn)
				if math.Abs(expectedAmount-bankStmt.Amount) > 0 {
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

func (r *Reconcile) getExpectedBankAmount(systemTxn internal.Transaction) float64 {
	if systemTxn.Type == internal.DEBIT {
		return -systemTxn.Amount
	}
	return systemTxn.Amount
}

func (r *Reconcile) WriteCSVReport(summary ReconcileSummary, output string) error {
	folder := filepath.Clean(output)
	dir := filepath.Dir(folder + "/")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("WriteCSVReport: failed to create directory %s: %w", dir, err)
	}

	file, err := os.Create(fmt.Sprintf("%s/%s_", output, fmt.Sprintf("recon_summary_%s_%s.csv", summary.StartDate, summary.EndDate)))
	if err != nil {
		return fmt.Errorf("WriteCSVReport: failed create file: %w", err)
	}
	defer file.Close()

	data := [][]string{
		{"Processed", "Matched", "Unmatched", "Discrepancy"},
	}

	data = append(data, []string{
		strconv.Itoa(summary.TotalTransactionsProcessed),
		strconv.Itoa(summary.TotalMatchedTransactions),
		strconv.Itoa(summary.TotalUnmatchedTransactions),
		strconv.FormatFloat(summary.TotalDiscrepancies, 'f', 2, 64),
	})

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(data); err != nil {
		return fmt.Errorf("WriteCSVReport: failed write file: %w", err)
	}

	file, err = os.Create(fmt.Sprintf("%s/%s_", output, fmt.Sprintf("sys_txns_missing_bank_stmts_%s_%s.csv", summary.StartDate, summary.EndDate)))
	if err != nil {
		return fmt.Errorf("WriteCSVReport: failed create file: %w", err)
	}
	defer file.Close()

	data = [][]string{
		{"trxid", "amount", "type", "date"},
	}

	for _, txn := range summary.UnmatchedSystemTxn {
		data = append(data, []string{
			txn.TrxID,
			strconv.FormatFloat(txn.Amount, 'f', 2, 64),
			string(txn.Type),
			txn.TransactionTime.Format("2006-01-02 15:04:05"),
		})
	}

	writer = csv.NewWriter(file)
	if err := writer.WriteAll(data); err != nil {
		return fmt.Errorf("WriteCSVReport: failed write file: %w", err)
	}

	file, err = os.Create(fmt.Sprintf("%s/%s_", output, fmt.Sprintf("bank_stmts_missing_system_txns_%s_%s.csv", summary.StartDate, summary.EndDate)))
	if err != nil {
		return fmt.Errorf("WriteCSVReport: failed create file: %w", err)
	}
	defer file.Close()

	data = [][]string{
		{"uniqueid", "amount", "date", "bank"},
	}

	for bankName, statements := range summary.UnmatchedBankTx {
		for _, stmt := range statements {
			data = append(data, []string{
				stmt.UniqueIdentifier,
				strconv.FormatFloat(stmt.Amount, 'f', 2, 64),
				stmt.Date.Format("2006-01-02"),
				bankName,
			})
		}
	}

	writer = csv.NewWriter(file)
	if err := writer.WriteAll(data); err != nil {
		return fmt.Errorf("WriteCSVReport: failed write file: %w", err)
	}

	file, err = os.Create(fmt.Sprintf("%s/%s_", output, fmt.Sprintf("discrepant_txns_%s_%s.csv", summary.StartDate, summary.EndDate)))
	if err != nil {
		return fmt.Errorf("WriteCSVReport: failed create file: %w", err)
	}
	defer file.Close()

	data = [][]string{
		{"txnid", "amount", "bankstmtid", "bankamount", "difference"},
	}

	for _, dt := range summary.DiscrepantTransactions {
		data = append(data, []string{
			dt.SystemTxn.TrxID,
			strconv.FormatFloat(dt.SystemTxn.Amount, 'f', 2, 64),
			dt.BankStatement.UniqueIdentifier,
			strconv.FormatFloat(dt.BankStatement.Amount, 'f', 2, 64),
			strconv.FormatFloat(dt.Difference, 'f', 2, 64),
		})
	}

	writer = csv.NewWriter(file)
	if err := writer.WriteAll(data); err != nil {
		return fmt.Errorf("WriteCSVReport: failed write file: %w", err)
	}
	return nil
}

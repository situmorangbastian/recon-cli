package reconcile

import (
	"fmt"

	"github.com/situmorangbastian/recon-cli/internal"
)

type UnmatchedTransaction struct {
	SystemTxn     *internal.Transaction
	BankStatement *internal.BankStatement
}

type DiscrepantTransaction struct {
	SystemTxn     internal.Transaction
	BankStatement *internal.BankStatement
	Difference    float64
}

type ReconcileSummary struct {
	TotalTransactionsProcessed int
	TotalMatchedTransactions   int
	TotalUnmatchedTransactions int
	UnmatchedSystemTxn         []internal.Transaction
	UnmatchedBankTx            map[string][]internal.BankStatement
	TotalDiscrepancies         float64
	DiscrepantTransactions     []DiscrepantTransaction
}

func (rs *ReconcileSummary) PrintSummary() {
	fmt.Println("\n==================== TRANSACTION RECONCILIATION SUMMARY ====================")
	fmt.Printf("Total Transactions Processed : %d\n", rs.TotalTransactionsProcessed)
	fmt.Printf("Total Matched Transactions   : %d\n", rs.TotalMatchedTransactions)
	fmt.Printf("Total Unmatched Transactions : %d\n", rs.TotalUnmatchedTransactions)
	fmt.Printf("Total Discrepancy Amount     : Rp %.2f\n", rs.TotalDiscrepancies)

	if len(rs.UnmatchedSystemTxn) > 0 {
		fmt.Printf("\nSYSTEM TRANSACTIONS MISSING IN BANK STATEMENTS (%d)\n", len(rs.UnmatchedSystemTxn))
		fmt.Printf("%-20s %-15s %-10s %-20s\n", "TrxID", "Amount (Rp)", "Type", "Date")
		for _, tx := range rs.UnmatchedSystemTxn {
			fmt.Printf("%-20s Rp %-12.2f %-10s %-20s\n",
				tx.TrxID, tx.Amount, tx.Type, tx.TransactionTime.Format("2006-01-02 15:04:05"))
		}
	}

	if len(rs.UnmatchedBankTx) > 0 {
		fmt.Printf("\nBANK STATEMENTS MISSING IN SYSTEM TRANSACTIONS (%d)\n", len(rs.UnmatchedBankTx))
		fmt.Printf("%-25s %-15s %-15s %-15s\n", "Unique ID", "Amount (Rp)", "Date", "Bank")
		for bankName, statements := range rs.UnmatchedBankTx {
			for _, stmt := range statements {
				fmt.Printf("%-25s Rp %-12.2f %-15s %-15s\n",
					stmt.UniqueIdentifier, stmt.Amount, stmt.Date.Format("2006-01-02"), bankName)
			}
		}
	}

	if len(rs.DiscrepantTransactions) > 0 {
		fmt.Printf("\nDISCREPANT TRANSACTIONS (%d)\n", len(rs.DiscrepantTransactions))
		fmt.Printf("%-20s %-20s %-20s %-17s %-20s\n", "System TrxID", "System Amount", "Bank Stmt ID", "Bank Amount", "Difference (Rp)")
		for _, dt := range rs.DiscrepantTransactions {
			fmt.Printf("%-20s Rp %-17.2f %-20s Rp %-14.2f Rp %-13.2f\n",
				dt.SystemTxn.TrxID, dt.SystemTxn.Amount,
				dt.BankStatement.UniqueIdentifier, dt.BankStatement.Amount,
				dt.Difference)
		}
	}
	fmt.Println("============================================================================")
}

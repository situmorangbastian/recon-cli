package reconcile

import (
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

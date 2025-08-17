package internal

import (
	"time"
)

type TransactionType string

const (
	DEBIT  TransactionType = "DEBIT"
	CREDIT TransactionType = "CREDIT"
)

type Transaction struct {
	TrxID           string
	Amount          float64
	Type            TransactionType
	TransactionTime time.Time
	BankRefNo       string
}

type BankStatement struct {
	UniqueIdentifier string
	Amount           float64
	Date             time.Time
	File             string
}

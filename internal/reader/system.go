package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/situmorangbastian/recon-cli/internal"
)

func (r *Reader) ReadSysTxnsCSV() ([]internal.Transaction, error) {
	f, err := os.Open(r.sysTxnFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("ReadSysTxnsCSV: failed to read header: %w", err)
	}

	indices := make(map[string]int)
	for i, col := range header {
		indices[strings.ToLower(strings.TrimSpace(col))] = i
	}

	requiredCols := []string{"trxid", "amount", "type", "transactiontime"}
	for _, col := range requiredCols {
		if _, exists := indices[col]; !exists {
			return nil, fmt.Errorf("ReadSysTxnsCSV: required column '%s' not found in system transactions", col)
		}
	}

	var transactions []internal.Transaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ReadSysTxnsCSV: failed read records: %w", err)
		}

		amount, err := strconv.ParseFloat(record[indices["amount"]], 64)
		if err != nil {
			log.Printf("ReadSysTxnsCSV: warning: invalid amount '%s', skipping", record[indices["amount"]])
			continue
		}

		txnType := internal.TransactionType(strings.ToUpper(strings.TrimSpace(record[indices["type"]])))
		if txnType != internal.DEBIT && txnType != internal.CREDIT {
			log.Printf("ReadSysTxnsCSV: warning: invalid transaction type '%s', skipping", record[indices["type"]])
			continue
		}

		txnTime, err := r.parseSysTxnDateTimeLayout(record[indices["transactiontime"]])
		if err != nil {
			log.Printf("ReadSysTxnsCSV: warning: invalid transaction time '%s', skipping: %v", err, record[indices["trxid"]])
			continue
		}

		transaction := internal.Transaction{
			TrxID:           strings.TrimSpace(record[indices["trxid"]]),
			Amount:          amount,
			Type:            txnType,
			TransactionTime: txnTime,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *Reader) parseSysTxnDateTimeLayout(dateTime string) (time.Time, error) {
	dateTime = strings.TrimSpace(dateTime)
	for _, format := range r.sysTxnDateTimeLayout {
		if t, err := time.Parse(format, dateTime); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", dateTime)
}

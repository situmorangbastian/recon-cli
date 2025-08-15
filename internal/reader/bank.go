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

func (r *Reader) ReadBankStmtsCSV() ([]internal.BankStatement, error) {
	var bankStmts []internal.BankStatement
	for _, filePath := range r.bankStmtFilePaths {
		statements, err := r.readSingleBankStatement(filePath)
		if err != nil {
			return nil, fmt.Errorf("ReadBankStmtsCSV: failed to read: %w", err)
		}
		bankStmts = append(bankStmts, statements...)
	}

	return bankStmts, nil
}

func (r *Reader) readSingleBankStatement(filePath string) ([]internal.BankStatement, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var statements []internal.BankStatement

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("readSingleBankStatement: failed to read header: %w", err)
	}

	indices := make(map[string]int)
	for i, col := range header {
		indices[strings.ToLower(strings.TrimSpace(col))] = i
	}

	requiredCols := []string{"unique_identifier", "amount", "date"}
	for _, col := range requiredCols {
		if _, exists := indices[col]; !exists {
			return nil, fmt.Errorf("readSingleBankStatement: required column '%s' not found in bank statement: %s", col, file.Name())
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("readSingleBankStatement: failed to read record: %w", err)
		}

		amount, err := strconv.ParseFloat(record[indices["amount"]], 64)
		if err != nil {
			log.Printf("readSingleBankStatement: warning: invalid amount '%s' in bank statement: %s, skipping", record[indices["amount"]], file.Name())
			continue
		}

		date, err := r.parseBankStmtDateLayout(record[indices["date"]])
		if err != nil {
			log.Printf("readSingleBankStatement: warning: invalid date '%s', skipping: %v", record[indices["date"]], err)
			continue
		}

		statement := internal.BankStatement{
			UniqueIdentifier: strings.TrimSpace(record[indices["unique_identifier"]]),
			Amount:           amount,
			Date:             date,
			File:             filePath,
		}

		statements = append(statements, statement)
	}

	return statements, nil
}

func (r *Reader) parseBankStmtDateLayout(dateTime string) (time.Time, error) {
	dateTime = strings.TrimSpace(dateTime)
	for _, format := range r.bankStmtDateLayout {
		if t, err := time.Parse(format, dateTime); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateTime)
}

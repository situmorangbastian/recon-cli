# recon-cli: Transaction Reconciliation Service

`recon-cli` is a CLI tool to match internal system transactions with external bank statement records.
It detects unmatched and discrepant transactions within a specified date range.

---

## Features

- Match transactions by **bankfrefno**
- Supports **multiple bank statement files** (comma-separated)
- Flags discrepancy transactions where bank reference matches but amounts differ.
- Optional `-output` flag to export result reconcile as csv file

---

## Usage

```bash
go run cmd/reconcli/main.go \
  -startdate="2025-08-01"  \
  -enddate="2025-08-31" \
  -systempath="testdata/system_transactions.csv" \
  -bankstmtpath="testdata/bank_a_statement.csv,testdata/bank_b_statement.csv"
```

### Flags

| Flag                    | Description                                                              |
|-------------------------|--------------------------------------------------------------------------|
| `-systempath`           | path file to system transactions (required)                              |
| `-bankstmtpath`.        | comma-separated path files to bank statement (required)                  |
| `-startdate`            | startdate in YYYY-MM-DD format (required)                                |
| `-enddate`              | enddate in YYYY-MM-DD format (required)                                  |
| `-output`               | Target folder for CSV output (optional, default: print to terminal)      |

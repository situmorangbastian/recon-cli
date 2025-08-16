# recon-cli: Transaction Reconciliation Service

`recon-cli` is a CLI tool to match internal system transactions with external bank statement records.
It detects unmatched and discrepant transactions within a specified date range.

---

## Features

- Match transactions by **date** and **amount**
- Supports **multiple bank statement files** (comma-separated)
- Includes a **1,000 tolerance threshold** for real-world discrepancies
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

---

## Matching Logic

### Match Criteria

1. **Date Match:**
   System transaction date must exactly match the bank statement date.
   (Time portion is ignored — dates are truncated to `YYYY-MM-DD`).

2. **Amount Match with Tolerance:**
   For transactions, a **tolerance threshold of 1,000 is applied:
   - Accepts small mismatches due to rounding, fees, or currency adjustment.

### Example

| System Transaction | Bank Statement | Difference | Match? |
|--------------------|----------------|------------|--------|
| 1000000            | 999200         | 800        | ✅     |
| 1000000            | 998500         | 1500      | ❌     |

---

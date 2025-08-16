package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/situmorangbastian/recon-cli/internal/reader"
	"github.com/situmorangbastian/recon-cli/internal/reconcile"
	"github.com/situmorangbastian/recon-cli/internal/service"
)

var (
	sysTxnDateTimeLayout = []string{
		"2006-01-02 15:04:05",
	}
	bankStmtDateLayout = []string{
		"2006-01-02",
	}
)

func main() {
	flag.Usage = func() {
		progName := filepath.Base(os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", progName)
		flag.PrintDefaults()
	}

	systemTxnPath := flag.String("systempath", "", "path file to system transactions (required)")
	bankStmtPath := flag.String("bankstmtpath", "", "comma-separated path files to bank statement (required)")
	startDate := flag.String("startdate", "", "startdate in YYYY-MM-DD format (required)")
	endDate := flag.String("enddate", "", "enddate in YYYY-MM-DD format (required)")
	output := flag.String("output", "", "Target folder for CSV output (optional, default: print to terminal)")

	flag.Parse()

	switch true {
	case *startDate == "":
		fmt.Println("Error: -startdate required")
		flag.Usage()
		os.Exit(1)
	case *endDate == "":
		fmt.Println("Error: -enddate required")
		flag.Usage()
		os.Exit(1)
	case *systemTxnPath == "":
		fmt.Println("Error: -systempath required")
		flag.Usage()
		os.Exit(1)
	case *bankStmtPath == "":
		fmt.Println("Error: -bankstmtpath required")
		flag.Usage()
		os.Exit(1)
	}

	bankStmtPaths := strings.Split(*bankStmtPath, ",")
	reader := reader.NewReader(sysTxnDateTimeLayout, bankStmtDateLayout, bankStmtPaths, *systemTxnPath)
	svc := service.NewService(reader)
	reconcile := reconcile.New(svc)

	resultReconciled, err := reconcile.Reconcile(*startDate, *endDate)
	if err != nil {
		log.Fatalf("reconciliation failed: %v", err)
	}

	if *output == "" {
		resultReconciled.PrintSummary()
		return
	}

	err = reconcile.WriteCSVReport(*resultReconciled, *output)
	if err != nil {
		log.Fatalf("write csv report failed: %v", err)
	}
}

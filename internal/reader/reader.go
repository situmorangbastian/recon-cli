package reader

type Reader struct {
	sysTxnDateTimeLayout []string
	bankStmtDateLayout   []string
	sysTxnFilePath       string
	bankStmtFilePaths    []string
}

func NewReader(sysTxnDateTimeLayout, bankStmtDateLayout, bankStmtFilePaths []string, sysTxnFilePath string) *Reader {
	return &Reader{
		sysTxnDateTimeLayout: sysTxnDateTimeLayout,
		bankStmtDateLayout:   bankStmtDateLayout,
		sysTxnFilePath:       sysTxnFilePath,
		bankStmtFilePaths:    bankStmtFilePaths,
	}
}

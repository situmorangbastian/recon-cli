package reader

type Reader struct {
	sysTxnDateTimeLayout []string
	bankStmtDateLayout   []string
}

func NewReader(sysTxnDateTimeLayout, bankStmtDateLayout []string) *Reader {
	return &Reader{
		sysTxnDateTimeLayout: sysTxnDateTimeLayout,
		bankStmtDateLayout:   bankStmtDateLayout,
	}
}

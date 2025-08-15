package reader

type Reader struct {
	sysTxnDateTimeLayout []string
}

func NewReader(sysTxnDateTimeLayout []string) *Reader {
	return &Reader{
		sysTxnDateTimeLayout: sysTxnDateTimeLayout,
	}
}

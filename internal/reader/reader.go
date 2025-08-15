package reader

type Reader struct {
	SysTxnDateTimeLayout []string
}

func NewReader(SysTxnDateTimeLayout []string) *Reader {
	return &Reader{
		SysTxnDateTimeLayout: SysTxnDateTimeLayout,
	}
}

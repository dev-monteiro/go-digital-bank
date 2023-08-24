package commons

type CoreBankingInvoiceResponse struct {
	CreditAccountId     int32
	ProcessingSituation string
	IsPaymentDone       bool
	DueDate             string
	ActualDueDate       string
	ClosingDate         string
	TotalAmount         float32
	InvoiceId           int
}

type CoreBankingInvoiceListResponse struct {
	Invoices []CoreBankingInvoiceResponse
}

type PurchaseEvent struct {
	PurchaseId          int
	CreditAccountId     int
	PurchaseDateTime    string
	Amount              float32
	NumInstallments     int
	MerchantDescription string
	Status              string
	Description         string
}

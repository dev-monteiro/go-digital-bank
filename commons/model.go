package commons

type CoreBankInvoiceResp struct {
	CreditAccountId     int32
	ProcessingSituation string // TODO: convert to some kind of enum?
	IsPaymentDone       bool
	DueDate             string // TODO: convert to date
	ActualDueDate       string
	ClosingDate         string
	TotalAmount         *MoneyAmount
	InvoiceId           int
}

type CoreBankInvoiceListResp struct {
	Invoices []CoreBankInvoiceResp
}

type PurchaseEvent struct {
	PurchaseId          int
	CreditAccountId     int
	PurchaseDateTime    string
	Amount              *MoneyAmount
	NumInstallments     int
	MerchantDescription string
	Status              string
	Description         string
}

type BatchEvent struct {
	BatchId       int
	ReferenceDate string
}

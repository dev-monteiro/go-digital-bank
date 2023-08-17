package main

type CustomerCredential struct {
	CustomerId      string
	CreditAccountId int
}

type CurrentInvoiceResponse struct {
	Id          string
	StatusLabel string
	Amount      float32
	ClosingDate string
}

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

type Purchase struct {
	PurchaseId          int
	CreditAccountId     int
	PurchaseDateTime    string
	Amount              float32
	NumInstallments     int
	MerchantDescription string
}

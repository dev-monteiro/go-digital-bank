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

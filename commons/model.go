package commons

import "dev-monteiro/go-digital-bank/commons/invostatus"

type CoreBankInvoiceResp struct {
	CustomerId    int32                 `json:"accountId"`
	Status        invostatus.InvoStatus `json:"processingSituation"`
	IsPaymentDone bool                  `json:"paymentDone"`
	DueDate       *LocalDate            `json:"invoiceDueDate"`
	ActualDueDate *LocalDate            `json:"realDueDate"`
	ClosingDate   *LocalDate            `json:"closingDate"`
	Amount        *MoneyAmount          `json:"totalAmount"`
}

type CoreBankInvoiceListResp struct {
	Invoices []CoreBankInvoiceResp `json:"content"`
}

type PurchaseEvent struct {
	Id                  int          `json:"purchase_id"`
	CustomerId          int          `json:"account_id"`
	DateTime            string       `json:"purchase_date"`
	Amount              *MoneyAmount `json:"amount"`
	NumInstallments     int          `json:"installment"`
	MerchantDescription string       `json:"merchant"`
	Status              string       `json:"status"`
}

type BatchEvent struct {
	Id            int        `json:"process_control_id"`
	ReferenceDate *LocalDate `json:"processing_date"`
}

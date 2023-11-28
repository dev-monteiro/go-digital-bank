package commons

import (
	"dev-monteiro/go-digital-bank/commons/invstat"
	"dev-monteiro/go-digital-bank/commons/ldate"
	"dev-monteiro/go-digital-bank/commons/ldatetime"
	"dev-monteiro/go-digital-bank/commons/mnyamnt"
)

type CoreBankInvoiceResp struct {
	CustomerId    int32              `json:"accountId"`
	Status        invstat.InvoStatus `json:"processingSituation"`
	IsPaymentDone bool               `json:"paymentDone"`
	DueDate       *ldate.LocDate     `json:"invoiceDueDate"`
	ActualDueDate *ldate.LocDate     `json:"realDueDate"`
	ClosingDate   *ldate.LocDate     `json:"closingDate"`
	Amount        *mnyamnt.MnyAmount `json:"totalAmount"`
}

type CoreBankInvoiceListResp struct {
	Invoices []CoreBankInvoiceResp `json:"content"`
}

type PurchaseEvent struct {
	Id                  int                    `json:"purchase_id"`
	CustomerId          int                    `json:"account_id"`
	DateTime            *ldatetime.LocDateTime `json:"purchase_date"`
	Amount              *mnyamnt.MnyAmount     `json:"amount"`
	NumInstallments     int                    `json:"installment"`
	MerchantDescription string                 `json:"merchant"`
	Status              string                 `json:"status"`
}

type BatchEvent struct {
	Id            int            `json:"process_control_id"`
	ReferenceDate *ldate.LocDate `json:"processing_date"`
}

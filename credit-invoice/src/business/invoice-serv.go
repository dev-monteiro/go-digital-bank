package business

import (
	comm "devv-monteiro/go-digital-bank/commons"
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"

	"encoding/json"
	"net/http"
	"strconv"
)

type InvoiceServ struct {
	credRepo  *data.CredentialRepo
	purchRepo *data.PurchaseRepo
}

func NewInvoiceServ(credRepo *data.CredentialRepo, purchRepo *data.PurchaseRepo) *InvoiceServ {
	return &InvoiceServ{
		credRepo:  credRepo,
		purchRepo: purchRepo,
	}
}

func (serv *InvoiceServ) GetCurrentInvoice(customerId string) (*CurrInvoiceResp, *conf.AppError) {
	creditAccountId, err := serv.credRepo.GetCreditAccountId(customerId)
	if err != nil {
		return nil, err
	}

	invoice, err := serv.getCoreBankInvoice(creditAccountId)
	if err != nil {
		return nil, err
	}

	updatedAmount, err := serv.updateAmount(creditAccountId, invoice.TotalAmount)
	if err != nil {
		return nil, err
	}

	response := CurrInvoiceResp{
		Id:          strconv.Itoa(invoice.InvoiceId),
		StatusLabel: invoice.ProcessingSituation,
		Amount:      updatedAmount,
		ClosingDate: invoice.ClosingDate,
	}

	return &response, nil
}

func (serv *InvoiceServ) updateAmount(creditAccountId int, currentAmount float32) (float32, *conf.AppError) {
	purchases, err := serv.purchRepo.FindAllByCreditAccountId(creditAccountId)

	if err != nil {
		return 0, err
	}

	for _, purchase := range purchases {
		currentAmount += purchase.Amount
	}
	return currentAmount, nil
}

func (InvoiceServ) getCoreBankInvoice(creditAccountId int) (*comm.CoreBankInvoiceResp, *conf.AppError) {
	url := "http://core_banking_mock/invoices?creditAccountId=" + strconv.Itoa(creditAccountId)

	coreBankResp, err := http.Get(url)
	if err != nil {
		return nil, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	var invoiceList comm.CoreBankInvoiceListResp
	json.NewDecoder(coreBankResp.Body).Decode(&invoiceList)
	invoice := invoiceList.Invoices[0]
	return &invoice, nil
}

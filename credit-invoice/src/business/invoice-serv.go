package business

import (
	comm "devv-monteiro/go-digital-bank/commons"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"

	"encoding/json"
	"fmt"
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

func (serv *InvoiceServ) GetCurrentInvoice(customerId string) (*CurrInvoiceResp, error) {
	creditAccountId, err := serv.credRepo.GetCreditAccountId(customerId)
	if err != nil {
		return nil, err
	}

	invoice, err := serv.getCoreBankInvoice(creditAccountId)
	if err != nil {
		return nil, err
	}

	updatedAmount := serv.updateAmount(creditAccountId, invoice.TotalAmount)

	response := CurrInvoiceResp{
		Id:          strconv.Itoa(invoice.InvoiceId),
		StatusLabel: invoice.ProcessingSituation,
		Amount:      updatedAmount,
		ClosingDate: invoice.ClosingDate,
	}

	return &response, nil
}

func (serv *InvoiceServ) updateAmount(creditAccountId int, currentAmount float32) float32 {
	purchases, err := serv.purchRepo.FindAllByCreditAccountId(creditAccountId)

	if err != nil {
		fmt.Println(err)
		return currentAmount
	}

	for _, purchase := range purchases {
		currentAmount += purchase.Amount
	}
	return currentAmount
}

func (InvoiceServ) getCoreBankInvoice(creditAccountId int) (*comm.CoreBankInvoiceResp, error) {
	url := "http://core_banking_mock/invoices?creditAccountId=" + strconv.Itoa(creditAccountId)

	coreBankResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var invoiceList comm.CoreBankInvoiceListResp
	json.NewDecoder(coreBankResp.Body).Decode(&invoiceList)
	invoice := invoiceList.Invoices[0]
	return &invoice, nil
}

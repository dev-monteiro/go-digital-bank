package main

import (
	"devv-monteiro/go-digital-bank/commons"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type InvoiceServ struct {
	credRepo  *CredentialRepo
	purchRepo *PurchaseRepo
}

func NewInvoiceServ(credRepo *CredentialRepo, purchRepo *PurchaseRepo) *InvoiceServ {
	return &InvoiceServ{
		credRepo:  credRepo,
		purchRepo: purchRepo,
	}
}

func (serv *InvoiceServ) GetCurrentInvoice(customerId string) (*CurrentInvoiceResponse, error) {
	creditAccountId, err := serv.credRepo.getCreditAccountId(customerId)
	if err != nil {
		return nil, err
	}

	invoice, err := serv.getCoreBankingInvoice(creditAccountId)
	if err != nil {
		return nil, err
	}

	updatedAmount := serv.updateAmount(creditAccountId, invoice.TotalAmount)

	response := CurrentInvoiceResponse{
		Id:          strconv.Itoa(invoice.InvoiceId),
		StatusLabel: invoice.ProcessingSituation,
		Amount:      updatedAmount,
		ClosingDate: invoice.ClosingDate,
	}

	return &response, nil
}

func (serv *InvoiceServ) updateAmount(creditAccountId int, currentAmount float32) float32 {
	purchases, err := serv.purchRepo.findAllByCreditAccountId(creditAccountId)

	if err != nil {
		fmt.Println(err)
		return currentAmount
	}

	for _, purchase := range purchases {
		currentAmount += purchase.Amount
	}
	return currentAmount
}

func (InvoiceServ) getCoreBankingInvoice(creditAccountId int) (*commons.CoreBankingInvoiceResponse, error) {
	url := "http://core_banking_mock/invoices?creditAccountId=" + strconv.Itoa(creditAccountId)

	invoiceResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var invoiceList commons.CoreBankingInvoiceListResponse
	json.NewDecoder(invoiceResp.Body).Decode(&invoiceList)
	invoice := invoiceList.Invoices[0]
	return &invoice, nil
}

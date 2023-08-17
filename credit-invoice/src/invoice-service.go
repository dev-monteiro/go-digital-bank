package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type InvoiceService struct {
	credentialRepo CredentialRepository
	purchaseRepo   PurchaseRepository
}

func NewInvoiceService(credentialRepo CredentialRepository, purchaseRepo PurchaseRepository) InvoiceService {
	return InvoiceService{
		credentialRepo: credentialRepo,
		purchaseRepo:   purchaseRepo,
	}
}

func (serv *InvoiceService) getCurrentInvoice(customerId string) (*CurrentInvoiceResponse, error) {
	creditAccountId, err := serv.credentialRepo.getCreditAccountId(customerId)
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

func (serv *InvoiceService) updateAmount(creditAccountId int, currentAmount float32) float32 {
	purchases, err := serv.purchaseRepo.findAllByCreditAccountId(creditAccountId)

	if err != nil {
		fmt.Println(err)
		return currentAmount
	}

	for _, purchase := range purchases {
		currentAmount += purchase.Amount
	}
	return currentAmount
}

func (InvoiceService) getCoreBankingInvoice(creditAccountId int) (*CoreBankingInvoiceResponse, error) {
	url := "http://core_banking_mock/invoices?creditAccountId=" + strconv.Itoa(creditAccountId)
	fmt.Println("Url = " + url)

	invoiceResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var invoiceList CoreBankingInvoiceListResponse
	json.NewDecoder(invoiceResp.Body).Decode(&invoiceList)
	invoice := invoiceList.Invoices[0]
	return &invoice, nil
}

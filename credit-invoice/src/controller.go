package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Controller struct {
	repo Repository
}

func NewController() Controller {
	return Controller{NewRepository()}
}

func (cont Controller) Close() {
	cont.repo.Close()
}

func (cont Controller) getCurrentInvoice(resWr http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	customerId := req.Form.Get("customerId")
	fmt.Println("CustomerId = " + customerId)

	creditAccountId, err := cont.repo.getCreditAccountId(customerId)
	if err != nil {
		resWr.WriteHeader(http.StatusBadRequest)
		return
	}

	url := "http://core_banking_mock/invoices?creditAccountId=" + strconv.Itoa(creditAccountId)
	fmt.Println("Url = " + url)

	cbListResWrapper, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		resWr.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cbListRes CoreBankingInvoiceListResponse
	json.NewDecoder(cbListResWrapper.Body).Decode(&cbListRes)
	cbRes := cbListRes.Invoices[0]

	response := CurrentInvoiceResponse{
		Id:          strconv.Itoa(cbRes.InvoiceId),
		StatusLabel: cbRes.ProcessingSituation,
		Amount:      cbRes.TotalAmount,
		ClosingDate: cbRes.ClosingDate,
	}

	resWr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resWr).Encode(response)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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

func getCurrentInvoice(resWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	customerId := request.Form.Get("customerId")

	if customerId != "abc-123-def" {
		resWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(customerId)
	cbListResWrapper, err := http.Get("http://core-banking-mock/invoices?creditAccountId=123")
	if err != nil {
		fmt.Println(err.Error())
		resWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(cbListResWrapper)

	var cbListRes CoreBankingInvoiceListResponse
	json.NewDecoder(cbListResWrapper.Body).Decode(&cbListRes)
	cbRes := cbListRes.Invoices[0]
	fmt.Println(cbRes)

	response := CurrentInvoiceResponse{
		Id:          strconv.Itoa(cbRes.InvoiceId),
		StatusLabel: cbRes.ProcessingSituation,
		Amount:      cbRes.TotalAmount,
		ClosingDate: cbRes.ClosingDate,
	}

	resWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resWriter).Encode(response)
}

func main() {
	http.HandleFunc("/invoices/current", getCurrentInvoice)

	http.ListenAndServe(":80", nil)
}

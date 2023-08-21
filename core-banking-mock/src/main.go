package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type InvoiceResponse struct {
	CreditAccountId     int
	ProcessingSituation string
	IsPaymentDone       bool
	DueDate             string
	ActualDueDate       string
	ClosingDate         string
	TotalAmount         float32
	InvoiceId           int
}

type InvoiceListResponse struct {
	Invoices []InvoiceResponse
}

func main() {
	fmt.Println("Setup completed")

	http.HandleFunc("/invoices", getInvoices)

	http.ListenAndServe(":80", nil)
}

func getInvoices(resWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	creditAccountId := request.Form.Get("creditAccountId")

	fmt.Println("CreditAccountId = " + creditAccountId)
	if creditAccountId != "123" {
		resWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	invoice := InvoiceResponse{
		CreditAccountId:     123,
		ProcessingSituation: "OPEN",
		IsPaymentDone:       false,
		DueDate:             "2023-08-20",
		ActualDueDate:       "2023-08-21",
		ClosingDate:         "2023-08-15",
		TotalAmount:         1234.56,
		InvoiceId:           1234,
	}

	invoiceList := InvoiceListResponse{
		Invoices: []InvoiceResponse{invoice},
	}

	resWriter.Header().Add("Content-Type", "application/json")
	json.NewEncoder(resWriter).Encode(invoiceList)
}

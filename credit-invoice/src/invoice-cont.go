package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type InvoiceCont struct {
	serv *InvoiceServ
}

func NewInvoiceCont(serv *InvoiceServ) *InvoiceCont {
	return &InvoiceCont{serv: serv}
}

func (cont *InvoiceCont) getCurrInvoice(resWr http.ResponseWriter, req *http.Request) {
	fmt.Println("Path: " + req.URL.Path)

	if req.Method != http.MethodGet {
		resWr.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req.ParseForm()
	customerId := req.Form.Get("customerId")
	fmt.Println("CustomerId: " + customerId)

	currentInvoice, err := cont.serv.GetCurrentInvoice(customerId)
	if err != nil {
		resWr.WriteHeader(http.StatusBadRequest)
		return
	}

	resWr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resWr).Encode(currentInvoice)
}

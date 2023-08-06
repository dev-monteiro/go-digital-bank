package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Controller struct {
	serv InvoiceService
}

func NewController(serv InvoiceService) Controller {
	return Controller{serv: serv}
}

func (cont Controller) getCurrentInvoice(resWr http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	customerId := req.Form.Get("customerId")
	fmt.Println("CustomerId = " + customerId)

	currentInvoice, err := cont.serv.getCurrentInvoice(customerId)
	if err != nil {
		resWr.WriteHeader(http.StatusBadRequest)
		return
	}

	resWr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resWr).Encode(currentInvoice)
}

package transport

import (
	"devv-monteiro/go-digital-bank/credit-invoice/src/business"
	"encoding/json"
	"log"
	"net/http"
)

type InvoiceCont struct {
	serv business.InvoiceServ
}

func NewInvoiceCont(serv business.InvoiceServ) *InvoiceCont {
	cont := &InvoiceCont{serv: serv}
	http.HandleFunc("/invoices/current", cont.GetCurrInvoice)
	return cont
}

func (cont *InvoiceCont) GetCurrInvoice(resWr http.ResponseWriter, req *http.Request) {
	log.Println("[InvoiceCont] GetCurrInvoice")

	resWr.Header().Set("Content-Type", "application/json")

	if req.Method != http.MethodGet {
		resWr.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req.ParseForm()
	customerId := req.Form.Get("customerId")
	log.Println("[InvoiceCont] CustomerId: " + customerId)

	if customerId == "" {
		resWr.WriteHeader(http.StatusBadRequest)
		return
	}

	currentInvoice, err := cont.serv.GetCurrInvoice(customerId)
	if err != nil {
		resWr.WriteHeader(err.StatusCode)
		json.NewEncoder(resWr).Encode(err)
		return
	}

	json.NewEncoder(resWr).Encode(currentInvoice)
}

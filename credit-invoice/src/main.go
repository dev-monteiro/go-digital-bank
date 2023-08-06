package main

import (
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	time.Sleep(10 * time.Second)
}

func main() {
	credentialRepo := NewCredentialRepo()
	purchaseRepo := NewPurchaseRepo()

	listener := NewListener(purchaseRepo)
	defer listener.Close()

	invoiceServ := NewInvoiceService(credentialRepo, purchaseRepo)
	controller := NewController(invoiceServ)

	http.HandleFunc("/invoices/current", controller.getCurrentInvoice)

	http.ListenAndServe(":80", nil)
}

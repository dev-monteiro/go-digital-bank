package main

import (
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: use some real logging library
// TODO: use best practices for constants and env variables
func main() {
	time.Sleep(20 * time.Second)

	dynamoCli := NewDynamoClient()
	sqsCli := NewSqsClient()

	credentialRepo := NewCredentialRepository(dynamoCli)
	purchaseRepo := NewPurchaseRepository(dynamoCli)

	NewListener(sqsCli, purchaseRepo)

	invoiceServ := NewInvoiceService(credentialRepo, purchaseRepo)
	controller := NewController(invoiceServ)

	http.HandleFunc("/invoices/current", controller.getCurrentInvoice)

	http.ListenAndServe(":80", nil)
}

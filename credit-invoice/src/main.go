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

	dynamoClient := NewDynamoDbClient()
	sqsClient := NewSqsClient()

	credentialRepo := NewCredentialRepo(dynamoClient)
	purchaseRepo := NewPurchaseRepo(dynamoClient)

	NewListener(sqsClient, purchaseRepo)

	invoiceServ := NewInvoiceService(credentialRepo, purchaseRepo)
	controller := NewController(invoiceServ)

	http.HandleFunc("/invoices/current", controller.getCurrentInvoice)

	http.ListenAndServe(":80", nil)
}

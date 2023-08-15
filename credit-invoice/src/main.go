package main

import (
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	time.Sleep(20 * time.Second)
}

func main() {
	dynamoDB := NewDynamoDbClient()

	credentialRepo := NewCredentialRepo(dynamoDB)
	purchaseRepo := NewPurchaseRepo(dynamoDB)

	//listener := NewListener(purchaseRepo)
	//defer listener.Close()

	invoiceServ := NewInvoiceService(credentialRepo, purchaseRepo)
	controller := NewController(invoiceServ)

	http.HandleFunc("/invoices/current", controller.getCurrentInvoice)

	http.ListenAndServe(":80", nil)
}

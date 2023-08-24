package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: use some real logging library
// TODO: use best practices for constants and env variables
func main() {
	controller := setup()

	http.HandleFunc("/invoices/current", controller.getCurrentInvoice)

	http.ListenAndServe(":80", nil)
}

func setup() *Controller {
	sqsCli := setupComponentWithRetry(NewSqsClient)
	dynamoCli := setupComponentWithRetry(NewDynamoClient)

	credentialRepo := NewCredentialRepository(dynamoCli)
	purchaseRepo := NewPurchaseRepository(dynamoCli)

	setupComponentWithRetry(func() (*Listener, error) { return NewListener(sqsCli, purchaseRepo) })

	invoiceServ := NewInvoiceService(credentialRepo, purchaseRepo)
	controller := NewController(invoiceServ)

	fmt.Println("Setup completed")
	return controller
}

func setupComponentWithRetry[T any](setupFunction func() (T, error)) T {
	for {
		component, err := setupFunction()

		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		return component
	}
}

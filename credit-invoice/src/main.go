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
	var controller *Controller

	for {
		dynamoCli, err := NewDynamoClient()
		sqsCli, err := NewSqsClient()

		credentialRepo := NewCredentialRepository(dynamoCli)
		purchaseRepo := NewPurchaseRepository(dynamoCli)

		NewListener(sqsCli, purchaseRepo)

		invoiceServ := NewInvoiceService(credentialRepo, purchaseRepo)
		controller = NewController(invoiceServ)

		if err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Second)
		}
		break
	}

	fmt.Println("Setup completed")
	return controller
}

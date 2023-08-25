package main

import (
	"fmt"
	"net/http"
	"time"
)

// TODO: use some real logging library
// TODO: use best practices for constants and env variables
func main() {
	invoCont := setup()

	http.HandleFunc("/invoices/current", invoCont.getCurrInvoice)

	http.ListenAndServe(":80", nil)
}

func setup() *InvoiceCont {
	sqsClnt := setupWithRetry(NewSqsClnt)
	dynaClnt := setupWithRetry(NewDynamoClnt)

	credRepo := NewCredentialRepo(dynaClnt)
	purchRepo := NewPurchaseRepo(dynaClnt)

	setupWithRetry(func() (*PurchaseListen, error) { return NewPurchaseListen(sqsClnt, purchRepo) })

	invoServ := NewInvoiceServ(credRepo, purchRepo)
	invoCont := NewInvoiceCont(invoServ)

	fmt.Println("Setup completed")
	return invoCont
}

func setupWithRetry[T any](setupFunc func() (T, error)) T {
	for {
		comp, err := setupFunc()

		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}

		return comp
	}
}

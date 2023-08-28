package main

import (
	busi "devv-monteiro/go-digital-bank/credit-invoice/src/business"
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"
	tran "devv-monteiro/go-digital-bank/credit-invoice/src/transport"
	"fmt"
	"net/http"
	"time"
)

// TODO: use some real logging library
// TODO: use best practices for constants and env variables
func main() {
	invoCont := setup()

	http.HandleFunc("/invoices/current", invoCont.GetCurrInvoice)

	http.ListenAndServe(":80", nil)
}

func setup() *tran.InvoiceCont {
	sqsClnt := setupWithRetry(conf.NewSqsClnt)
	dynaClnt := setupWithRetry(conf.NewDynamoClnt)

	credRepo := data.NewCredentialRepo(dynaClnt)
	purchRepo := data.NewPurchaseRepo(dynaClnt)

	invoServ := busi.NewInvoiceServ(credRepo, purchRepo)

	setupWithRetry(func() (*tran.PurchaseListen, error) { return tran.NewPurchaseListen(sqsClnt, purchRepo) })
	invoCont := tran.NewInvoiceCont(invoServ)

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

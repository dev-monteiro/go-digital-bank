package main

import (
	busi "devv-monteiro/go-digital-bank/credit-invoice/src/business"
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"
	tran "devv-monteiro/go-digital-bank/credit-invoice/src/transport"
	"log"
	"net/http"
	"time"
)

// TODO: use best practices for constants and env variables
// TODO: analyze concurrent behaviour
func main() {
	setupComponents()

	http.ListenAndServe(":80", nil)
}

func setupComponents() {
	sqsClnt := setupWithRetry(conf.NewSqsClnt)
	dynaClnt := setupWithRetry(conf.NewDynamoClnt)

	custRepo := data.NewCustomerRepo(dynaClnt)
	transcRepo := data.NewTransactionRepo(dynaClnt)

	invoServ := busi.NewInvoiceServ(custRepo, transcRepo)

	setupWithRetry(func() (*tran.PurchaseListen, error) { return tran.NewPurchaseListen(sqsClnt, transcRepo) })
	tran.NewInvoiceCont(invoServ)

	log.Println("[Main] Setup completed!")
}

func setupWithRetry[T any](setupFunc func() (T, error)) T {
	for {
		comp, err := setupFunc()

		if err != nil {
			log.Println("[Main] Setting up...")
			time.Sleep(5 * time.Second)
			continue
		}

		return comp
	}
}

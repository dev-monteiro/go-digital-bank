package main

import (
	comm "dev-monteiro/go-digital-bank/commons"
	busi "dev-monteiro/go-digital-bank/credit-invoice/src/business"
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	tran "dev-monteiro/go-digital-bank/credit-invoice/src/transport"
	"log"
	"net/http"
	"time"
)

// TODO: analyze concurrent behaviour
func main() {
	setupComponents()

	http.ListenAndServe(":80", nil)
}

func setupComponents() {
	sqsConn := setupWithRetry(conn.NewSqsConn)
	dynaConn := setupWithRetry(conn.NewDynamoConn)
	coreBankConn := conn.NewCoreBankConn()

	custRepo := data.NewCustomerRepo(dynaConn)
	transcRepo := data.NewTransactionRepo(dynaConn)

	invoServ := busi.NewInvoiceServ(custRepo, transcRepo, coreBankConn)
	transcServ := busi.NewTransactionServ(custRepo, transcRepo)

	setupWithRetry(func() (tran.Listen[comm.PurchaseEvent], error) {
		return tran.NewPurchaseListen(sqsConn, transcServ)
	})
	setupWithRetry(func() (tran.Listen[comm.BatchEvent], error) {
		return tran.NewBatchListen(sqsConn, transcServ)
	})

	tran.NewInvoiceCont(invoServ)
	tran.NewHealthCont()

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

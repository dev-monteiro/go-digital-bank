package transport

import (
	comm "dev-monteiro/go-digital-bank/commons"
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	"log"
)

type purchaseListen struct {
	sqsConn    conn.SqsConn
	transcRepo data.TransactionRepo
}

func NewPurchaseListen(sqsConn conn.SqsConn, transcRepo data.TransactionRepo) (Listen[comm.PurchaseEvent], error) {
	return SetupListen[comm.PurchaseEvent](sqsConn, &purchaseListen{
		sqsConn:    sqsConn,
		transcRepo: transcRepo,
	})
}

func (listen *purchaseListen) getQueueNameEnv() string {
	return "AWS_PURCHASES_QUEUE_NAME"
}

func (listen *purchaseListen) process(event comm.PurchaseEvent) {
	log.Println("[PurchaseListen] Event received")

	transc := data.Transaction{
		PurchaseId:         event.PurchaseId,
		CustomerCoreBankId: event.CreditAccountId,
		Amount:             event.Amount,
	}

	err := listen.transcRepo.Save(transc)
	if err != nil {
		log.Println("[PurchaseListen] " + err.Error())
	}
}

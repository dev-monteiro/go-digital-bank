package transport

import (
	comm "dev-monteiro/go-digital-bank/commons"
	busn "dev-monteiro/go-digital-bank/credit-invoice/src/business"
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	"log"
)

type purchaseListen struct {
	sqsConn    conn.SqsConn
	transcServ busn.TransactionServ
}

func NewPurchaseListen(sqsConn conn.SqsConn, transcServ busn.TransactionServ) (Listen[comm.PurchaseEvent], error) {
	return SetupListen[comm.PurchaseEvent](sqsConn, &purchaseListen{
		sqsConn:    sqsConn,
		transcServ: transcServ,
	})
}

func (listen *purchaseListen) getQueueNameEnv() string {
	return "AWS_PURCHASES_QUEUE_NAME"
}

func (listen *purchaseListen) process(event comm.PurchaseEvent) {
	err := listen.transcServ.CreateFromPurchase(event)
	if err != nil {
		log.Println("[PurchaseListen] Error: " + err.Message)
	}
}

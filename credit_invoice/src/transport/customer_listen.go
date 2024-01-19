package transport

import (
	comm "dev-monteiro/go-digital-bank/commons"
	busn "dev-monteiro/go-digital-bank/credit-invoice/src/business"
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	"log"
)

type customerListen struct {
	sqsConn  conn.SqsConn
	custServ busn.CustomerServ
}

func NewCustomerListen(sqsConn conn.SqsConn, custServ busn.CustomerServ) (Listen[comm.CustomerEvent], error) {
	return SetupListen[comm.CustomerEvent](sqsConn, &customerListen{
		sqsConn:  sqsConn,
		custServ: custServ,
	})
}

func (listen *customerListen) getQueueNameEnv() string {
	return "AWS_CUSTOMERS_QUEUE_NAME"
}

func (listen *customerListen) process(event comm.CustomerEvent) {
	err := listen.custServ.CreateFromEvent(&event)
	if err != nil {
		log.Println("[CustomerListen] Error: " + err.Message)
	}
}

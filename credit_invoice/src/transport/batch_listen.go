package transport

import (
	comm "dev-monteiro/go-digital-bank/commons"
	busn "dev-monteiro/go-digital-bank/credit-invoice/src/business"
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	"log"
)

type batchListen struct {
	sqsConn    conn.SqsConn
	transcServ busn.TransactionServ
}

func NewBatchListen(sqsConn conn.SqsConn, transcServ busn.TransactionServ) (Listen[comm.BatchEvent], error) {
	return SetupListen[comm.BatchEvent](sqsConn, &batchListen{
		sqsConn:    sqsConn,
		transcServ: transcServ,
	})
}

func (listen *batchListen) getQueueNameEnv() string {
	return "AWS_BATCHES_QUEUE_NAME"
}

func (listen *batchListen) process(event comm.BatchEvent) {
	err := listen.transcServ.ClearByBatch(&event)
	if err != nil {
		log.Println("[BatchListen] Error: " + err.Message)
	}
}

package transport

import (
	comm "dev-monteiro/go-digital-bank/commons"
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	"log"
)

type batchListen struct {
	sqsConn    conn.SqsConn
	custRepo   data.CustomerRepo
	transcRepo data.TransactionRepo
}

func NewBatchListen(sqsConn conn.SqsConn, custRepo data.CustomerRepo, transcRepo data.TransactionRepo) (Listen[comm.BatchEvent], error) {
	return SetupListen[comm.BatchEvent](sqsConn, &batchListen{
		sqsConn:    sqsConn,
		custRepo:   custRepo,
		transcRepo: transcRepo,
	})
}

func (listen *batchListen) getQueueNameEnv() string {
	return "AWS_BATCHES_QUEUE_NAME"
}

func (listen *batchListen) process(event comm.BatchEvent) {
	log.Println("[BatchListen] Event received")

	custArr, err := listen.custRepo.FindAllByCoreBankBatchId(event.BatchId)
	if err != nil {
		log.Println("[BatchListen] Error: " + err.Error())
	}

	for _, cust := range custArr {
		transcArr, err := listen.transcRepo.FindAllByCustomerCoreBankId(cust.CoreBankId)
		if err != nil {
			log.Println("[BatchListen] Error: " + err.Error())
		}

		for _, transc := range transcArr {
			err = listen.transcRepo.Delete(transc)
			if err != nil {
				log.Println("[BatchListen] Error: " + err.Error())
			}
		}
	}
}

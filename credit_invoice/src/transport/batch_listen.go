package transport

import (
	comm "dev-monteiro/go-digital-bank/commons"
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type BatchListen struct {
	sqsClnt    conn.SqsConn
	custRepo   data.CustomerRepo
	transcRepo data.TransactionRepo
}

func NewBatchListen(sqsClnt conn.SqsConn, custRepo data.CustomerRepo, transcRepo data.TransactionRepo) (*BatchListen, error) {
	queueName := os.Getenv("AWS_BATCHES_QUEUE_NAME")
	sqsUrlInput := sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}

	sqsUrlOutput, err := sqsClnt.GetQueueUrl(&sqsUrlInput)
	if err != nil {
		return nil, err
	}

	// TODO: make a pool of goroutines to handle more than 10 messages concurrently
	// TODO: create some safe exiting mechanism for the main goroutine
	go func() {
		for {
			sqsRecvInput := sqs.ReceiveMessageInput{
				QueueUrl:            sqsUrlOutput.QueueUrl,
				MaxNumberOfMessages: aws.Int64(10),
			}

			sqsRecvOutput, err := sqsClnt.ReceiveMessage(&sqsRecvInput)
			if err != nil {
				log.Println(err)
			}

			if len(sqsRecvOutput.Messages) == 0 {
				time.Sleep(1 * time.Second)
			}

			for _, msg := range sqsRecvOutput.Messages {
				if aws.StringValue(msg.Body) == "" {
					continue
				}
				log.Println("[BatchListen] Event received")

				batch := comm.BatchEvent{}
				json.Unmarshal([]byte(*msg.Body), &batch)

				custArr, err := custRepo.FindAllByCoreBankBatchId(batch.BatchId)
				if err != nil {
					log.Println("[BatchListen] " + err.Error())
				}

				for _, cust := range custArr {
					transcArr, err := transcRepo.FindAllByCustomerCoreBankId(cust.CoreBankId)
					if err != nil {
						log.Println("[BatchListen] " + err.Error())
					}

					for _, transc := range transcArr {
						err = transcRepo.Delete(transc)
						if err != nil {
							log.Println("[BatchListen] " + err.Error())
						}
					}
				}

				sqsDelInput := sqs.DeleteMessageInput{
					QueueUrl:      sqsUrlOutput.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				}

				_, err = sqsClnt.DeleteMessage(&sqsDelInput)
				if err != nil {
					log.Println("[BatchListen] " + err.Error())
				}
			}
		}
	}()

	return &BatchListen{
		sqsClnt:    sqsClnt,
		custRepo:   custRepo,
		transcRepo: transcRepo,
	}, nil
}

package transport

import (
	comm "devv-monteiro/go-digital-bank/commons"
	conn "devv-monteiro/go-digital-bank/credit-invoice/src/connector"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type PurchaseListen struct {
	sqsClnt    conn.SqsConn
	transcRepo data.TransactionRepo
}

func NewPurchaseListen(sqsClnt conn.SqsConn, transcRepo data.TransactionRepo) (*PurchaseListen, error) {
	queueName := os.Getenv("AWS_PURCHASES_QUEUE_NAME")
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
				log.Println("[PurchaseListen] Event received")

				purch := comm.PurchaseEvent{}
				json.Unmarshal([]byte(*msg.Body), &purch)

				transc := data.Transaction{
					PurchaseId:         purch.PurchaseId,
					CustomerCoreBankId: purch.CreditAccountId,
					Amount:             purch.Amount,
				}

				err := transcRepo.Save(transc)
				if err != nil {
					log.Println("[PurchaseListen] " + err.Error())
				}

				sqsDelInput := sqs.DeleteMessageInput{
					QueueUrl:      sqsUrlOutput.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				}

				_, err = sqsClnt.DeleteMessage(&sqsDelInput)
				if err != nil {
					log.Println("[PurchaseListen] " + err.Error())
				}
			}
		}
	}()

	return &PurchaseListen{
		sqsClnt:    sqsClnt,
		transcRepo: transcRepo,
	}, nil
}

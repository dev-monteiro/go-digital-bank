package transport

import (
	comm "devv-monteiro/go-digital-bank/commons"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type PurchaseListen struct {
	sqsClnt    *sqs.SQS
	transcRepo *data.TransactionRepo
}

func NewPurchaseListen(sqsClnt *sqs.SQS, transcRepo *data.TransactionRepo) (*PurchaseListen, error) {
	queueName := "purchases-queue"
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
				fmt.Println(err)
			}

			if len(sqsRecvOutput.Messages) == 0 {
				time.Sleep(1 * time.Second)
			}

			for _, msg := range sqsRecvOutput.Messages {
				if aws.StringValue(msg.Body) == "" {
					continue
				}

				purch := comm.PurchaseEvent{}
				json.Unmarshal([]byte(*msg.Body), &purch)

				transc := data.Transaction{
					PurchaseId:         purch.PurchaseId,
					CustomerCoreBankId: purch.CreditAccountId,
					Amount:             purch.Amount,
				}

				err := transcRepo.Save(transc)
				if err != nil {
					fmt.Println(err)
				}

				sqsDelInput := sqs.DeleteMessageInput{
					QueueUrl:      sqsUrlOutput.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				}

				_, err = sqsClnt.DeleteMessage(&sqsDelInput)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	return &PurchaseListen{
		sqsClnt:    sqsClnt,
		transcRepo: transcRepo,
	}, nil
}

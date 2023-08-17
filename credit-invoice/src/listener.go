package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Listener struct {
	sqsClient *sqs.SQS
	repo      PurchaseRepository
}

func NewListener(sqsClient *sqs.SQS, repo PurchaseRepository) Listener {
	queueName := "purchases-queue"

	queueUrlResp, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		fmt.Println(err)
	}

	// TODO: make a pool of goroutines to handle more than 10 messages concurrently
	// TODO: create some safe exiting mechanism for the main goroutine
	go func() {
		for {
			output, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            queueUrlResp.QueueUrl,
				MaxNumberOfMessages: aws.Int64(10),
			})
			if err != nil {
				fmt.Println(err)
			}

			if len(output.Messages) == 0 {
				time.Sleep(1 * time.Second)
			}

			for _, msg := range output.Messages {
				if aws.StringValue(msg.Body) != "" {
					pur := Purchase{}
					json.Unmarshal([]byte(*msg.Body), &pur)
					fmt.Println(pur)
					repo.save(pur)

					_, err = sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
						QueueUrl:      queueUrlResp.QueueUrl,
						ReceiptHandle: msg.ReceiptHandle,
					})
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}()

	fmt.Println("Listener set up.")
	return Listener{
		sqsClient: sqsClient,
		repo:      repo,
	}
}

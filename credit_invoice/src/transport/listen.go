package transport

import (
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Listen[T any] interface {
	getQueueNameEnv() string
	process(event T)
}

func SetupListen[T any](sqsConn conn.SqsConn, listen Listen[T]) (Listen[T], error) {
	queueName := os.Getenv(listen.getQueueNameEnv())
	sqsUrlInput := sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}

	sqsUrlOutput, err := sqsConn.GetQueueUrl(&sqsUrlInput)
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

			sqsRecvOutput, err := sqsConn.ReceiveMessage(&sqsRecvInput)
			if err != nil {
				log.Println("[BaseListen] " + err.Error())
			}

			if len(sqsRecvOutput.Messages) == 0 {
				time.Sleep(1 * time.Second)
			}

			for _, msg := range sqsRecvOutput.Messages {
				if aws.StringValue(msg.Body) == "" {
					continue
				}

				var event T
				json.Unmarshal([]byte(*msg.Body), &event)

				listen.process(event)

				sqsDelInput := sqs.DeleteMessageInput{
					QueueUrl:      sqsUrlOutput.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				}

				_, err = sqsConn.DeleteMessage(&sqsDelInput)
				if err != nil {
					log.Println("[BaseListen] " + err.Error())
				}
			}
		}
	}()

	return listen, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type PurchaseEvent struct {
	PurchaseId          int
	CreditAccountId     int
	PurchaseDateTime    string
	Amount              float32
	NumInstallments     int
	MerchantDescription string
	Status              string
}

func main() {
	time.Sleep(20 * time.Second)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String("http://localstack:4566"),
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connected to SQS.")
	}

	sqsClient := sqs.New(sess)
	queueName := "purchases-queue"

	queueUrlResp, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 50; i++ {
		time.Sleep(2 * time.Second)

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		event := PurchaseEvent{
			PurchaseId:          i,
			CreditAccountId:     123,
			PurchaseDateTime:    time.Now().String(),
			Amount:              float32(r.Intn(10000) / 100),
			NumInstallments:     1,
			MerchantDescription: "Acme Shop",
			Status:              "APPROVED",
		}
		body, err := json.Marshal(event)
		if err != nil {
			fmt.Println(err)
		}

		output, err := sqsClient.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String(string(body)),
			QueueUrl:    queueUrlResp.QueueUrl,
		})

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(output)
		}
	}
}

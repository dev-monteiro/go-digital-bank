package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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

var sqsClient *sqs.SQS
var queueUrl *string

func sendEvent(resWr http.ResponseWriter, req *http.Request) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	event := PurchaseEvent{
		PurchaseId:          r.Intn(10000),
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
		QueueUrl:    queueUrl,
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}

	resWr.WriteHeader(200)
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

	sqsClient = sqs.New(sess)

	queueName := "purchases-queue"
	queueUrlResp, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		fmt.Println(err)
	}
	queueUrl = queueUrlResp.QueueUrl

	http.HandleFunc("/", sendEvent)

	http.ListenAndServe(":80", nil)
}

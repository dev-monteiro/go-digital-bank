package main

import (
	"devv-monteiro/go-digital-bank/commons"
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

var sqsClient *sqs.SQS
var queueUrl *string
var random *rand.Rand

func main() {
	setup()

	http.HandleFunc("/", sendEvent)

	http.ListenAndServe(":80", nil)
}

func setup() {
	for {
		err := setupSqs()

		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
	}

	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Setup completed")
}

func setupSqs() error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String("http://localstack:4566"),
	})

	if err != nil {
		return err
	}
	sqsClient = sqs.New(sess)

	queueName := "purchases-queue"
	queueUrlResp, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		return err
	}
	queueUrl = queueUrlResp.QueueUrl

	return nil
}

func sendEvent(resWr http.ResponseWriter, req *http.Request) {
	event := commons.PurchaseEvent{
		PurchaseId:          random.Intn(10000),
		CreditAccountId:     123,
		PurchaseDateTime:    time.Now().String(),
		Amount:              float32(random.Intn(10000) / 100),
		NumInstallments:     1,
		MerchantDescription: "Acme Mall",
		Status:              "APPROVED",
		Description:         "love generation",
	}
	body, err := json.Marshal(event)

	if err != nil {
		fmt.Println(err)
		resWr.WriteHeader(500)
		return
	}

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    queueUrl,
	})

	if err != nil {
		fmt.Println(err)
		resWr.WriteHeader(500)
		return
	}

	resWr.WriteHeader(200)
}

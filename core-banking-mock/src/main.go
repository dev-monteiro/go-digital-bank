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

	fmt.Println("setup completed!")

	http.HandleFunc("/invoices", getInvoices)
	http.HandleFunc("/trigger/purchase", createPurchase)

	http.ListenAndServe(":80", nil)
}

func setup() {
	for {
		err := setupSqs()

		if err == nil {
			break
		}

		fmt.Println("setting up...")
		time.Sleep(5 * time.Second)
	}

	random = rand.New(rand.NewSource(time.Now().UnixNano()))
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

func getInvoices(resWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	creditAccountId := request.Form.Get("creditAccountId")

	fmt.Println("CreditAccountId: " + creditAccountId)
	if creditAccountId != "123" {
		resWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	invoice := commons.CoreBankInvoiceResp{
		CreditAccountId:     123,
		ProcessingSituation: "OPEN",
		IsPaymentDone:       false,
		DueDate:             genDueDate(5),
		ActualDueDate:       genDueDate(5),
		ClosingDate:         genClosingDate(30),
		TotalAmount:         1234.5,
		InvoiceId:           1234,
	}

	invoiceList := commons.CoreBankInvoiceListResp{
		Invoices: []commons.CoreBankInvoiceResp{invoice},
	}

	resWriter.Header().Add("Content-Type", "application/json")
	json.NewEncoder(resWriter).Encode(invoiceList)
}

func genClosingDate(closingDay int) string {
	closingDate := time.Date(time.Now().Year(), time.Now().Month(), closingDay, 0, 0, 0, 0, time.UTC)
	return closingDate.Format("2006-01-02")
}

func genDueDate(dueDay int) string {
	dueDate := time.Date(time.Now().Year(), time.Now().Month()+1, dueDay, 0, 0, 0, 0, time.UTC)
	return dueDate.Format("2006-01-02")
}

func createPurchase(resWr http.ResponseWriter, req *http.Request) {
	event := commons.PurchaseEvent{
		PurchaseId:          random.Intn(10000),
		CreditAccountId:     123,
		PurchaseDateTime:    time.Now().String(),
		Amount:              float64(random.Intn(10000) / 100),
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

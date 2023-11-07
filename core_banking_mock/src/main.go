package main

import (
	"dev-monteiro/go-digital-bank/commons"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
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

	log.Println("[Mock] Setup completed!")

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

		log.Println("[Mock] Setting up...")
		time.Sleep(5 * time.Second)
	}

	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func setupSqs() error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_LOGIN"), os.Getenv("AWS_PASS"), ""),
		Endpoint:    aws.String(os.Getenv("AWS_ENDPOINT")),
	})

	if err != nil {
		return err
	}
	sqsClient = sqs.New(sess)

	queueName := os.Getenv("AWS_PURCHASES_QUEUE_NAME")
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
	log.Println("[Mock] GetInvoices")

	request.ParseForm()
	creditAccountId := request.Form.Get("creditAccountId")

	log.Println("[Mock] CreditAccountId: " + creditAccountId)
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
	log.Println("[Mock] CreatePurchase")

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
		log.Println(err)
		resWr.WriteHeader(500)
		return
	}

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    queueUrl,
	})

	if err != nil {
		log.Println(err)
		resWr.WriteHeader(500)
		return
	}

	resWr.WriteHeader(200)
}

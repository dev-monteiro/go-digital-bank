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

type MockData struct {
	sqsClient   *sqs.SQS
	queueUrlMap map[string]*string
	random      *rand.Rand
}

var data *MockData

func main() {
	setup()

	log.Println("[Mock] Setup completed!")

	http.HandleFunc("/invoices", getInvoices)
	http.HandleFunc("/trigger/purchase", createPurchase)
	http.HandleFunc("/trigger/batch", createBatch)

	http.ListenAndServe(":80", nil)
}

func setup() {
	data = &MockData{queueUrlMap: make(map[string]*string)}

	for {
		err := setupSqs()

		if err == nil {
			break
		}

		log.Println("[Mock] Setting up...")
		time.Sleep(5 * time.Second)
	}

	data.random = rand.New(rand.NewSource(time.Now().UnixNano()))
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
	data.sqsClient = sqs.New(sess)

	err = setupQueue("AWS_PURCHASES_QUEUE_NAME")
	if err != nil {
		return err
	}

	err = setupQueue("AWS_BATCHES_QUEUE_NAME")
	if err != nil {
		return err
	}

	return nil
}

func setupQueue(queueNameEnv string) error {
	queueName := os.Getenv(queueNameEnv)
	queueUrlResp, err := data.sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		return err
	}

	data.queueUrlMap[queueNameEnv] = queueUrlResp.QueueUrl

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
	return closingDate.Format(time.DateOnly)
}

func genDueDate(dueDay int) string {
	dueDate := time.Date(time.Now().Year(), time.Now().Month()+1, dueDay, 0, 0, 0, 0, time.UTC)
	return dueDate.Format(time.DateOnly)
}

func createPurchase(resWr http.ResponseWriter, req *http.Request) {
	log.Println("[Mock] CreatePurchase")

	event := commons.PurchaseEvent{
		PurchaseId:          data.random.Intn(10000),
		CreditAccountId:     123,
		PurchaseDateTime:    time.Now().String(),
		Amount:              float64(data.random.Intn(10000) / 100),
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

	_, err = data.sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    data.queueUrlMap["AWS_PURCHASES_QUEUE_NAME"],
	})

	if err != nil {
		log.Println(err)
		resWr.WriteHeader(500)
		return
	}

	resWr.WriteHeader(200)
}

func createBatch(resWr http.ResponseWriter, req *http.Request) {
	log.Println("[Mock] CreateBatch")

	event := commons.BatchEvent{
		BatchId:       789,
		ReferenceDate: time.Now().Format(time.DateOnly),
	}
	body, err := json.Marshal(event)

	if err != nil {
		log.Println(err)
		resWr.WriteHeader(500)
		return
	}

	_, err = data.sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    data.queueUrlMap["AWS_BATCHES_QUEUE_NAME"],
	})

	if err != nil {
		log.Println(err)
		resWr.WriteHeader(500)
		return
	}

	resWr.WriteHeader(200)
}

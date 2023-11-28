package main

import (
	comm "dev-monteiro/go-digital-bank/commons"
	"dev-monteiro/go-digital-bank/commons/invstat"
	"dev-monteiro/go-digital-bank/commons/ldate"
	"dev-monteiro/go-digital-bank/commons/mnyamnt"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MockData struct {
	sqsClient            *sqs.SQS
	queueUrlMap          map[string]*string
	random               *rand.Rand
	pendingPurchaseMap   map[int][]comm.PurchaseEvent
	processedPurchaseMap map[int][]comm.PurchaseEvent
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
	data = &MockData{
		queueUrlMap:          make(map[string]*string),
		pendingPurchaseMap:   make(map[int][]comm.PurchaseEvent),
		processedPurchaseMap: make(map[int][]comm.PurchaseEvent),
	}

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

	log.Println("[Mock] GetInvoices CreditAccountId: " + creditAccountId)
	if creditAccountId != "123" {
		resWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	totalAmount := mnyamnt.NewMnyAmount("1234.5")
	for _, purchase := range data.processedPurchaseMap[123] {
		totalAmount = totalAmount.Add(purchase.Amount)
	}

	invoice := comm.CoreBankInvoiceResp{
		CustomerId:    123,
		Status:        invstat.OPEN,
		IsPaymentDone: false,
		DueDate:       genDueDate(5),
		ActualDueDate: genDueDate(5),
		ClosingDate:   genClosingDate(30),
		Amount:        totalAmount,
	}

	invoiceList := comm.CoreBankInvoiceListResp{
		Invoices: []comm.CoreBankInvoiceResp{invoice},
	}

	resWriter.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(resWriter).Encode(invoiceList)
	if err != nil {
		log.Println(err)
	}
}

func genDueDate(dueDay int) *ldate.LocDate {
	return ldate.NewLocDate(ldate.Today().Year(), ldate.Today().Month()+1, dueDay)
}

func genClosingDate(closingDay int) *ldate.LocDate {
	return ldate.NewLocDate(ldate.Today().Year(), ldate.Today().Month(), closingDay)
}

func createPurchase(resWr http.ResponseWriter, req *http.Request) {
	log.Println("[Mock] CreatePurchase")

	purchase := comm.PurchaseEvent{
		Id:                  data.random.Intn(10000),
		CustomerId:          123,
		DateTime:            time.Now().String(),
		Amount:              mnyamnt.NewMnyAmount(strconv.FormatFloat(float64(data.random.Intn(10000))/100.0, 'f', 2, 64)),
		NumInstallments:     1,
		MerchantDescription: "Acme Mall",
		Status:              "APPROVED",
	}
	body, err := json.Marshal(purchase)

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

	data.pendingPurchaseMap[123] = append(data.pendingPurchaseMap[123], purchase)

	resWr.WriteHeader(200)
}

func createBatch(resWr http.ResponseWriter, req *http.Request) {
	log.Println("[Mock] CreateBatch")

	event := comm.BatchEvent{
		Id:            789,
		ReferenceDate: ldate.Today(),
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

	data.processedPurchaseMap[123] = append(data.processedPurchaseMap[123], data.pendingPurchaseMap[123]...)
	data.pendingPurchaseMap[123] = []comm.PurchaseEvent{}

	resWr.WriteHeader(200)
}

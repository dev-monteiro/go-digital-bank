package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type CurrentInvoiceResponse struct {
	Id          string
	StatusLabel string
	Amount      float32
	ClosingDate string
}

type CoreBankingInvoiceResponse struct {
	CreditAccountId     int32
	ProcessingSituation string
	IsPaymentDone       bool
	DueDate             string
	ActualDueDate       string
	ClosingDate         string
	TotalAmount         float32
	InvoiceId           int
}

type CoreBankingInvoiceListResponse struct {
	Invoices []CoreBankingInvoiceResponse
}

func getCurrentInvoice(resWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	customerId := request.Form.Get("customerId")
	fmt.Println("CustomerId = " + customerId)

	query := "select credit_account_id from customer_credentials where customer_id = '" + customerId + "';"
	fmt.Println("Query = " + query)

	row := db.QueryRow(query)
	var creditAccountId int
	err := row.Scan(&creditAccountId)

	if err != nil {
		resWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	url := "http://core_banking_mock/invoices?creditAccountId=" + strconv.Itoa(creditAccountId)
	fmt.Println("Url = " + url)

	cbListResWrapper, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		resWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cbListRes CoreBankingInvoiceListResponse
	json.NewDecoder(cbListResWrapper.Body).Decode(&cbListRes)
	cbRes := cbListRes.Invoices[0]

	response := CurrentInvoiceResponse{
		Id:          strconv.Itoa(cbRes.InvoiceId),
		StatusLabel: cbRes.ProcessingSituation,
		Amount:      cbRes.TotalAmount,
		ClosingDate: cbRes.ClosingDate,
	}

	resWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resWriter).Encode(response)
}

func init() {
	time.Sleep(2 * time.Second)

	var err error
	db, err = sql.Open("mysql", "root:root@tcp(mysql)/credit_invoice")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Database connected.")
	}
}

func main() {
	http.HandleFunc("/invoices/current", getCurrentInvoice)

	http.ListenAndServe(":80", nil)
}

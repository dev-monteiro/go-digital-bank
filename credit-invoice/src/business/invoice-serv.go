package business

import (
	"devv-monteiro/go-digital-bank/commons"
	comm "devv-monteiro/go-digital-bank/commons"
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"encoding/json"
	"net/http"
	"strconv"
)

type InvoiceServ struct {
	custRepo   *data.CustomerRepo
	transcRepo *data.TransactionRepo
}

func NewInvoiceServ(custRepo *data.CustomerRepo, transcRepo *data.TransactionRepo) *InvoiceServ {
	return &InvoiceServ{
		custRepo:   custRepo,
		transcRepo: transcRepo,
	}
}

func (serv *InvoiceServ) GetCurrInvoice(custId string) (*CurrInvoiceResp, *conf.AppError) {
	log.Println("[InvoiceServ] GetCurrInvoice")

	cust, err := serv.custRepo.FindById(custId)
	if err != nil {
		return nil, err
	}

	invoArr, err := serv.getCoreBankInvoices(cust.CoreBankId)
	if err != nil {
		return nil, err
	}

	invo, err := serv.getCurrInvoice(invoArr)
	if err != nil {
		return nil, err
	}

	amount := invo.TotalAmount
	if invo.ProcessingSituation == "OPEN" {
		amount, err = serv.updateInvoiceAmount(cust.CoreBankId, amount)
		if err != nil {
			return nil, err
		}
	}

	closDate, err := serv.convertClosingDate(invo.ClosingDate)
	if err != nil {
		return nil, err
	}

	resp := CurrInvoiceResp{
		StatusLabel: strings.Title(strings.ToLower(invo.ProcessingSituation)),
		Amount:      fmt.Sprintf("$ %.2f", amount),
		ClosingDate: closDate,
	}

	return &resp, nil
}

func (InvoiceServ) getCoreBankInvoices(custCoreBankId int) ([]comm.CoreBankInvoiceResp, *conf.AppError) {
	log.Println("[InvoiceServ] GetCoreBankInvoices")

	url := "http://" + os.Getenv("CORE_BANKING_HOST") + "/invoices"
	url = url + "?creditAccountId=" + strconv.Itoa(custCoreBankId)

	resp, err := http.Get(url)
	if err != nil {
		return nil, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	var invoListResp comm.CoreBankInvoiceListResp
	json.NewDecoder(resp.Body).Decode(&invoListResp)
	return invoListResp.Invoices, nil
}

func (serv *InvoiceServ) getCurrInvoice(invoArr []commons.CoreBankInvoiceResp) (*commons.CoreBankInvoiceResp, *conf.AppError) {
	log.Println("[InvoiceServ] GetCurrInvoice")

	var openInvo commons.CoreBankInvoiceResp

	for _, invo := range invoArr {
		dueDate, err := time.Parse(time.DateOnly, invo.ActualDueDate)
		if err != nil {
			return nil, &conf.AppError{
				Message:    "Unknown error: " + err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}

		if invo.ProcessingSituation == "CLOSED" && !dueDate.Before(time.Now().Truncate(24*time.Hour)) {
			return &invo, nil
		} else if invo.ProcessingSituation == "OPEN" {
			openInvo = invo
		}
	}

	return &openInvo, nil
}
func (serv *InvoiceServ) updateInvoiceAmount(custCoreBankId int, invoAmount float64) (float64, *conf.AppError) {
	log.Println("[InvoiceServ] UpdateInvoiceAmount")

	transcArr, err := serv.transcRepo.FindAllByCustomerCoreBankId(custCoreBankId)

	if err != nil {
		return 0, err
	}

	for _, transc := range transcArr {
		invoAmount += transc.Amount
	}
	return invoAmount, nil
}

func (InvoiceServ) convertClosingDate(closingDate string) (string, *conf.AppError) {
	log.Println("[InvoiceServ] ConvertClosingDate")

	parsedDate, err := time.Parse(time.DateOnly, closingDate)
	if err != nil {
		return "", &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return strings.ToUpper(parsedDate.Format("Jan 02")), nil
}

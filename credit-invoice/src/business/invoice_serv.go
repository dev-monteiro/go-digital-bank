package business

import (
	"devv-monteiro/go-digital-bank/commons"
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	conn "devv-monteiro/go-digital-bank/credit-invoice/src/connector"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type InvoiceServ interface {
	GetCurrInvoice(custId string) (*CurrInvoiceResp, *conf.AppError)
}

type invoiceServ struct {
	custRepo     data.CustomerRepo
	transcRepo   data.TransactionRepo
	coreBankConn conn.CoreBankConn
}

func NewInvoiceServ(custRepo data.CustomerRepo, transcRepo data.TransactionRepo, coreBankConn conn.CoreBankConn) InvoiceServ {
	return &invoiceServ{
		custRepo:     custRepo,
		transcRepo:   transcRepo,
		coreBankConn: coreBankConn,
	}
}

func (serv *invoiceServ) GetCurrInvoice(custId string) (*CurrInvoiceResp, *conf.AppError) {
	log.Println("[InvoiceServ] GetCurrInvoice")

	cust, err := serv.custRepo.FindById(custId)
	if cust == nil {
		return nil, &conf.AppError{Message: "Customer not found.", StatusCode: http.StatusNotFound}
	}
	if err != nil {
		return nil, conf.NewUnknownError(err)
	}

	invoArr, err := serv.coreBankConn.GetAllInvoices(cust.CoreBankId)
	if err != nil {
		return nil, conf.NewUnknownError(err) // TODO: improve error handling
	}

	invo, err := serv.getCurrInvoice(invoArr)
	if err != nil {
		return nil, conf.NewUnknownError(err)
	}

	amount := invo.TotalAmount
	if invo.ProcessingSituation == "OPEN" {
		amount, err = serv.updateInvoiceAmount(cust.CoreBankId, amount)
		if err != nil {
			return nil, conf.NewUnknownError(err)
		}
	}

	closDate, err := serv.convertClosingDate(invo.ClosingDate)
	if err != nil {
		return nil, conf.NewUnknownError(err)
	}

	resp := CurrInvoiceResp{
		StatusLabel: strings.Title(strings.ToLower(invo.ProcessingSituation)),
		Amount:      fmt.Sprintf("$ %.2f", amount),
		ClosingDate: closDate,
	}

	return &resp, nil
}

func (serv *invoiceServ) getCurrInvoice(invoArr []commons.CoreBankInvoiceResp) (*commons.CoreBankInvoiceResp, error) {
	log.Println("[InvoiceServ] GetCurrInvoice")

	var openInvo commons.CoreBankInvoiceResp

	for _, invo := range invoArr {
		dueDate, err := time.Parse(time.DateOnly, invo.ActualDueDate)
		if err != nil {
			return nil, err
		}

		if invo.ProcessingSituation == "CLOSED" && !dueDate.Before(time.Now().Truncate(24*time.Hour)) {
			return &invo, nil
		} else if invo.ProcessingSituation == "OPEN" {
			openInvo = invo
		}
	}

	return &openInvo, nil
}
func (serv *invoiceServ) updateInvoiceAmount(custCoreBankId int, invoAmount float64) (float64, error) {
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

func (serv *invoiceServ) convertClosingDate(closingDate string) (string, error) {
	log.Println("[InvoiceServ] ConvertClosingDate")

	parsedDate, err := time.Parse(time.DateOnly, closingDate)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(parsedDate.Format("Jan 02")), nil
}

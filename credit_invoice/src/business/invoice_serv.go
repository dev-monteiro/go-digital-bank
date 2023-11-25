package business

import (
	comm "dev-monteiro/go-digital-bank/commons"
	conf "dev-monteiro/go-digital-bank/credit-invoice/src/configuration"
	conn "dev-monteiro/go-digital-bank/credit-invoice/src/connector"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	"log"
	"net/http"
	"strings"
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
	if err != nil {
		return nil, conf.NewUnknownError(err)
	}
	if cust == nil {
		return nil, &conf.AppError{Message: conf.CUSTOMER_NOT_FOUND, StatusCode: http.StatusNotFound}
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

	resp := CurrInvoiceResp{
		StatusLabel: strings.Title(strings.ToLower(invo.ProcessingSituation)),
		Amount:      "$ " + amount.String(),
		ClosingDate: strings.ToUpper(invo.ClosingDate.Format(comm.MonLitCapsDayNum)),
	}

	return &resp, nil
}

func (serv *invoiceServ) getCurrInvoice(invoArr []comm.CoreBankInvoiceResp) (*comm.CoreBankInvoiceResp, error) {
	log.Println("[InvoiceServ] GetCurrInvoice")

	var openInvo comm.CoreBankInvoiceResp

	for _, invo := range invoArr {
		if invo.ProcessingSituation == "CLOSED" && !comm.Today().After(invo.ActualDueDate) {
			return &invo, nil
		} else if invo.ProcessingSituation == "OPEN" {
			openInvo = invo
		}
	}

	return &openInvo, nil
}
func (serv *invoiceServ) updateInvoiceAmount(custCoreBankId int, invoAmount *comm.MoneyAmount) (*comm.MoneyAmount, error) {
	log.Println("[InvoiceServ] UpdateInvoiceAmount")

	transcArr, err := serv.transcRepo.FindAllByCustomerCoreBankId(custCoreBankId)

	if err != nil {
		return nil, err
	}

	sum := invoAmount
	for _, transc := range transcArr {
		sum = sum.Add(transc.Amount)
	}

	return sum, nil
}

package business

import (
	comm "dev-monteiro/go-digital-bank/commons"
	"dev-monteiro/go-digital-bank/commons/invstat"
	"dev-monteiro/go-digital-bank/commons/ldate"
	"dev-monteiro/go-digital-bank/commons/mnyamnt"
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

	amount := invo.Amount
	if invo.Status == invstat.OPEN {
		amount, err = serv.updateInvoiceAmount(cust.CoreBankId, amount)
		if err != nil {
			return nil, conf.NewUnknownError(err)
		}
	}

	resp := CurrInvoiceResp{
		StatusLabel:    strings.Title(strings.ToLower(string(invo.Status))),
		Amount:         "$ " + amount.String(),
		FmtClosingDate: strings.ToUpper(invo.ClosingDate.Format(ldate.MMdd)),
	}

	return &resp, nil
}

func (serv *invoiceServ) getCurrInvoice(invoArr []comm.CoreBankInvoiceResp) (*comm.CoreBankInvoiceResp, error) {
	log.Println("[InvoiceServ] GetCurrInvoice")

	var openInvo comm.CoreBankInvoiceResp

	for _, invo := range invoArr {
		if invo.Status == invstat.CLOSED && !ldate.Today().After(invo.ActualDueDate) {
			return &invo, nil
		} else if invo.Status == invstat.OPEN {
			openInvo = invo
		}
	}

	return &openInvo, nil
}
func (serv *invoiceServ) updateInvoiceAmount(custCoreBankId int, invoAmount *mnyamnt.MnyAmount) (*mnyamnt.MnyAmount, error) {
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

package business

import (
	comm "devv-monteiro/go-digital-bank/commons"
	conf "devv-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "devv-monteiro/go-digital-bank/credit-invoice/src/database"
	"fmt"

	"encoding/json"
	"net/http"
	"strconv"
)

type InvoiceServ struct {
	credRepo  *data.CustomerRepo
	purchRepo *data.PurchaseRepo
	invoRepo  *data.InvoiceRepo
}

func NewInvoiceServ(custRepo *data.CustomerRepo, purchRepo *data.PurchaseRepo, invoRepo *data.InvoiceRepo) *InvoiceServ {
	return &InvoiceServ{
		credRepo:  custRepo,
		purchRepo: purchRepo,
		invoRepo:  invoRepo,
	}
}

func (serv *InvoiceServ) GetCurrInvoice(custId string) (*CurrInvoiceResp, *conf.AppError) {
	cbCustId, err := serv.credRepo.GetCoreBankId(custId)
	if err != nil {
		return nil, err
	}

	cbInvo, err := serv.getCoreBankInvoice(cbCustId)
	if err != nil {
		return nil, err
	}

	id, err := serv.invoRepo.GetId(cbInvo.InvoiceId)
	if err != nil {
		return nil, err
	}
	if id == "" {
		invo := data.NewInvoice(cbInvo.InvoiceId)

		err = serv.invoRepo.Save(*invo)
		if err != nil {
			return nil, err
		}

		id = invo.Id
	}

	updaAmount, err := serv.updateAmount(cbCustId, cbInvo.TotalAmount)
	if err != nil {
		return nil, err
	}

	resp := CurrInvoiceResp{
		Id:          id,
		StatusLabel: cbInvo.ProcessingSituation,
		Amount:      updaAmount,
		ClosingDate: cbInvo.ClosingDate,
	}

	return &resp, nil
}

func (serv *InvoiceServ) updateAmount(creditAccountId int, currentAmount float32) (float32, *conf.AppError) {
	purchases, err := serv.purchRepo.FindAllByCreditAccountId(creditAccountId)

	if err != nil {
		return 0, err
	}

	for _, purchase := range purchases {
		currentAmount += purchase.Amount
	}
	return currentAmount, nil
}

func (InvoiceServ) getCoreBankInvoice(creditAccountId int) (*comm.CoreBankInvoiceResp, *conf.AppError) {
	fmt.Println("getCoreBankInvoice")

	url := "http://core_banking_mock/invoices?creditAccountId=" + strconv.Itoa(creditAccountId)

	coreBankResp, err := http.Get(url)
	if err != nil {
		return nil, &conf.AppError{
			Message:    "Unknown error: " + err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	var invoiceList comm.CoreBankInvoiceListResp
	json.NewDecoder(coreBankResp.Body).Decode(&invoiceList)
	invoice := invoiceList.Invoices[0]
	return &invoice, nil
}

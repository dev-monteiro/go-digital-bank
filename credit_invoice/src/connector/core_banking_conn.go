package connector

import (
	comm "dev-monteiro/go-digital-bank/commons"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type CoreBankConn interface {
	GetAllInvoices(custCoreBankId int) ([]comm.CoreBankInvoiceResp, error)
}

type coreBankConn struct {
	coreBankHost string
}

func NewCoreBankConn() CoreBankConn {
	return &coreBankConn{coreBankHost: os.Getenv("CORE_BANKING_HOST")}
}

func (conn *coreBankConn) GetAllInvoices(custCoreBankId int) ([]comm.CoreBankInvoiceResp, error) {
	log.Println("[InvoiceServ] GetCoreBankInvoices")

	url := conn.coreBankHost + "/invoices"
	url = url + "?creditAccountId=" + strconv.Itoa(custCoreBankId)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var invoListResp comm.CoreBankInvoiceListResp
	json.NewDecoder(resp.Body).Decode(&invoListResp)
	return invoListResp.Invoices, nil
}

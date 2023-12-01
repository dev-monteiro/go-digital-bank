package system_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type currInvoiceResp struct {
	Status      string
	Amount      string
	ClosingDate string
}

type purchaseReq struct {
	CustId int
	Amount float64
}

func TestSystem_NoAction(t *testing.T) {
	invoResp, err := getCurrInvoice("abc-123-def")
	if err != nil {
		fmt.Println(err)
	}

	require.Equal(t, *invoResp, currInvoiceResp{
		Status:      "Open",
		Amount:      "$ 1234.50",
		ClosingDate: "DEC 30",
	})
}

func TestSystem_TriggerOnePurchase(t *testing.T) {
	err := triggerPurchase(123, 11.11)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(1 * time.Second)

	invoResp, err := getCurrInvoice("abc-123-def")
	if err != nil {
		fmt.Println(err)
	}

	require.Equal(t, *invoResp, currInvoiceResp{
		Status:      "Open",
		Amount:      "$ 1245.61",
		ClosingDate: "DEC 30",
	})
}

func getCurrInvoice(custId string) (*currInvoiceResp, error) {
	url := "http://localhost:8080/invoices/current?customerId=" + custId

	httpResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var invoResp currInvoiceResp
	json.NewDecoder(httpResp.Body).Decode(&invoResp)

	return &invoResp, nil
}

func triggerPurchase(cbCustId int, amount float64) error {
	url := "http://localhost:9090/trigger/purchase"

	payload, err := json.Marshal(purchaseReq{
		CustId: cbCustId,
		Amount: amount,
	})
	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	return nil
}

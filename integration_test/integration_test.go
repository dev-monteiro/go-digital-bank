package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type currInvoiceResp struct {
	Amount string
}

type purchaseReq struct {
	CustId int
	Amount float64
}

type coreBankInvoiceResp struct {
	TotalAmount float64
}

type coreBankInvoiceListResp struct {
	Content []coreBankInvoiceResp
}

func TestSystem(t *testing.T) {
	fmt.Println("----------------------------")

	for {
		time.Sleep(5 * time.Second)
		if isSystemReady() {
			break
		} else {
			fmt.Println("waiting for system to be ready...")
		}
	}

	fmt.Println("----------------------------")
	fmt.Println("Init")

	cbInvoResp := getCoreBankInvoice(t, 123)
	require.Equal(t, *cbInvoResp, coreBankInvoiceResp{
		TotalAmount: 1234.50,
	})

	invoResp := getCurrInvoice(t, "abc-123-def")
	require.Equal(t, *invoResp, currInvoiceResp{
		Amount: "$ 1234.50",
	})

	fmt.Println("OK")

	fmt.Println("----------------------------")
	fmt.Println("Trigger one purchase")

	triggerPurchase(t, 123, 11.11)

	cbInvoResp = getCoreBankInvoice(t, 123)
	require.Equal(t, *cbInvoResp, coreBankInvoiceResp{
		TotalAmount: 1234.50,
	})

	invoResp = getCurrInvoice(t, "abc-123-def")
	require.Equal(t, *invoResp, currInvoiceResp{
		Amount: "$ 1245.61",
	})

	fmt.Println("OK")

	fmt.Println("----------------------------")
	fmt.Println("Trigger multiple purchases")

	triggerPurchase(t, 123, 22.22)
	triggerPurchase(t, 123, 33.33)
	triggerPurchase(t, 123, 44.44)

	cbInvoResp = getCoreBankInvoice(t, 123)
	require.Equal(t, *cbInvoResp, coreBankInvoiceResp{
		TotalAmount: 1234.50,
	})

	invoResp = getCurrInvoice(t, "abc-123-def")
	require.Equal(t, *invoResp, currInvoiceResp{
		Amount: "$ 1345.60",
	})

	fmt.Println("OK")

	fmt.Println("----------------------------")
	fmt.Println("Trigger batch")

	triggerBatch(t)

	cbInvoResp = getCoreBankInvoice(t, 123)
	require.Equal(t, *cbInvoResp, coreBankInvoiceResp{
		TotalAmount: 1345.60,
	})

	invoResp = getCurrInvoice(t, "abc-123-def")
	require.Equal(t, *invoResp, currInvoiceResp{
		Amount: "$ 1345.60",
	})

	fmt.Println("OK")
}

func isSystemReady() bool {
	url := "http://localhost:8080/health"

	httpResp, err := http.Get(url)
	if err != nil {
		return false
	}

	return httpResp.StatusCode == 200
}

func getCoreBankInvoice(t *testing.T, cbCustId int) *coreBankInvoiceResp {
	url := "http://localhost:9090/invoices?creditAccountId=" + strconv.Itoa(cbCustId)

	httpResp, err := http.Get(url)
	if err != nil {
		onError(t, err)
		return nil
	}

	var cbInvoListResp coreBankInvoiceListResp
	json.NewDecoder(httpResp.Body).Decode(&cbInvoListResp)

	return &cbInvoListResp.Content[0]
}

func getCurrInvoice(t *testing.T, custId string) *currInvoiceResp {
	url := "http://localhost:8080/invoices/current?customerId=" + custId

	httpResp, err := http.Get(url)
	if err != nil {
		onError(t, err)
		return nil
	}

	var invoResp currInvoiceResp
	json.NewDecoder(httpResp.Body).Decode(&invoResp)

	return &invoResp
}

func triggerPurchase(t *testing.T, cbCustId int, amount float64) {
	url := "http://localhost:9090/trigger/purchase"

	payload, err := json.Marshal(purchaseReq{
		CustId: cbCustId,
		Amount: amount,
	})
	if err != nil {
		onError(t, err)
		return
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		onError(t, err)
		return
	}

	time.Sleep(1 * time.Second)
}

func triggerBatch(t *testing.T) {
	url := "http://localhost:9090/trigger/batch"

	_, err := http.Get(url)
	if err != nil {
		onError(t, err)
		return
	}

	time.Sleep(2 * time.Second)
}

func onError(t *testing.T, err error) {
	fmt.Println(err)
	t.FailNow()
}

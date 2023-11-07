package business_test

import (
	comm "dev-monteiro/go-digital-bank/commons"
	mock_imp "dev-monteiro/go-digital-bank/credit-invoice/mock"
	busn "dev-monteiro/go-digital-bank/credit-invoice/src/business"
	conf "dev-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	"errors"
	"net/http"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testData struct {
	Name               string
	Input              string
	CustRepoOutput     *data.Customer
	CustRepoError      error
	CoreBankConnOutput []comm.CoreBankInvoiceResp
	CoreBankConnError  error
	TranscRepoOutput   []data.Transaction
	TranscRepoError    error
	ExpectedOutput     *busn.CurrInvoiceResp
	ExpectedError      *conf.AppError
}

func TestInvoiceServ(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2023, 10, 31, 12, 0, 0, 0, time.UTC)
	})

	tests := []testData{
		{
			Name:               "should_ReturnNotFoundError_When_CustomerIsNotFound",
			Input:              "abc-123-def",
			CustRepoOutput:     nil,
			CustRepoError:      nil,
			CoreBankConnOutput: nil,
			CoreBankConnError:  nil,
			TranscRepoOutput:   nil,
			TranscRepoError:    nil,
			ExpectedOutput:     nil,
			ExpectedError: &conf.AppError{
				Message:    conf.CUSTOMER_NOT_FOUND,
				StatusCode: http.StatusNotFound,
			},
		},
		{
			Name:               "should_ReturnUnknownError_When_CustomerRepoReturnsError",
			Input:              "abc-123-def",
			CustRepoOutput:     nil,
			CustRepoError:      errors.New("error"),
			CoreBankConnOutput: nil,
			CoreBankConnError:  nil,
			TranscRepoOutput:   nil,
			TranscRepoError:    nil,
			ExpectedOutput:     nil,
			ExpectedError: &conf.AppError{
				Message:    "Unknown error: error",
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			Name:               "should_ReturnUnknownError_When_CoreBankConnReturnsError",
			Input:              "abc-123-def",
			CustRepoOutput:     &data.Customer{CoreBankId: 123},
			CustRepoError:      nil,
			CoreBankConnOutput: nil,
			CoreBankConnError:  errors.New("error"),
			TranscRepoOutput:   nil,
			TranscRepoError:    nil,
			ExpectedOutput:     nil,
			ExpectedError: &conf.AppError{
				Message:    "Unknown error: error",
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			Name:           "should_ReturnClosedInvoice_When_ThereIsValidClosedInvoice",
			Input:          "abc-123-def",
			CustRepoOutput: &data.Customer{CoreBankId: 123},
			CustRepoError:  nil,
			CoreBankConnOutput: []comm.CoreBankInvoiceResp{
				{
					ProcessingSituation: "CLOSED",
					ClosingDate:         date(2023, 9, 30),
					ActualDueDate:       date(2023, 10, 5),
					TotalAmount:         999.99,
				},
				{
					ProcessingSituation: "CLOSED",
					ClosingDate:         date(2023, 10, 30),
					ActualDueDate:       date(2023, 11, 6),
					TotalAmount:         456.78,
				},
				{
					ProcessingSituation: "OPEN",
					ClosingDate:         date(2023, 11, 30),
					ActualDueDate:       date(2023, 12, 5),
					TotalAmount:         123.45,
				},
			},
			CoreBankConnError: nil,
			TranscRepoOutput:  nil,
			TranscRepoError:   nil,
			ExpectedOutput: &busn.CurrInvoiceResp{
				StatusLabel: "Closed",
				Amount:      "$ 456.78",
				ClosingDate: "OCT 30",
			},
			ExpectedError: nil,
		},
		{
			Name:           "should_ReturnUnknownError_When_ThereIsNotValidClosedInvoice_And_TransactionRepoReturnsError",
			Input:          "abc-123-def",
			CustRepoOutput: &data.Customer{CoreBankId: 123},
			CustRepoError:  nil,
			CoreBankConnOutput: []comm.CoreBankInvoiceResp{
				{
					ProcessingSituation: "CLOSED",
					ClosingDate:         date(2023, 9, 10),
					ActualDueDate:       date(2023, 9, 15),
					TotalAmount:         888.88,
				},
				{
					ProcessingSituation: "CLOSED",
					ClosingDate:         date(2023, 10, 10),
					ActualDueDate:       date(2023, 10, 16),
					TotalAmount:         777.77,
				},
				{
					ProcessingSituation: "OPEN",
					ClosingDate:         date(2023, 11, 10),
					ActualDueDate:       date(2023, 11, 15),
					TotalAmount:         100,
				},
			},
			CoreBankConnError: nil,
			TranscRepoOutput:  nil,
			TranscRepoError:   errors.New("error"),
			ExpectedOutput:    nil,
			ExpectedError: &conf.AppError{
				Message:    "Unknown error: error",
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			Name:           "should_ReturnOpenInvoice_With_OriginalAmount_When_ThereIsNotValidClosedInvoice_And_ThereAreNoPendingTransactions",
			Input:          "abc-123-def",
			CustRepoOutput: &data.Customer{CoreBankId: 123},
			CustRepoError:  nil,
			CoreBankConnOutput: []comm.CoreBankInvoiceResp{
				{
					ProcessingSituation: "CLOSED",
					ClosingDate:         date(2023, 9, 10),
					ActualDueDate:       date(2023, 9, 15),
					TotalAmount:         888.88,
				},
				{
					ProcessingSituation: "CLOSED",
					ClosingDate:         date(2023, 10, 10),
					ActualDueDate:       date(2023, 10, 16),
					TotalAmount:         777.77,
				},
				{
					ProcessingSituation: "OPEN",
					ClosingDate:         date(2023, 11, 10),
					ActualDueDate:       date(2023, 11, 15),
					TotalAmount:         100,
				},
			},
			CoreBankConnError: nil,
			TranscRepoOutput:  []data.Transaction{},
			TranscRepoError:   nil,
			ExpectedOutput: &busn.CurrInvoiceResp{
				StatusLabel: "Open",
				Amount:      "$ 100.00",
				ClosingDate: "NOV 10",
			},
			ExpectedError: nil,
		},
		{
			Name:           "should_ReturnOpenInvoice_With_UpdatedAmount_When_ThereIsNotValidClosedInvoice_And_ThereArePendingTransactions",
			Input:          "abc-123-def",
			CustRepoOutput: &data.Customer{CoreBankId: 123},
			CustRepoError:  nil,
			CoreBankConnOutput: []comm.CoreBankInvoiceResp{
				{
					ProcessingSituation: "CLOSED",
					ClosingDate:         date(2023, 9, 10),
					ActualDueDate:       date(2023, 9, 15),
					TotalAmount:         888.88,
				},
				{
					ProcessingSituation: "CLOSED",
					ClosingDate:         date(2023, 10, 10),
					ActualDueDate:       date(2023, 10, 16),
					TotalAmount:         777.77,
				},
				{
					ProcessingSituation: "OPEN",
					ClosingDate:         date(2023, 11, 10),
					ActualDueDate:       date(2023, 11, 15),
					TotalAmount:         100,
				},
			},
			CoreBankConnError: nil,
			TranscRepoOutput:  []data.Transaction{{Amount: 13.4}, {Amount: 10}},
			TranscRepoError:   nil,
			ExpectedOutput: &busn.CurrInvoiceResp{
				StatusLabel: "Open",
				Amount:      "$ 123.40",
				ClosingDate: "NOV 10",
			},
			ExpectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(*testing.T) {
			custRepoMock := &mock_imp.CustRepoMock{}
			custRepoMock.On("FindById", mock.Anything).Return(test.CustRepoOutput, test.CustRepoError)

			coreBankConnMock := &mock_imp.CoreBankConnMock{}
			coreBankConnMock.On("GetAllInvoices", mock.Anything).Return(test.CoreBankConnOutput, test.CoreBankConnError)

			transcRepoMock := &mock_imp.TranscRepoMock{}
			transcRepoMock.On("FindAllByCustomerCoreBankId", mock.Anything).Return(test.TranscRepoOutput, test.TranscRepoError)

			invoServ := busn.NewInvoiceServ(custRepoMock, transcRepoMock, coreBankConnMock)

			resp, err := invoServ.GetCurrInvoice(test.Input)

			assert.Equal(t, test.ExpectedOutput, resp)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}

func date(y int, m time.Month, d int) string {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
}

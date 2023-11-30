package business_test

import (
	comm "dev-monteiro/go-digital-bank/commons"
	"dev-monteiro/go-digital-bank/commons/mnyamnt"
	mock_imp "dev-monteiro/go-digital-bank/credit-invoice/mock"
	busn "dev-monteiro/go-digital-bank/credit-invoice/src/business"
	conf "dev-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionServ_CreateFromPurchase(t *testing.T) {
	type testData struct {
		Name                    string
		Input                   *comm.PurchaseEvent
		TranscRepoExpectedInput *data.Transaction
		TranscRepoError         error
		ExpectedError           *conf.AppError
	}

	tests := []testData{
		{
			Name: "should_ReturnUnknownError_When_TransactionRepoReturnsError",
			Input: &comm.PurchaseEvent{
				Id:         123,
				CustomerId: 456,
				Amount:     mnyamnt.NewMnyAmount("19.99"),
			},
			TranscRepoExpectedInput: &data.Transaction{
				PurchaseId:         123,
				CustomerCoreBankId: 456,
				Amount:             mnyamnt.NewMnyAmount("19.99"),
			},
			TranscRepoError: errors.New("error"),
			ExpectedError: &conf.AppError{
				Message:    "Unknown error: error",
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			Name: "should_NotReturnError_When_TransactionRepoDoesNotReturnsError",
			Input: &comm.PurchaseEvent{
				Id:         123,
				CustomerId: 456,
				Amount:     mnyamnt.NewMnyAmount("19.99"),
			},
			TranscRepoExpectedInput: &data.Transaction{
				PurchaseId:         123,
				CustomerCoreBankId: 456,
				Amount:             mnyamnt.NewMnyAmount("19.99"),
			},
			TranscRepoError: nil,
			ExpectedError:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(*testing.T) {
			custRepoMock := &mock_imp.CustRepoMock{}

			transcRepoMock := &mock_imp.TranscRepoMock{}
			transcRepoMock.On("Save", test.TranscRepoExpectedInput).Return(test.TranscRepoError)

			transcServ := busn.NewTransactionServ(custRepoMock, transcRepoMock)

			err := transcServ.CreateFromPurchase(test.Input)

			assert.Equal(t, test.ExpectedError, err)
		})
	}
}

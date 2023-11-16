package mock_imp

import (
	comm "dev-monteiro/go-digital-bank/commons"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"

	"github.com/stretchr/testify/mock"
)

type CustRepoMock struct {
	mock.Mock
}

func (mock *CustRepoMock) FindById(id string) (*data.Customer, error) {
	args := mock.Called(id)
	return args.Get(0).(*data.Customer), args.Error(1)
}

func (mock *CustRepoMock) FindAllByCoreBankBatchId(coreBankBatchId int) ([]data.Customer, error) {
	args := mock.Called(coreBankBatchId)
	return args.Get(0).([]data.Customer), args.Error(1)
}

type TranscRepoMock struct {
	mock.Mock
}

func (mock *TranscRepoMock) Save(transc data.Transaction) error {
	args := mock.Called(transc)
	return args.Error(0)
}

func (mock *TranscRepoMock) FindAllByCustomerCoreBankId(custCoreBankId int) ([]data.Transaction, error) {
	args := mock.Called(custCoreBankId)
	return args.Get(0).([]data.Transaction), args.Error(1)
}

func (mock *TranscRepoMock) Delete(transc data.Transaction) error {
	args := mock.Called(transc)
	return args.Error(1)
}

type CoreBankConnMock struct {
	mock.Mock
}

func (mock *CoreBankConnMock) GetAllInvoices(custCoreBankId int) ([]comm.CoreBankInvoiceResp, error) {
	args := mock.Called(custCoreBankId)
	return args.Get(0).([]comm.CoreBankInvoiceResp), args.Error(1)
}

package business

import (
	comm "dev-monteiro/go-digital-bank/commons"
	conf "dev-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	"log"
)

type TransactionServ interface {
	CreateFromPurchase(purchase comm.PurchaseEvent) *conf.AppError
	ClearByBatch(batch comm.BatchEvent) *conf.AppError
}

type transactionServ struct {
	custRepo   data.CustomerRepo
	transcRepo data.TransactionRepo
}

func NewTransactionServ(custRepo data.CustomerRepo, transcRepo data.TransactionRepo) TransactionServ {
	return &transactionServ{
		custRepo:   custRepo,
		transcRepo: transcRepo,
	}
}

func (serv *transactionServ) CreateFromPurchase(purchase comm.PurchaseEvent) *conf.AppError {
	log.Println("[TransactionServ] CreateFromPurchase")

	transc := data.Transaction{
		PurchaseId:         purchase.Id,
		CustomerCoreBankId: purchase.CustomerId,
		Amount:             purchase.Amount,
	}

	err := serv.transcRepo.Save(transc)
	if err != nil {
		return conf.NewUnknownError(err)
	}

	return nil
}

func (serv *transactionServ) ClearByBatch(batch comm.BatchEvent) *conf.AppError {
	log.Println("[TransactionServ] ClearByBatch")

	custArr, err := serv.custRepo.FindAllByCoreBankBatchId(batch.Id)
	if err != nil {
		return conf.NewUnknownError(err)
	}

	for _, cust := range custArr {
		transcArr, err := serv.transcRepo.FindAllByCustomerCoreBankId(cust.CoreBankId)
		if err != nil {
			return conf.NewUnknownError(err)
		}

		for _, transc := range transcArr {
			err = serv.transcRepo.Delete(transc)
			if err != nil {
				return conf.NewUnknownError(err)
			}
		}
	}

	return nil
}

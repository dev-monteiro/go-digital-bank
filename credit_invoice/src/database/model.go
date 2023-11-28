package database

import (
	"dev-monteiro/go-digital-bank/commons/mnyamnt"
)

type Customer struct {
	Id              string
	CoreBankId      int
	CoreBankBatchId int
}

type Transaction struct {
	PurchaseId         int
	CustomerCoreBankId int
	Amount             *mnyamnt.MnyAmount
}

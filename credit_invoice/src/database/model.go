package database

import comm "dev-monteiro/go-digital-bank/commons"

type Customer struct {
	Id              string
	CoreBankId      int
	CoreBankBatchId int
}

type Transaction struct {
	PurchaseId         int
	CustomerCoreBankId int
	Amount             *comm.MoneyAmount
}

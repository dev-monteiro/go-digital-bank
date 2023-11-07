package database

type Customer struct {
	Id              string
	CoreBankId      int
	CoreBankBatchId int
}

type Transaction struct {
	PurchaseId         int
	CustomerCoreBankId int
	Amount             float64
}

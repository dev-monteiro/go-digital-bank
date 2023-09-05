package database

import "github.com/google/uuid"

type Customer struct {
	Id         string
	CoreBankId int
}

type Invoice struct {
	Id         string
	CoreBankId int
}

func NewInvoice(cbId int) *Invoice {
	return &Invoice{Id: uuid.New().String(), CoreBankId: cbId}
}

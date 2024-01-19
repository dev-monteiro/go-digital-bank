package business

import (
	comm "dev-monteiro/go-digital-bank/commons"
	conf "dev-monteiro/go-digital-bank/credit-invoice/src/configuration"
	data "dev-monteiro/go-digital-bank/credit-invoice/src/database"
	"log"
)

type CustomerServ interface {
	CreateFromEvent(event *comm.CustomerEvent) *conf.AppError
}

type customerServ struct {
	custRepo data.CustomerRepo
}

func NewCustomerServ(custRepo data.CustomerRepo) CustomerServ {
	return &customerServ{
		custRepo: custRepo,
	}
}

func (serv *customerServ) CreateFromEvent(event *comm.CustomerEvent) *conf.AppError {
	log.Println("[CustomerServ] CreateFromEvent")

	cust := &data.Customer{
		Id:              event.Id,
		CoreBankId:      event.CoreBankingCreditId,
		CoreBankBatchId: event.CoreBankingBatchId,
	}

	err := serv.custRepo.Save(cust)
	if err != nil {
		return conf.NewUnknownError(err)
	}

	return nil
}

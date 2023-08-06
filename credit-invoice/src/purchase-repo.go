package main

import (
	"errors"
	"fmt"
)

type PurchaseRepo struct {
	table map[int][]Purchase
}

func NewPurchaseRepo() PurchaseRepo {
	return PurchaseRepo{table: make(map[int][]Purchase)}
}

func (repo PurchaseRepo) save(purchase Purchase) {
	_, exists := repo.table[purchase.CreditAccountId]

	if exists {
		repo.table[purchase.CreditAccountId] = append(repo.table[purchase.CreditAccountId], purchase)
	} else {
		repo.table[purchase.CreditAccountId] = []Purchase{purchase}
	}
	fmt.Println("Saved on repo.")
}

func (repo PurchaseRepo) findAllByCreditAccountId(creditAccountId int) ([]Purchase, error) {
	purchases, exists := repo.table[creditAccountId]

	if exists {
		return purchases, nil
	} else {
		return nil, errors.New("not found")
	}
}

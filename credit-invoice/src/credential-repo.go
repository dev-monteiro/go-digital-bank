package main

import "errors"

type CredentialRepo struct {
	table map[string]int
}

func NewCredentialRepo() CredentialRepo {
	table := make(map[string]int)
	bootstrap(table)
	return CredentialRepo{table: table}
}

func (repo CredentialRepo) getCreditAccountId(customerId string) (int, error) {
	creditAccountId, exists := repo.table[customerId]

	if !exists {
		return 0, errors.New("not found")
	}

	return creditAccountId, nil
}

func bootstrap(table map[string]int) {
	table["abc-123-def"] = 123
}

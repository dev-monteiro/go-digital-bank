package main

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() Repository {
	db := NewDatabase()

	return Repository{
		db: db,
	}
}

func (repo Repository) Close() {
	repo.db.Close()
}

func (repo Repository) getCreditAccountId(customerId string) (int, error) {
	query := "select credit_account_id from customer_credentials where customer_id = '" + customerId + "';"
	row := repo.db.QueryRow(query)

	var creditAccountId int
	err := row.Scan(&creditAccountId)

	return creditAccountId, err
}

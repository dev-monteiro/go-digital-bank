package main

import (
	"database/sql"
	"fmt"
)

func NewDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(mysql)/credit_invoice")
	if err != nil {
		fmt.Println(err)
	}

	return db
}

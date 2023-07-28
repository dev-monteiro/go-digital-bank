package main

import (
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	time.Sleep(10 * time.Second)
}

func main() {
	listener := NewListener()
	defer listener.Close()

	controller := NewController()
	defer controller.Close()

	http.HandleFunc("/invoices/current", controller.getCurrentInvoice)

	http.ListenAndServe(":80", nil)
}

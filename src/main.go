package main

import (
	"io"
	"net/http"
)

func helloWorld(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Hello World!\n")
}

func main() {
	http.HandleFunc("/", helloWorld)

	http.ListenAndServe(":8096", nil)
}

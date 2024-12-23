package main

import (
	"fmt"
	"net/http"
)

func main() {

	mx := http.NewServeMux()
	mx.Handle("POST /edi-to-json", http.HandlerFunc(ediToJsonHandler))

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", mx)
}

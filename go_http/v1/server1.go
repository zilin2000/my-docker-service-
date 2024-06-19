package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello, this is v1\n")
}

func main() {
	port := 8080
	fmt.Printf("start listening on port %d...\n", port)
	http.HandleFunc("/v1", hello)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

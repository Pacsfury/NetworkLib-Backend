package main

import (
	"fmt"
	"net/http"
)


func main() {
	http.HandleFunc("/", listen)

	fmt.Println("Server running on :8080...")
	go http.ListenAndServe(":8080", nil)
	test()
}

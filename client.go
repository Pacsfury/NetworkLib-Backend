package main

/* 
Example client for tests.
This client tests all the implemented operations.
You can use this as an example client for creating you own clients in the future.
*/

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func test() {
	fmt.Println(" --- Test started  --- ")

	// 1. Test SET
	respSET, err := http.Post("http://localhost:8080", "text/plain", bytes.NewBuffer([]byte("SET a b")))
	if err != nil {
		fmt.Println("ERROR SET:", err)
		return
	}
	defer respSET.Body.Close()

	bodySET, _ := io.ReadAll(respSET.Body)
	fmt.Println("[Client got from server - SET]:")
	fmt.Println(string(bodySET))

	// 2. Test GET
	respGET, err := http.Post("http://localhost:8080", "text/plain", bytes.NewBuffer([]byte("GET a")))
	if err != nil {
		fmt.Println("ERROR GET:", err)
		return
	}
	defer respGET.Body.Close()

	bodyGET, _ := io.ReadAll(respGET.Body)
	fmt.Println("[Client got from server - GET]:")
	fmt.Println(string(bodyGET))
}

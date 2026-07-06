package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var vals = make(map[string]string)

func listen(w http.ResponseWriter, r *http.Request) {
	bytes, _ := io.ReadAll(r.Body)
	fmt.Println("Received: ", string(bytes))

	recstr := string(bytes)
	ops := strings.Split(recstr, " ")
	if ops[0] == "SET" {
		vals[string(ops[1])] = string(ops[2])
	} else if ops[0] == "GET" {
		sendBinary(w, []byte(vals[string(ops[1])]))
	}
	
}

func sendBinary(w http.ResponseWriter, data []byte) {
	w.Write(data)
	fmt.Println("Sent.")
}

func main() {
	http.HandleFunc("/", listen)

	fmt.Println("Server running on :8080...")
	http.ListenAndServe(":8080", nil)
}

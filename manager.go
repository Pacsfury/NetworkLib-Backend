package main

import (
	"net/http"
	"io"
	"fmt"
	"strings"
)

func listen(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bytes, _ := io.ReadAll(r.Body)
	fmt.Println("Received: ", string(bytes))

	recstr := strings.TrimSpace(string(bytes))

	ops := strings.Split(recstr, " ")

	command := ops[0]

	switch command {
	case "SET":
		if len(ops) < 3 {
			http.Error(w, "Missing value for SET", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		vals[ops[1]] = ops[2]
		mutex.Unlock()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")

	case "GET":
		mutex.RLock()
		val, exists := vals[ops[1]]
		mutex.RUnlock()

		if !exists {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		sendBinary(w, []byte(val))

	default:
		http.Error(w, "Unknown command", http.StatusBadRequest)
		fmt.Println(command)
	}

	
}

func sendBinary(w http.ResponseWriter, data []byte) {
	w.Write(data)
	fmt.Println("Sent.")
}
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	vals  = make(map[string]string)
	mutex sync.RWMutex
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

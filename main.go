package main

import (
	"fmt"
	"net"
)


func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Server error:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server running on :8080...")
	go test()
	go test2()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}
package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func test2() {
	fmt.Println(" --- [Client 2] Started --- ")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("ERROR [Client 2] connecting:", err)
		return
	}
	defer conn.Close()

	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadBytes(DEL)
			if err != nil {
				fmt.Println("\n[Client 2] Connection closed or error:", err)
				return
			}
			fmt.Println("[Client 2 received]:", string(msg[:len(msg)-1]))
		}
	}()

	time.Sleep(1 * time.Second)

	fmt.Println("[Client 2] Subscribing to 'pos_x'...")
	sendCmd(conn, SUB, "pos_x")

	time.Sleep(10 * time.Second)
}
package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// sendCmd writes an opcode followed by its DEL-terminated args.
func sendCmd(conn net.Conn, opcode byte, args ...string) {
	msg := []byte{opcode}
	for _, a := range args {
		msg = append(msg, []byte(a)...)
		msg = append(msg, DEL)
	}
	conn.Write(msg)
}

func test() {
	fmt.Println(" --- Test started  --- ")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("ERROR connecting:", err)
		return
	}
	defer conn.Close()

	go func() {
		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadBytes(DEL)
			if err != nil {
				fmt.Println("[Client] connection closed because this error:", err)
				return
			}
			fmt.Println("[Server]:", string(msg[:len(msg)-1]))
		}
	}()

	// 1. Test SET
	sendCmd(conn, SET, "a", "b")
	time.Sleep(50 * time.Millisecond)

	// 2. Test GET
	sendCmd(conn, GET, "a")
	time.Sleep(50 * time.Millisecond)

	// 3. Test TEMP
	sendCmd(conn, TEMP, "er", "b")
	time.Sleep(50 * time.Millisecond)

	// 4. Test CONST
	sendCmd(conn, CONST, "ar", "tr")
	time.Sleep(50 * time.Millisecond)

	// 5. Test edit CONST
	sendCmd(conn, SET, "ar", "trdf")
	time.Sleep(50 * time.Millisecond)

	// 6. Double GET at TEMP
	sendCmd(conn, GET, "er")
	time.Sleep(50 * time.Millisecond)

	sendCmd(conn, GET, "er")
	time.Sleep(50 * time.Millisecond)

	// 7. Test SIGNAL
	sendCmd(conn, SIGNAL, "#er")

	// 8. Test SUB
	sendCmd(conn, SET, "pos_x", "10")
	time.Sleep(1000 * time.Millisecond)

	sendCmd(conn, SET, "pos_x", "10")
	time.Sleep(1000 * time.Millisecond)

	time.Sleep(100 * time.Millisecond)
}
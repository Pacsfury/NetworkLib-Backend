package main

/* 
Example client for tests.
This client tests all the implemented operations.
You can use this as an example client for creating you own clients in the future.
*/

import (
	"bufio"
	"fmt"
	"net"
)

func test() {
	fmt.Println(" --- Test started  --- ")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("ERROR connecting:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// 1. Test SET
	fmt.Fprintln(conn, "SET a b")
	respSET, _ := reader.ReadString('\n')
	fmt.Println("[Client got from server - SET]:")
	fmt.Println(respSET)

	// 2. Test GET
	fmt.Fprintln(conn, "GET a")
	respGET, _ := reader.ReadString('\n')
	fmt.Println("[Client got from server - GET]:")
	fmt.Println(respGET)

	// 3. Test TEMP
	fmt.Fprintln(conn, "TEMP er b")
	respSET3, _ := reader.ReadString('\n')
	fmt.Println("[Client got from server - TEMP]:")
	fmt.Println(respSET3)

	// 4. Test CONST
	fmt.Fprintln(conn, "CONST ar tr")
	respGET4, _ := reader.ReadString('\n')
	fmt.Println("[Client got from server - CONST]:")
	fmt.Println(respGET4)

	// 5. Test edit CONST
	fmt.Fprintln(conn, "SET ar trdf")
	respGET5, _ := reader.ReadString('\n')
	fmt.Println("[Client got from server - SET]:")
	fmt.Println(respGET5)

	// 6. Double GET at TEMP
	fmt.Fprintln(conn, "GET er")
	respGET6, _ := reader.ReadString('\n')
	fmt.Println("[Client got from server - GET]:")
	fmt.Println(respGET6)

	fmt.Fprintln(conn, "GET er")
	respGET7, _ := reader.ReadString('\n')
	fmt.Println("[Client got from server - GET]:")
	fmt.Println(respGET7)

	// 7. Test SIGNAL
	fmt.Fprintln(conn, "SIGNAL #er")
	respGET8, _ := reader.ReadString('\n')
	fmt.Println("[Client got from server - SIGNAL]:")
	fmt.Println(respGET8)
}
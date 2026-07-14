package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"github.com/google/uuid"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		_, ok := connections[conn]
		if !ok {
			connections[conn] = connection {
				conn: conn,
				subscriptions: make(map[string]variable),
				uid: uuid.New().String(),
			}
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			return 
		}

		recstr := strings.TrimSpace(line)
		if recstr == "" {
			continue
		}

		fmt.Println("Received: ", recstr)

		ops := strings.Split(recstr, " ")
		command := ops[0]

		switch command {
		case "SET":
			if len(ops) < 3 {
				fmt.Fprintln(conn, "ERROR Missing value for SET")
				continue
			}

			mutex.Lock()
			if current, exists := vals[ops[1]]; exists && current.isconst {
				fmt.Fprintln(conn, "CONST var can't change")
				mutex.Unlock() 
				continue
			}

			vals[ops[1]] = variable{
				name:    ops[1],
				value:   ops[2],
				istemp:  false,
				isconst: false,
			}
			mutex.Unlock()
			fmt.Fprintln(conn, "OK")

		case "GET":
			if len(ops) < 2 {
				fmt.Fprintln(conn, "ERROR Missing key for GET")
				continue
			}

			mutex.RLock()
			val, exists := vals[ops[1]]
			mutex.RUnlock()

			if !exists {
				fmt.Fprintln(conn, "ERROR Key not found")
				continue
			}

			if val.istemp {
				mutex.Lock()
				if current, encorat := vals[ops[1]]; encorat && current.istemp {
					fmt.Println("Deleting TEMP variable...")
					delete(vals, ops[1])
				}
				mutex.Unlock()
			}

			fmt.Fprintln(conn, val.value)

		case "TEMP":
			if len(ops) < 3 {
				fmt.Fprintln(conn, "ERROR Missing value for SET")
				continue
			}

			mutex.Lock()
			vals[ops[1]] = variable{
				name:    ops[1],
				value:   ops[2],
				istemp:  true,
				isconst: false,
			}
			mutex.Unlock()
			fmt.Fprintln(conn, "OK")

		case "CONST":
			if len(ops) < 3 {
				fmt.Fprintln(conn, "ERROR Missing value for SET")
				continue
			}

			mutex.Lock()
			vals[ops[1]] = variable{
				name:    ops[1],
				value:   ops[2],
				istemp:  false,
				isconst: true,
			}
			mutex.Unlock()
			fmt.Fprintln(conn, "OK")
		
		case "SIGNAL":
			for connect, _ := range connections {
				connect.Write([]byte(ops[1]))
				connect.Close()
			}

		default:
			fmt.Fprintln(conn, "ERROR Unknown command")
			fmt.Println(command)
		}
	}
}
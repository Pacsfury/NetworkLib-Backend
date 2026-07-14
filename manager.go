package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"github.com/google/uuid"
)

func handleConn(conn net.Conn) {
	mutex.Lock()
	connections[conn] = connection{
		conn:          conn,
		subscriptions: make(map[string]variable),
		uid:           uuid.New().String(),
	}
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(connections, conn)
		mutex.Unlock()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	for {
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
			varName := ops[1]
			varValue := ops[2]

			mutex.Lock()
			if current, exists := vals[varName]; exists && current.isconst {
				fmt.Fprintln(conn, "CONST var can't change")
				mutex.Unlock() 
				continue
			}

			vals[varName] = variable{
				name:    varName,
				value:   varValue,
				istemp:  false,
				isconst: false,
			}

			for connect, clientCtx := range connections {
				if _, isSubscribed := clientCtx.subscriptions[varName]; isSubscribed {
					fmt.Fprintf(connect, "#subscribed_var_changed %s %s\n", varName, varValue)
				}
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
			if len(ops) < 2 {
				fmt.Fprintln(conn, "ERROR Missing value for SIGNAL")
				continue
			}

			mutex.Lock()
			for connect := range connections {
				fmt.Fprintln(connect, ops[1]) 
			}
			mutex.Unlock()

		case "SUB":
			if len(ops) < 2 {
				fmt.Fprintln(conn, "ERROR Missing variable name for SUB")
				continue
			}
			variableToSub := ops[1]

			mutex.Lock()
			if clientConn, exists := connections[conn]; exists {
				clientConn.subscriptions[variableToSub] = variable{
					name: variableToSub,
				}
				fmt.Fprintln(conn, "OK SUB")
			} else {
				fmt.Fprintln(conn, "ERROR Connection not registered")
			}
			mutex.Unlock()

		default:
			fmt.Fprintln(conn, "ERROR Unknown command")
			fmt.Println(command)
		}
	}
}

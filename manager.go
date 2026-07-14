package main

import (
	"bufio"
	"fmt"
	"net"
	"github.com/google/uuid"
)

// DEL is the field/record delimiter used to separate ascii-encoded
// arguments in both directions: [opcode][arg1][DEL][arg2][DEL]...
const DEL byte = 0x1F

// send writes a raw byte message terminated with DEL.
func send(conn net.Conn, msg string) {
	conn.Write(append([]byte(msg), DEL))
}

// sendf is a convenience wrapper for formatted messages.
func sendf(conn net.Conn, format string, a ...interface{}) {
	send(conn, fmt.Sprintf(format, a...))
}

// readArg reads a single DEL-terminated field and strips the delimiter.
func readArg(reader *bufio.Reader) (string, error) {
	data, err := reader.ReadBytes(DEL)
	if err != nil {
		return "", err
	}
	return string(data[:len(data)-1]), nil
}

// argCountFor returns how many DEL-delimited arguments follow a given opcode.
func argCountFor(opcode byte) int {
	switch opcode {
	case SET, TEMP, CONST:
		return 2
	case GET, SIGNAL, SUB:
		return 1
	default:
		return 0
	}
}

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
		opcode, err := reader.ReadByte()
		if err != nil {
			return
		}

		argCount := argCountFor(opcode)
		if argCount == 0 && opcode != SIGNAL {
			// unknown opcode with no defined arg count
			if _, known := map[byte]bool{SET: true, GET: true, TEMP: true, CONST: true, SIGNAL: true, SUB: true}[opcode]; !known {
				send(conn, "ERROR Unknown command")
				continue
			}
		}

		args := make([]string, 0, argCount)
		ok := true
		for i := 0; i < argCount; i++ {
			arg, err := readArg(reader)
			if err != nil {
				return
			}
			args = append(args, arg)
		}
		if !ok {
			continue
		}

		fmt.Printf("Received: opcode=0x%X args=%v\n", opcode, args)

		switch opcode {
		case SET:
			varName := args[0]
			varValue := args[1]

			mutex.Lock()
			if current, exists := vals[varName]; exists && current.isconst {
				send(conn, "CONST var can't change")
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
					sendf(connect, "#subscribed_var_changed %s %s", varName, varValue)
				}
			}
			mutex.Unlock()
			send(conn, "OK")

		case GET:
			key := args[0]

			mutex.RLock()
			val, exists := vals[key]
			mutex.RUnlock()

			if !exists {
				send(conn, "ERROR Key not found")
				continue
			}

			if val.istemp {
				mutex.Lock()
				if current, ok := vals[key]; ok && current.istemp {
					fmt.Println("Deleting TEMP variable...")
					delete(vals, key)
				}
				mutex.Unlock()
			}

			send(conn, val.value)

		case TEMP:
			varName := args[0]
			varValue := args[1]

			mutex.Lock()
			vals[varName] = variable{
				name:    varName,
				value:   varValue,
				istemp:  true,
				isconst: false,
			}
			mutex.Unlock()
			send(conn, "OK")

		case CONST:
			varName := args[0]
			varValue := args[1]

			mutex.Lock()
			vals[varName] = variable{
				name:    varName,
				value:   varValue,
				istemp:  false,
				isconst: true,
			}
			mutex.Unlock()
			send(conn, "OK")

		case SIGNAL:
			msg := args[0]

			mutex.Lock()
			for connect := range connections {
				send(connect, msg)
			}
			mutex.Unlock()

		case SUB:
			variableToSub := args[0]

			mutex.Lock()
			if clientConn, exists := connections[conn]; exists {
				clientConn.subscriptions[variableToSub] = variable{
					name: variableToSub,
				}
				send(conn, "OK SUB")
			} else {
				send(conn, "ERROR Connection not registered")
			}
			mutex.Unlock()

		default:
			send(conn, "ERROR Unknown command")
			fmt.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
)

// DEL separates individual ascii-encoded args: [opcode][arg1][DEL][arg2][DEL]...
const DEL byte = 0x1F

// END terminates a variable-length argument list (used by SIGNAL).
const END byte = 0x00

func send(conn net.Conn, msg string) {
	conn.Write(append([]byte(msg), DEL))
}

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

// readVarArgs reads DEL-terminated fields until it hits an END marker byte.
func readVarArgs(reader *bufio.Reader) ([]string, error) {
	args := make([]string, 0)
	for {
		peek, err := reader.Peek(1)
		if err != nil {
			return nil, err
		}
		if peek[0] == END {
			_, err := reader.ReadByte() // consume END
			if err != nil {
				return nil, err
			}
			return args, nil
		}
		arg, err := readArg(reader)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
}

// argCountFor returns (count, ok) for fixed-arity opcodes.
// SIGNAL is handled separately since it's variable-arity.
func argCountFor(opcode byte) (int, bool) {
	switch opcode {
	case SET, TEMP, CONST:
		return 2, true
	case GET, SUB:
		return 1, true
	default:
		return 0, false
	}
}

func handleConn(conn net.Conn) {
	mutex.Lock()
	connections[conn] = connection{
		conn:          conn,
		subscriptions: make(map[string]variable),
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

		var args []string

		if opcode == SIGNAL {
			args, err = readVarArgs(reader)
			if err != nil {
				return
			}
		} else {
			count, known := argCountFor(opcode)
			if !known {
				send(conn, "ERROR Unknown command")
				fmt.Printf("Unknown opcode: 0x%X\n", opcode)
				continue
			}
			args = make([]string, 0, count)
			for i := 0; i < count; i++ {
				arg, err := readArg(reader)
				if err != nil {
					return
				}
				args = append(args, arg)
			}
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
			msg := ""
			for i, a := range args {
				if i > 0 {
					msg += " "
				}
				msg += a
			}

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
		}
	}
}
package main

import (
	"sync" 
	"net"
)

type variable struct {
	name    string
	value   string
	istemp  bool
	isconst bool
}

type connection struct {
	conn            net.Conn
	subscriptions   map[string]variable
	uid             string
}

var (
	vals  = make(map[string]variable)
	mutex sync.RWMutex

	connections = make(map[net.Conn]connection)
)

const (
	SET = 0x1
	GET = 0x2
	TEMP = 0x3
	CONST = 0x4
	SIGNAL = 0x5
	SUB = 0x6
)
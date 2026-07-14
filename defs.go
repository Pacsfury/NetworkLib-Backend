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

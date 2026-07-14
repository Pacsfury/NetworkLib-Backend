package main

import "sync"

type variable struct {
	name    string
	value   string
	istemp  bool
	isconst bool
}

var (
	vals  = make(map[string]variable)
	mutex sync.RWMutex
)

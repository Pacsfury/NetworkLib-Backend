package main

import "sync"

var (
	vals  = make(map[string]string)
	mutex sync.RWMutex
)

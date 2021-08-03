package net

import (
	"fmt"
)

const addr = "127.0.0.1"
const port = 12345

func serverAddr() string {
	return fmt.Sprintf("%s:%d", addr, port)
}
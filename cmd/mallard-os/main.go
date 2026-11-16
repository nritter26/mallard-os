package main

import (
	"time"
)

func main() {
	for {
		println("Hello from TinyGo!")
		time.Sleep(1 * time.Second)
	}
}

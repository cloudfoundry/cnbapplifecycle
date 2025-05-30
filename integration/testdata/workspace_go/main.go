package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("hello world")
		time.Sleep(10 * time.Second)
	}

}

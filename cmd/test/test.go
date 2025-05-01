package main

import (
	"fmt"
	"time"
)

func AddFunc(taskName string, cmd func()) (int, error) {
	id := 101
	cmd()
	return id, nil
}

func main() {
	id := 100

	id, _ = AddFunc("test", func() {
		fmt.Println("test:", id)
	})

	fmt.Println("----->", id)

	// Wait to observe task execution
	time.Sleep(3 * time.Second)
}

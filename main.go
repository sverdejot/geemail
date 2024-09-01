package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	msgsSrv := NewMessageService(b)
	senders, err := msgsSrv.GetCountBySender()
	fmt.Println(senders)
}

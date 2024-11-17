package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	geemail "github.com/sverdejot/geemail/pkg"
	"github.com/sverdejot/geemail/pkg/auth"
)

func main() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	client, err := auth.NewHTTPClient(b)
	if err != nil {
		log.Fatalf("Unable to create default HTTP client: %v", err)
	}

	service := geemail.NewMessageService(client)

	t0 := time.Now()
	contents, err := service.GetContent(context.TODO())
	t1 := time.Now().Sub(t0)

	if err != nil {
		log.Fatalf("Unable to retrieve latest messages: %v", err)
	}

	for _, c := range contents {
		fmt.Println(c)
	}

	fmt.Println("took: ", t1)
}

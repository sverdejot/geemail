package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sverdejot/geemail/internal"
	"github.com/sverdejot/geemail/internal/auth"
)

func main() {
    credentialsPath := "config/credentials.json"
    if envCredentialsPath := os.Getenv("CREDENTIALS_FILE"); envCredentialsPath != "" {
        credentialsPath = envCredentialsPath
    }
    b, err := os.ReadFile(credentialsPath)
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

    client, err := auth.NewHTTPClient(b)
    if err != nil {
        log.Fatalf("Unable to create default HTTP client: %v", err)
    }

    service, err := geemail.NewMessageService(client)
    if err != nil {
        log.Fatalf("Unable to create message service: %v", err)
    }

    t0 := time.Now()
    contents, err := service.GetContent(context.TODO())
    t1 := time.Now().Sub(t0)

    if err != nil {
        log.Fatalf("Unable to retrieve latest messages: %v", err)
    }

    fmt.Println(geemail.FormatCount(contents))

    fmt.Println("took: ", t1)
}

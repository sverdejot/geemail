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

    for _, c := range contents {
        fmt.Println(c)
    }

    fmt.Println("took: ", t1)
}

package main

import (
	"log"
	"os"

	"github.com/sverdejot/geemail/pkg/auth"
)

func main() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	if err := auth.GenerateToken(b); err != nil {
		log.Fatalf("Cannot generate token: %v", err)
	}
}

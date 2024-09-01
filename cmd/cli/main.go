package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sverdejot/geemail/internal"
	"github.com/sverdejot/geemail/internal/auth"
	"gopkg.in/yaml.v3"
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

	service := internal.NewMessageService(client)
	senders, err := service.GetCountBySender()
	if err != nil {
		log.Fatalf("Unable to retrieve latest messages: %v", err)
	}

	requested := make([]string, 0, len(senders))
	for _, s := range senders {
		if requestAdding(s.Key) {
			requested = append(requested, s.Key)
		}
	}

	f := internal.NewRuleFile()
	for _, r := range requested {
		f = f.AddRule(r, true)
	}
	fb, _ := os.Create("rules.yaml")
	yaml.NewEncoder(fb).Encode(f)
}

func requestAdding(sender string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("do you want to add %s to removal set [y/N]: ", sender)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" || response == "" {
			return false
		}
	}
}

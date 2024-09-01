package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	user = "me"
)

var (
	scopes []string = []string{
		gmail.MailGoogleComScope,
		gmail.GmailAddonsCurrentMessageReadonlyScope,
		gmail.GmailAddonsCurrentMessageMetadataScope,
	}
)

type MessageService struct {
	srv *gmail.Service
}

func NewMessageService(credentials []byte) *MessageService {
	config, err := google.ConfigFromJSON(credentials, scopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return &MessageService{
		srv: srv,
	}
}

func (s *MessageService) GetCountBySender() (map[string]int, error) {
	req := s.srv.Users.Messages.
		List(user).
		MaxResults(500)

	msgs, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages for user %s: %w", user, err)
	}

	senders := make(map[string]int)
	for _, rawMsg := range msgs.Messages {
		msg, err := s.srv.Users.Messages.Get(user, rawMsg.Id).Do()
		if err != nil {
			log.Printf("failed to fetch message id [%s]: %v\n", rawMsg.Id, err)
			continue
		}
		sender := getSenderFromHeader(msg)
		domain := sender[strings.Index(sender, "@")+1:]

		senders[domain] += 1
	}
	return senders, nil
}

func getSenderFromHeader(message *gmail.Message) string {
	for _, h := range message.Payload.Headers {
		if h.Name == "From" {
			return strings.Trim(h.Value, "<>")
		}
	}
	return ""
}

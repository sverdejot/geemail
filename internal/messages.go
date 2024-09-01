package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	user = "me"
)

type MessageService struct {
	srv *gmail.Service
}

func NewMessageService(client *http.Client) *MessageService {
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return &MessageService{
		srv: srv,
	}
}

func (s *MessageService) GetCountBySender() ([]KVPair, error) {
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
	return sort(senders), nil
}

func getSenderFromHeader(message *gmail.Message) string {
	for _, h := range message.Payload.Headers {
		if h.Name == "From" {
			return strings.Trim(h.Value, "<>")
		}
	}
	return ""
}

type KVPair struct {
	Key   string
	Value int
}

func sort(senders map[string]int) []KVPair {
	sl := make([]KVPair, 0, len(senders))

	for k, v := range senders {
		sl = append(sl, KVPair{k, v})
	}

	slices.SortFunc(sl, func(a, b KVPair) int {
		switch {
		case a.Value > b.Value:
			return -1
		case a.Value < b.Value:
			return 1
		default:
			return 0
		}
	})

	return sl
}

func (kv KVPair) String() string {
	return fmt.Sprintf("%s: %d", kv.Key, kv.Value)
	}

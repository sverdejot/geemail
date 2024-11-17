package geemail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	user       = "me"
	query      = "is:unread has:nouserlabels"
	maxResults = 20
)

type MessageService struct {
	srv *gmail.Service
}

func NewMessageService(client *http.Client) (*MessageService, error) {
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create gmail service: %w", err)
	}
	return &MessageService{
		srv: srv,
	}, nil
}

func (s *MessageService) GetContent(ctx context.Context) ([]Content, error) {
	req := s.srv.Users.Messages.
		List(user).
		Q(query).
		MaxResults(maxResults).
		Context(ctx)

	res, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages for user %s: %w", user, err)
	}

	var wg sync.WaitGroup
	wg.Add(len(res.Messages))

	contents := make([]Content, len(res.Messages))
	for i, rawMsg := range res.Messages {
		go func() {
			defer wg.Done()
			msg, err := s.srv.Users.Messages.Get(user, rawMsg.Id).Do()
			if err != nil || msg.Id == "" {
				log.Printf("failed to fetch message id [%s]: %v\n", rawMsg.Id, err)
				return
			}
			contents[i] = Content{msg.Id, getSubject(msg), getSnippet(msg)}
		}()
	}
	wg.Wait()

	return contents, nil
}

func getSubject(message *gmail.Message) string {
	if message == nil || message.Payload == nil {
		return ""
	}

	for _, h := range message.Payload.Headers {
		if h.Name == "Subject" {
			return h.Value
		}
	}

	return ""
}

func getSnippet(msg *gmail.Message) string {
	if msg == nil {
		return ""
	}

	return msg.Snippet
}

type Content struct {
	ID, Subject, Snippet string
}

func (c Content) String() string {
	return fmt.Sprintf("%s: the subject of this message is \"%s\" and the snippet is \"%s\"", c.ID, c.Subject, c.Snippet)
}


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
    maxResults = 100
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
            contents[i] = Content{msg.Id, getSubject(msg), getSnippet(msg), getSender(msg)}
        }()
    }
    wg.Wait()

    return contents, nil
}

func getSender(message *gmail.Message) string {
    if message == nil || message.Payload == nil {
        return ""
    }

    for _, h := range message.Payload.Headers {
        if h.Name == "From" {
            return h.Value
        }
    }

    return ""
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
    ID, Subject, Snippet, From string
}

func (c Content) String() string {
    return fmt.Sprintf("%s [%s]: the subject of this message is \"%s\" and the snippet is \"%s\"", c.From, c.ID, c.Subject, c.Snippet)
}

func countBySender(msgs []Content) map[string]int {
    counts := make(map[string]int)

    for _, msg := range msgs {
        if _, ok := counts[msg.From]; ok {
            counts[msg.From] += 1
            continue
        }
        counts[msg.From] = 1
    }

    return counts
}

func FormatCount(c []Content) string {
    cm := countBySender(c)
    var s string
    for k, v := range cm {
        s += fmt.Sprintf("%s: %d\n", k, v)
    }
    return s
}

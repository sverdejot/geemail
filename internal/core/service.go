package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	"golang.org/x/time/rate"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	poolSize = runtime.NumCPU()
)

const (
	user = "me"

	// non-classified emails nor reads ones, reasoning is that read messages
	// may be interesting for the user
	query      = "is:unread has:nouserlabels"
	inboxLabel = "INBOX"

	// max query results, set to the API maximum, default is 100
	maxResults = 500

	// usage limit is 15_000 u/min per user, let's leave some room
	apiQuotaLimitPerMinute = 14_900

	// quota consumption per operation
	messagesListQuotaUsage = 5
	messagesGetQuotaUsage  = 5
)

type MailService struct {
	srv *gmail.Service
	lim *rate.Limiter
}

func NewMessageService(ctx context.Context, client *http.Client) (*MailService, error) {
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create gmail service: %w", err)
	}

	lim := rate.NewLimiter(
		rate.Limit(250), 250,
	)

	return &MailService{
		srv: srv,
		lim: lim,
	}, nil
}

func (s *MailService) StreamUnreadMessages(ctx context.Context) (chan RawMail, error) {
	if err := s.lim.WaitN(ctx, messagesListQuotaUsage); err != nil {
		return nil, err
	}

    ids, err := s.GetUnreadMessageIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages for user %s: %w", user, err)
	}

	var wg sync.WaitGroup
	wg.Add(poolSize)
	jobs := make(chan string, poolSize)
	results := make(chan RawMail, poolSize)

	for range poolSize {
		go s.getMessageWorker(ctx, &wg, jobs, results)
	}

	go func() {
		for _, id := range ids {
			jobs <- id
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return results, nil
}

func (s *MailService) GetUnreadMessages(ctx context.Context) ([]RawMail, error) {
	results, err := s.StreamUnreadMessages(ctx)
	if err != nil {
		return nil, err
	}

	contents := make([]RawMail, 0)
	for msg := range results {
		contents = append(contents, msg)
	}

	return contents, nil
}

func (s *MailService) getMessageWorker(
	ctx context.Context,
	wg *sync.WaitGroup,
	jobs <-chan string,
	results chan<- RawMail,
) {
	defer wg.Done()

	const maxRetries = 3

	for msgID := range jobs {
		var lastErr error

		for attempt := 0; attempt < maxRetries; attempt++ {
			if err := s.lim.WaitN(ctx, messagesGetQuotaUsage); err != nil {
				lastErr = err
				continue
			}

			msg, err := s.srv.Users.Messages.
				Get(user, msgID).
				Context(ctx).
				Do()

			if err == nil && msg.Id != "" {
				mail, err := NewRawMail(
					WithID(msg),
					WithSender(msg),
					WithSnippet(msg),
					WithSubject(msg),
					WithHeaders(msg),
				)
				if err != nil {
					// something about the mail cannot be parsed, continue
					// without retrying
					continue
				}
				results <- mail
				lastErr = nil
				break
			}

			lastErr = err
		}

		if lastErr != nil {
			log.Printf("message %s failed after retries: %v", msgID, lastErr)
		}
	}
}

func (s *MailService) GetTotalUnreads(ctx context.Context) (int64, error) {
	req := s.srv.Users.Labels.
		Get(user, inboxLabel).
		Context(ctx)

	label, err := req.Do()
	if err != nil {
		return 0, fmt.Errorf("error fetching labels: %w", err)
	}
	return label.MessagesUnread, nil
}

func (s *MailService) GetUnreadMessageIDs(ctx context.Context) ([]string, error) {
	var currentMailCount int64
	var pageToken string

	totalMailCount, err := s.GetTotalUnreads(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch total mail count: %w", err)
	}

	req := s.srv.Users.Messages.
		List(user).
		Q(query).
		MaxResults(maxResults).
		Context(ctx)

	mailIDs := make([]string, 0, totalMailCount)
	for currentMailCount < totalMailCount {
		if err := s.lim.WaitN(ctx, messagesListQuotaUsage); err != nil {
			// probably a ctx.Canceled, so better to return i.o. continue
			return nil, fmt.Errorf("cannot fetch whole mailing list: %w", err)
		}
		resp, err := req.
			PageToken(pageToken).
			Do()
		if err != nil {
			continue
		}
		pageToken = resp.NextPageToken
		currentMailCount += int64(len(resp.Messages))

		mailIDs = append(mailIDs, getIds(resp.Messages)...)
	}

	return mailIDs, nil
}

func getIds(msgs []*gmail.Message) []string {
	ids := make([]string, 0, len(msgs))
	for _, m := range msgs {
		ids = append(ids, m.Id)
	}
	return ids
}

package inbox

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var ErrNoUnsubscriber = errors.New("no unsubscriber")

const (
	unsubscribeHeaderKey     = "list-unsubscribe"
	postUnsubscribeHeaderKey = "list-unsubscribe-post"
)

type unsubscriber struct {
	target *url.URL
}

func NewOneClickUnsubscriber(headers map[string][]string) *unsubscriber {
	postHeader, ok := headers[postUnsubscribeHeaderKey]
	if !ok || len(postHeader) == 0 {
		return nil
	}

	if strings.TrimSpace(strings.ToLower(postHeader[0])) != "list-unsubscribe=one-click" {
		return nil
	}

	listVals, ok := headers[unsubscribeHeaderKey]
	if !ok || len(listVals) == 0 {
		return nil
	}

	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindAllStringSubmatch(listVals[0], -1)
	if len(matches) == 0 {
		return nil
	}

	var target *url.URL
	for _, m := range matches {
		u, err := url.Parse(strings.TrimSpace(m[1]))
		if err != nil {
			continue
		}

		if u.Scheme == "https" {
			target = u
			break
		}
	}

	if target == nil {
		return nil
	}

	return &unsubscriber{target: target}
}

func (u *unsubscriber) Do(ctx context.Context) error {
	if u == nil {
		return ErrNoUnsubscriber
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u.target.String(),
		strings.NewReader("List-Unsubscribe=One-Click"),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "GheeMail/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return fmt.Errorf("failed req, status %s", resp.Status)
	}
	return nil
}

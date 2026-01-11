package inbox

import (
	"fmt"
	"net/mail"
	"strings"

	"google.golang.org/api/gmail/v1"
)

type RawMailList []RawMail

type RawMail struct {
	ID               string
	From             string
	To               []string
	Subject, Snippet string
	Headers          map[string][]string
}

type RawMailOpt func(*RawMail) error

func NewRawMail(opts ...RawMailOpt) (r RawMail, err error) {
	for _, fn := range opts {
		err = fn(&r)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func WithID(msg *gmail.Message) RawMailOpt {
	return func(rm *RawMail) error {
		rm.ID = msg.Id
		return nil
	}
}

func WithHeaders(msg *gmail.Message) (fn RawMailOpt) {
	hs := make(map[string][]string)
	for _, h := range msg.Payload.Headers {
		name := strings.ToLower(h.Name)
		hs[name] = append(hs[name], h.Value)
	}
	return func(rm *RawMail) error {
		rm.Headers = hs
		return nil
	}
}

func WithSubject(msg *gmail.Message) (fn RawMailOpt) {
	if msg == nil || msg.Payload == nil {
		return func(rm *RawMail) error {
			return fmt.Errorf("malformed mail: empty or nil message: %v", msg)
		}
	}

	for _, h := range msg.Payload.Headers {
		if h.Name == "Subject" {
			return func(c *RawMail) error {
				c.Subject = h.Value
				return nil
			}
		}
	}

	return func(rm *RawMail) error {
		return fmt.Errorf("malformed mail: no subject")
	}
}

func WithSnippet(msg *gmail.Message) (fn RawMailOpt) {
	if msg == nil || msg.Payload == nil {
		return func(rm *RawMail) error {
			return fmt.Errorf("malformed mail: emtpy message")
		}
	}
	return func(rm *RawMail) error {
		rm.Snippet = msg.Snippet
		return nil
	}
}

func WithSender(msg *gmail.Message) RawMailOpt {
	if msg == nil || msg.Payload == nil {
		return func(rm *RawMail) error {
			return fmt.Errorf("malformed mail: emtpy message")
		}
	}

	for _, h := range msg.Payload.Headers {
		if h.Name == "From" {
			addr, err := mail.ParseAddress(h.Value)
			if err != nil {
				return func(rm *RawMail) error {
					return fmt.Errorf("malformed mail: unparseable sender")
				}
			}
			return func(rm *RawMail) error {
				rm.From = addr.Address
				return nil
			}
		}
	}

	return func(rm *RawMail) error {
		return fmt.Errorf("malformed mail: no sender found")
	}
}

func (rm RawMail) FilterValue() string {
	return rm.From
}

func (rm RawMail) Title() string {
	return rm.From
}

func (rm RawMail) Description() string {
	return rm.Snippet
}

func (rm RawMail) String() string {
	return fmt.Sprintf("%s [%s]: the subject of this message is \"%s\" and the snippet is \"%s\"", rm.From, rm.ID, rm.Subject, rm.Snippet)
}

func (rml RawMailList) GroupBySender() map[string][]RawMail {
	senders := make(map[string][]RawMail)

	for _, msg := range rml {
		if _, ok := senders[msg.From]; ok {
			senders[msg.From] = append(senders[msg.From], msg)
			continue
		}
		senders[msg.From] = make([]RawMail, 0)
		senders[msg.From] = append(senders[msg.From], msg)
	}

	return senders
}

func (rml RawMailList) String() string {
	var s string
	for k, v := range rml.GroupBySender() {
		s += fmt.Sprintf("%s: %d\n", k, len(v))
	}
	return s
}

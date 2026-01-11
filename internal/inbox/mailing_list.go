package inbox

import (
	"context"
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
)

var (
	unsubscribeListElemstyle = lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("red"))
)

type MailingList struct {
	From              string
	TotalUnreads      int
	UnreadMessagesIDs []string
	Unsubscriber      *unsubscriber
}

func GetMailingList(l RawMailList) []MailingList {
	lists := make([]MailingList, 0)
	mails := l.GroupBySender()

	var total int
	for s, sm := range mails {
		ids := make([]string, 0, len(sm))
		var unsubscriber *unsubscriber
		for _, rm := range sm {
			if isMailingList(rm) {
				total += 1
				ids = append(ids, rm.ID)
				if unsubscriber == nil {
					unsubscriber = NewOneClickUnsubscriber(rm.Headers)
				}
			}
		}
		if total > 0 {
			lists = append(lists, MailingList{
				From:              s,
				TotalUnreads:      total,
				UnreadMessagesIDs: ids,
				Unsubscriber:      unsubscriber,
			})
		}
		total = 0
	}
	sortAscendingByTotalUnreads(lists)
	return lists
}

func (m MailingList) Unsubscribe(ctx context.Context) error {
	return m.Unsubscriber.Do(ctx)
}

func sortAscendingByTotalUnreads(l []MailingList) {
	sort.Slice(l, func(i, j int) bool {
		return l[i].TotalUnreads > l[j].TotalUnreads
	})
}

func isMailingList(rm RawMail) bool {
	_, ok := rm.Headers[unsubscribeHeaderKey]
	return ok
}

func (rm MailingList) FilterValue() string {
	return rm.From
}

func (rm MailingList) Title() string {
	if rm.Unsubscriber != nil {
		return unsubscribeListElemstyle.Render(rm.From)
	}
	return rm.From
}

func (rm MailingList) Description() string {
	return fmt.Sprintf("%d unread", rm.TotalUnreads)
}

func (rm MailingList) UnsubscribeAvailable() bool {
	return rm.Unsubscriber != nil
}

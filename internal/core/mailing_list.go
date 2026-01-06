package core

import "fmt"

type MailingList struct {
	From         string
	TotalUnreads int
}

func GetMailingList(l RawMailList) []MailingList {
	lists := make([]MailingList, 0)
	mails := l.GroupBySender()

	var total int
	for s, sm := range mails {
		for _, rm := range sm {
			if isMailingList(rm) {
				total += 1
			}
		}
		if total > 0 {
			lists = append(lists, MailingList{
				From:         s,
				TotalUnreads: total,
			})
		}
		total = 0
	}
	return lists
}

func isMailingList(rm RawMail) bool {
	_, ok := rm.Headers["list-unsubscribe"]
	return ok
}

func (rm MailingList) FilterValue() string {
	return rm.From
}

func (rm MailingList) Title() string {
	return rm.From
}

func (rm MailingList) Description() string {
	return fmt.Sprintf("%d unread", rm.TotalUnreads)
}

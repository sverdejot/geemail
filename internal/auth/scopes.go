package auth

import "google.golang.org/api/gmail/v1"

var (
    scopes []string = []string{
        gmail.MailGoogleComScope,
        gmail.GmailAddonsCurrentMessageReadonlyScope,
        gmail.GmailAddonsCurrentMessageMetadataScope,
    }
)

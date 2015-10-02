package model

import "github.com/jordan-wright/email"

// Message represents an email message parsed.
type Message struct {
	ID       string         `json:"id"`
	ThreadID string         `json:"thread_id,omitempty"`
	Tags     []string       `json:"tags,omitempty"`
	Emails   []*email.Email `json:"emails,omitempty"`
}

package model

// Thread represents a thread.
type Thread struct {
	ID      string   `json:"id,omitempty"`
	Subject string   `json:"subject,omitempty"`
	Authors []string `json:"authors,omitempty"`
}

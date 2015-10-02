package server

import "github.com/gmuch/gmuch/model"

// QueryRequest represents a query request.
type QueryRequest struct {
	Q      string `json:"q"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

// QueryResponse represents a query response.
type QueryResponse struct {
	Threads []*model.Thread `json:"threads,omitempty"`
}

// ThreadRequest represents a thread request.
type ThreadRequest struct {
	ID string `json:"id"`
}

// ThreadResponse represents a thread response.
type ThreadResponse struct {
	model.Thread
}

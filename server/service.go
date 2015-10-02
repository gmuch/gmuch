package server

// GmuchService is an interface representing the Gmuch API service.
type GmuchService interface {
	Query(qs string, offset, limit int) (*QueryResponse, error)
	Thread(id string) (*ThreadResponse, error)
}

// ServiceMiddleware defines Gmuch API middleware.
type ServiceMiddleware func(GmuchService) GmuchService

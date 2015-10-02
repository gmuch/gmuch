package server

import "github.com/gmuch/gmuch/model"

// GmuchService is an interface representing the Gmuch API service.
type GmuchService interface {
	Query(qs string, offset, limit int) ([]*model.Thread, error)
	Thread(id string) (*model.Thread, error)
}

// ServiceMiddleware defines Gmuch API middleware.
type ServiceMiddleware func(GmuchService) GmuchService

package http

import (
	"github.com/gmuch/gmuch/server"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// EndpointenizeQuery transforms Query to an Endpoint.
func EndpointenizeQuery(gmuch server.GmuchService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(QueryRequest)
		ts, err := gmuch.Query(req.Query, req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}
		return QueryResponse{ts}, nil
	}
}

// EndpointenizeThread transforms Thread to an Endpoint.
func EndpointenizeThread(gmuch server.GmuchService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ThreadRequest)
		t, err := gmuch.Thread(req.ID)
		if err != nil {
			return nil, err
		}
		return ThreadResponse{t}, nil
	}
}

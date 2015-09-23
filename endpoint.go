package gmuch

import (
	"github.com/gmuch/gmuch/server"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// EndpointenizeQuery transforms Query to an Endpoint.
func EndpointenizeQuery(gmuch server.GmuchService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(server.QueryRequest)
		return gmuch.Query(req.Q, req.Offset, req.Limit)
	}
}

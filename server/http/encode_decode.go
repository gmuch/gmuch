package http

import (
	"encoding/json"
	"net/http"

	"github.com/gmuch/gmuch/server"
)

const defaultQueryLimit = 100

// DecodeQueryRequest decodes a query request from JSON and validates the
// offset and limit.
func DecodeQueryRequest(r *http.Request) (interface{}, error) {
	var request server.QueryRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	if request.Limit <= request.Offset {
		request.Limit = request.Offset + defaultQueryLimit
	}
	return request, err
}

// EncodeQueryResponse encodes a request into JSON and writes it the http ResponseWriter.
func EncodeQueryResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

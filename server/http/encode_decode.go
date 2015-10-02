package http

import (
	"encoding/json"
	"net/http"
)

const defaultQueryLimit = 100

// DecodeQueryRequest decodes a query request from JSON and validates the
// offset and limit.
func DecodeQueryRequest(r *http.Request) (interface{}, error) {
	var request QueryRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	if request.Limit <= request.Offset {
		request.Limit = request.Offset + defaultQueryLimit
	}
	return request, err
}

// EncodeQueryResponse encodes a request into JSON and writes it to the http
// ResponseWriter.
func EncodeQueryResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// DecodeThreadRequest decodes a thread request from JSON.
func DecodeThreadRequest(r *http.Request) (interface{}, error) {
	var request ThreadRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}

// EncodeThreadResponse encodes a request into JSON and writes it to the http
// ResponseWriter.
func EncodeThreadResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

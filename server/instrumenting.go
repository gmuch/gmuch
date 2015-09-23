package server

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.TimeHistogram
	GmuchService
}

func InstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.TimeHistogram) ServiceMiddleware {
	return func(next GmuchService) GmuchService {
		return instrumentingMiddleware{requestCount, requestLatency, next}
	}
}

func (im instrumentingMiddleware) Query(qs string, offset, limit int) (*QueryResponse, error) {
	var (
		qr  *QueryResponse
		err error
	)

	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "query"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		im.requestCount.With(methodField).With(errorField).Add(1)
		im.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	qr, err = im.GmuchService.Query(qs, offset, limit)
	return qr, err
}

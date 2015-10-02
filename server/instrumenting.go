package server

import (
	"fmt"
	"time"

	"github.com/gmuch/gmuch/model"
	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.TimeHistogram
	GmuchService
}

// InstrumentingMiddleware provides instrumentation for the service.
func InstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.TimeHistogram) ServiceMiddleware {
	return func(next GmuchService) GmuchService {
		return instrumentingMiddleware{requestCount, requestLatency, next}
	}
}

// Query implements GmuchService.Query with instrumentation decorator.
func (im instrumentingMiddleware) Query(qs string, offset, limit int) ([]*model.Thread, error) {
	var (
		qr  []*model.Thread
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

// Thread implements GmuchService.Thread with instrumentation decorator.
func (im instrumentingMiddleware) Thread(id string) (*model.Thread, error) {
	var (
		tr  *model.Thread
		err error
	)

	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "thread"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		im.requestCount.With(methodField).With(errorField).Add(1)
		im.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	tr, err = im.GmuchService.Thread(id)
	return tr, err
}

package server

import (
	"time"

	stdlog "log"

	"github.com/gmuch/gmuch/model"
	"github.com/go-kit/kit/log"
)

type logMiddleware struct {
	logger log.Logger
	GmuchService
}

// LoggingMiddleware represents a logging middleware that wraps the service.
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next GmuchService) GmuchService {
		return logMiddleware{logger, next}
	}
}

// Query wraps GmuchService.Query with logging logic.
func (lm logMiddleware) Query(qs string, offset, limit int) ([]*model.Thread, error) {
	defer func(begin time.Time) {
		if err := lm.logger.Log(
			"method", "query",
			"query", qs,
			"offset", offset,
			"limit", limit,
			"took", time.Since(begin),
		); err != nil {
			stdlog.Printf("error logging: %s", err)
		}
	}(time.Now())

	return lm.GmuchService.Query(qs, offset, limit)
}

// Thread wraps GmuchService.Thread with logging logic.
func (lm logMiddleware) Thread(id string) (*model.Thread, error) {
	defer func(begin time.Time) {
		if err := lm.logger.Log(
			"method", "thread",
			"id", id,
			"took", time.Since(begin),
		); err != nil {
			stdlog.Printf("error logging: %s", err)
		}
	}(time.Now())

	return lm.GmuchService.Thread(id)
}

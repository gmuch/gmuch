package gmuch

import "github.com/go-kit/kit/log"

// Gmuch represents the Gmuch service providing the business logic.
type Gmuch struct {
	dbPath string
	logger log.Logger
}

// New returns a new Gmuch.
func New(dbPath string, logger log.Logger) *Gmuch {
	return &Gmuch{dbPath, logger}
}

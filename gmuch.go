package gmuch

// Gmuch represents the Gmuch service providing the business logic.
type Gmuch struct {
	dbPath string
}

// New returns a new Gmuch.
func New(dbPath string) *Gmuch {
	return &Gmuch{dbPath}
}

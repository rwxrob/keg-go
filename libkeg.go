package libkeg

// Index interface assumes any source of relational data store
type Index interface {
}

type index struct {
	Path string
}

// CreateIndex creates a Index implementation that is backed by
// a sqlite3 database file at location specified by path. If it does not
// exist and returns a reference to it. Otherwise, just opens and
// returns a new reference to it.
func CreateIndex(path string) (*index, error) {

	index := new(index)

	// detect existing db and just open if so

	// create a new sqlite3 db

	return index, nil
}

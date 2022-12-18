package libkeg

import (
	"database/sql"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Index interface assumes any source of relational data store
type Index interface {
}

type index struct {
	File string  // full path to the sqlite3 DB file
	DB   *sql.DB // open handle to the DB in memory
}

// OpenIndex creates a Index implementation that is backed by a sqlite3
// database file at location specified by path plus dex/nodes.db. If it
// does not exist it creates and returns a reference to it. Otherwise,
// just opens and returns a new reference to it. If the file exists the
// proper organization is assumed. File is set to the full path to the
// dex/nodex.db file itself. DB is set to an opened sqlite3 *sql.DB handle.
func OpenIndex(kegpath string) (*index, error) {

	index := new(index)
	index.File = filepath.Join(kegpath, `dex`, `nodes.db`)

	db, err := sql.Open(`sqlite`, index.File)
	if err != nil {
		return nil, err
	}

	if !tableExists(db, `nodes`) {
		if _, err := db.Exec(`
		create table nodes(
			id integer primary key autoincrement,
			title text,
			created datetime default current_timestamp,
			updated datetime default current_timestamp
		);
	`); err != nil {
			return nil, err
		}
	}

	index.DB = db
	return index, nil
}

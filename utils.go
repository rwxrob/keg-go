package libkeg

import (
	"database/sql"
	"log"
	"os"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func notExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func tableExists(db *sql.DB, name string) bool {
	rows, err := db.Query(`
		select name from sqlite_master
		where type='table' and name=? ;
	`, name)
	if err != nil {
		log.Fatal(err)
	}
	rows.Next()
	var first string
	rows.Scan(&first)
	return first == name
}

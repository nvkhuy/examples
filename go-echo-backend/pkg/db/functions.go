package db

import (
	_ "embed"
)

//go:embed sql/count_element.sql
var countElement string

//go:embed sql/array_unique.sql
var arrayUniq string

//go:embed sql/create_unaccent_extensions.sql
var createUnaccent string

func (db *DB) setupFunctions() {

	db.Exec(countElement)
	db.Exec(createUnaccent)
	db.Exec(arrayUniq)
}

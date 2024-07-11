package db

func (db *DB) setupExtensions() {
	db.Exec("CREATE EXTENSION IF NOT EXISTS citext;")
}

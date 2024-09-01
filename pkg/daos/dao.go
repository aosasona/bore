package daos

import "database/sql"

type Dao struct {
	db *sql.DB
}

func New(dbConnection *sql.DB) *Dao {
	return &Dao{db: dbConnection}
}

package db

import (
	"database/sql"
)

type PostgreSQL struct {
	Db *sql.DB
}

func NewPostgreSqlDbConnection(connectionString string) *PostgreSQL {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &PostgreSQL{
		Db: db,
	}
}

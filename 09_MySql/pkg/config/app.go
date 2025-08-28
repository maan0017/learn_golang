package config

import (
	"database/sql"
	"fmt"
)

var (
	db *sql.DB
)

func ConnectToDB() {
	dsn := "username:password@tcp(127.0.0.1:3306)/testdb"

	d, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// Test connection
	if err := d.Ping(); err != nil {
		panic(err)
	}

	db = d
	fmt.Println("Connected to MySQL successfully!")
}

func GetDB() *sql.DB {
	if db == nil {
		return nil
	}

	return db
}

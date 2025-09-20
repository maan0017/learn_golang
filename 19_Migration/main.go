package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
)

func main() {
	if len(os.Args) < 2 { // os.Args ?
		log.Fatal("please provide a migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]
	fmt.Println(direction)

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}

	_migrate, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		if err := _migrate.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("[ERROR]::Migration Direction Error >", err)
		}
	case "down":
		if err := _migrate.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("[ERROR]::Migration Direction Error >", err)
		}
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}
}

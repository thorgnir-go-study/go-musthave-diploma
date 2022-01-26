package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"os"
)

func main() {
	connStr, ok := os.LookupEnv("DATABASE_URI")
	if !ok {
		log.Fatal("database uri not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All migrations done")
	fmt.Println("All migrations done")
}

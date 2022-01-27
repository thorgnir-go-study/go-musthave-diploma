package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
)

func main() {
	//connStr, ok := os.LookupEnv("DATABASE_URI")
	var connStr string
	flag.StringVar(&connStr, "database_uri", "", "blabla")
	//if !ok {
	//	log.Fatal("database uri not set")
	//}
	flag.Parse()

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

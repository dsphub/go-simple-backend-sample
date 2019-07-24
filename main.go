package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	. "github.com/dsphub/go-simple-crud-sample/store"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "any-password"
	dbname   = "crud"
)

func main() {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	server := NewPostServer(NewPostgresPostStore(db))
	fmt.Println("Start service")

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

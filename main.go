package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	. "github.com/dsphub/go-simple-crud-sample/store"
	_ "github.com/lib/pq"
)

func main() {
	connInfo := initOptions()
	postStore, err := NewPostgresPostStore(connInfo)
	if err != nil {
		log.Panic(err)
	}

	server := NewPostServer(postStore)
	fmt.Println("Start service")

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
	//FIXIT close db
}

type options struct {
	host       *string
	portNumber *int
	user       *string
	password   *string
	dbname     *string
	ssl        *bool
}

func initOptions() string {
	opts := &options{}
	opts.host = flag.String("host", "localhost", "service host name")
	opts.portNumber = flag.Int("port", 5432, "service port number")
	opts.dbname = flag.String("dbname", "crud", "db name")
	opts.user = flag.String("user", "postgres", "db user")
	opts.password = flag.String("password", "", "db password")
	opts.ssl = flag.Bool("ssl", false, "db ssl support")

	port := strconv.Itoa(*opts.portNumber)

	var sslmode string
	if *opts.ssl {
		sslmode = "enable"
	} else {
		sslmode = "disable"
	}

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", *opts.host, port, *opts.user, *opts.password, *opts.dbname, sslmode)
	return dbinfo
}

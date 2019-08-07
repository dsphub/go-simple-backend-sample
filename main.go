package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

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

const domainName = "localhost"
const httpServerPort = "5000"
const logFileName = "log.out"

func main() {
	log.Println("Start service")

	log := initLogger(logFileName)
	connInfo := initOptions(log)
	store := initStore(log, connInfo)
	server := NewPostServer(log, store)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", domainName, httpServerPort), server); err != nil {
		store.Disconnect()
		log.Fatalf("could not listen on port %s %v", httpServerPort, err)
	}
	waitTerminateSignal(log, store)
}

func initStore(log *log.Logger, connInfo string) *PostgresPostStore {
	postStore, err := NewPostgresPostStore(connInfo)
	if err != nil {
		log.Panic(err)
	}
	if err := postStore.Connect(); err != nil {
		log.Panic(err)
	}
	return postStore
}

type options struct {
	host       *string
	portNumber *int
	user       *string
	password   *string
	dbname     *string
	ssl        *bool
}

func initOptions(log *log.Logger) string {
	log.Println("Parse command-line options")
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

func initLogger(fileName string) *log.Logger {
	if fileName != "" {
		log.Println("Create log file")

		filePath, err:= getLogFilePath()
		if err != nil {
			panic(err)
		}

		logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		
		return log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func getLogFilePath() (string, error) {
	projectPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return projectPath + string(filepath.Separator)+logFileName, nil
}

func waitTerminateSignal(log *log.Logger, store *PostgresPostStore) {
	// After setting everything up!
	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("Received an interrupt, stopping services...")
		store.Disconnect()
		close(cleanupDone)
	}()
	<-cleanupDone
}

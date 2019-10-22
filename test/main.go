package test

import (
	"os"
	"time"

	"github.com/communitybridge/ledger/api"
	"github.com/communitybridge/ledger/gen/restapi"
	log "github.com/communitybridge/ledger/logging"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	// BaseURL for all endpoints
	BaseURL  = "http://localhost:8080/api/"
	testPort = 8080
)

// InitTestDB ...
func initTestDB() (*sqlx.DB, error) {
	log.Println("Initializing Test DB")

	db, err := sqlx.Connect("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		log.Fatal("err", err)
		return nil, err
	}
	db.SetMaxOpenConns(2)

	return db, nil
}

// Runs instance of api just for tests
func init() {

	log.SetLogLevel(0)

	// DB setup
	pDB, err := initTestDB()
	if err != nil {
		log.Fatal("couldn't connect to database", err)
	}

	api := api.ConfigureAPI(pDB)

	go func() {
		server := restapi.NewServer(api)
		defer func() {
			err := server.Shutdown()
			if err != nil {
				log.Printf("Error with server.Shutdown(): %s", err)
				log.Fatal(err)
			}
		}()

		server.Port = testPort
		err := server.Serve()
		if err == nil {
			logrus.Panicln(err)
		}
	}()

	time.Sleep(2 * time.Second)
}

package test

import (
	"log"
	"time"

	"github.com/communitybridge/ledger/api"
	"github.com/communitybridge/ledger/gen/restapi"
	"github.com/sirupsen/logrus"
)

const (
	// BaseURL for all endpoints
	BaseURL  = "http://localhost:8080/api/"
	testPort = 8080
)

// Runs instance of api just for tests
func init() {

	// DB setup
	pDB, err := api.InitDB()
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

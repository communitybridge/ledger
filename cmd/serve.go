// +build !aws_lambda

package cmd

import (
	"github.com/communitybridge/ledger/gen/restapi"
	"github.com/communitybridge/ledger/gen/restapi/operations"
	log "github.com/communitybridge/ledger/logging"
)

// Start the API server
func Start(api *operations.LedgerAPI, portFlag int) error {
	server := restapi.NewServer(api)
	defer func() {
		err := server.Shutdown()
		if err != nil {
			log.Printf("Error with server.Shutdown(): %s", err)
			log.Fatal(err)
		}
	}()

	log.Info("Starting Standard")
	server.Port = portFlag
	return server.Serve()
}

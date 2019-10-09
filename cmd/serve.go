// +build !aws_lambda

package cmd

import (
	"log"

	"github.com/communitybridge/ledger/gen/restapi"
	"github.com/communitybridge/ledger/gen/restapi/operations"
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

	server.Port = portFlag
	return server.Serve()
}

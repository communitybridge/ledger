package api

import (
	"github.com/communitybridge/ledger/balance"
	"github.com/communitybridge/ledger/gen/restapi"
	"github.com/communitybridge/ledger/gen/restapi/operations"
	"github.com/communitybridge/ledger/health"
	log "github.com/communitybridge/ledger/logging"
	"github.com/communitybridge/ledger/transaction"
	"github.com/go-openapi/loads"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // import postgres driver
)

var (
	// BuildStamp is a timestamp (injected by go) of the build time
	BuildStamp = "None"
	// GitHash is the tag for current hash the build represents
	GitHash = "None"
)

// ConfigureAPI ...
func ConfigureAPI(pDB *sqlx.DB) *operations.LedgerAPI {
	// loads generated Swagger API specifications
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal("Invalid swagger file for initializing", err)
	}
	api := operations.NewLedgerAPI(swaggerSpec)

	// Health setup
	healthService := health.New()
	health.Configure(api, healthService)

	// Transactions package endpoints
	transactionRepo := transaction.NewRepository(pDB)
	transactionService := transaction.NewService(transactionRepo)
	transaction.Configure(api, transactionService)

	// Balance package endpoints
	balanceRepo := balance.NewRepository(pDB)
	balanceService := balance.NewService(balanceRepo)
	balance.Configure(api, balanceService)

	return api
}

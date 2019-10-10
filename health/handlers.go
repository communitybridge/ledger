package health

import (
	"github.com/communitybridge/ledger/gen/restapi/operations"
	"github.com/communitybridge/ledger/gen/restapi/operations/health"
	log "github.com/communitybridge/ledger/logging"
	"github.com/communitybridge/ledger/swagger"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with Service
func Configure(api *operations.LedgerAPI, service Service) {

	api.HealthGetHealthHandler = health.GetHealthHandlerFunc(func(params health.GetHealthParams) middleware.Responder {
		log.Info("entered GetHealthHandler")
		result, err := service.GetHealth(params.HTTPRequest.Context())
		if err != nil {
			return swagger.HealthErrorHandler("GetHealth", err)
		}
		return health.NewGetHealthOK().WithPayload(result)
	})

}

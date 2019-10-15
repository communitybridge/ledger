package balance

import (
	"fmt"

	"github.com/communitybridge/ledger/gen/restapi/operations/balance"

	"github.com/communitybridge/ledger/gen/restapi/operations"
	log "github.com/communitybridge/ledger/logging"
	"github.com/communitybridge/ledger/swagger"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"
)

// Configure setups handlers on api with Service
func Configure(api *operations.LedgerAPI, service Service) {

	api.BalanceGetBalanceHandler = balance.GetBalanceHandlerFunc(func(params balance.GetBalanceParams) middleware.Responder {
		log.Info("entering BalanceGetBalanceHandler")

		log.WithFields(logrus.Fields{
			"HttpRequest": fmt.Sprintf("%#v", *params.HTTPRequest),
			"EntityID":    params.EntityID,
		}).Info("GetTransactionHandlerFunc")

		result, err := service.GetEntityBalance(params.HTTPRequest.Context(), &params)
		if err != nil {
			return swagger.BalanceErrorHandler("GetBalance", err)
		}
		return balance.NewGetBalanceOK().WithPayload(result)
	})

}

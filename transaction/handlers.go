package transaction

import (
	"fmt"

	"github.com/communitybridge/ledger/gen/restapi/operations/transactions"

	"github.com/communitybridge/ledger/gen/restapi/operations"
	log "github.com/communitybridge/ledger/logging"
	"github.com/communitybridge/ledger/swagger"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"
)

// Configure setups handlers on api with Service
func Configure(api *operations.LedgerAPI, service Service) {

	api.TransactionsListTransactionsHandler = transactions.ListTransactionsHandlerFunc(func(params transactions.ListTransactionsParams) middleware.Responder {
		log.Info("entering ListTransactionsHandler")

		log.WithFields(logrus.Fields{
			"Offset":   *params.Offset,
			"PageSize": *params.PageSize,
		}).Info("ListTransactionsHandler")

		log.Info(fmt.Sprintf("{URL: %#v}",
			*params.HTTPRequest.URL))

		result, err := service.ListTransactions(params.HTTPRequest.Context(), &params)
		if err != nil {
			return swagger.TransactionErrorHandler("ListTransactions", err)
		}
		return transactions.NewListTransactionsOK().WithPayload(result)
	})

	api.TransactionsCreateTransactionHandler = transactions.CreateTransactionHandlerFunc(func(params transactions.CreateTransactionParams) middleware.Responder {
		log.Info("entering transactionsCreateTransactionHandler")

		log.WithFields(logrus.Fields{
			"HttpRequest":       fmt.Sprintf("%#v", *params.HTTPRequest),
			"createTransaction": fmt.Sprintf("%#v", *params.Transaction),
		}).Info("CreateTransactionHandlerFunc")

		result, err := service.CreateTransaction(params.HTTPRequest.Context(), &params)
		if err != nil {
			return swagger.TransactionErrorHandler("CreateTransaction", err)
		}
		return transactions.NewCreateTransactionCreated().WithPayload(result)
	})
}

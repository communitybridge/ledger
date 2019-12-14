package docs

import (
	"github.com/communitybridge/ledger/gen/restapi/operations"
	"github.com/communitybridge/ledger/gen/restapi/operations/docs"
	"github.com/go-openapi/runtime/middleware"
)

// Configure setups handlers on api with Service
func Configure(api *operations.LedgerAPI) {

	api.DocsGetDocHandler = docs.GetDocHandlerFunc(func(params docs.GetDocParams) middleware.Responder {
		return NewGetDocOK()
	})

	api.DocsGetSwaggerHandler = docs.GetSwaggerHandlerFunc(func(params docs.GetSwaggerParams) middleware.Responder {
		return NewGetSwaggerOK()
	})
}

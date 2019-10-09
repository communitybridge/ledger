// +build aws_lambda

package cmd

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/communitybridge/ledger/gen/restapi/operations"
	"github.com/sirupsen/logrus"
)

func Start(api *operations.LedgerAPI, _ int) error {
	adapter := httpadapter.New(api.Serve(nil))

	logrus.Debug("Starting Lambda")
	lambda.Start(adapter.Proxy)
	return nil
}

package main

import (
	"flag"
	"log"
	"os"

	"github.com/communitybridge/ledger/balance"
	"github.com/communitybridge/ledger/cmd"
	"github.com/communitybridge/ledger/docs"
	"github.com/communitybridge/ledger/gen/restapi"
	"github.com/communitybridge/ledger/gen/restapi/operations"
	"github.com/communitybridge/ledger/health"
	"github.com/communitybridge/ledger/transaction"
	"github.com/go-openapi/loads"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// BuildStamp is a timestamp (injected by go) of the build time
	BuildStamp = "None"
	// GitHash is the tag for current hash the build represents
	GitHash = "None"
)

func main() {

	host, err := os.Hostname()
	if err != nil {
		logrus.Panicln("unable to get Hostname", err)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.WithFields(logrus.Fields{
		"BuildTime": BuildStamp,
		"GitHash":   GitHash,
		"Host":      host,
	}).Info("Start Service")

	// Configures Viper, the configuration management tool, and set some app defaults
	viperConfig := viper.New()
	viperConfig.AutomaticEnv()
	viperConfig.SetEnvPrefix("LS") // this prefix is specific to the Ledger Service
	defaults := map[string]interface{}{
		"PORT":     8080,
		"USE_MOCK": "False",
	}
	for key, value := range defaults {
		viperConfig.SetDefault(key, value)
	}

	// DB setup
	pDB, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.Fatal("err", err)
	}
	pDB.SetMaxOpenConns(3)

	// loads generated Swagger API specifications
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal("Invalid swagger file for initializing", err)
	}
	api := operations.NewLedgerAPI(swaggerSpec)

	docs.Configure(api)

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

	var portFlag = flag.Int("port", viperConfig.GetInt("PORT"), "Port to listen for web requests on")
	flag.Parse()

	logrus.Info("Start Service")
	if err := cmd.Start(api, *portFlag); err != nil {
		logrus.Panicln(err)
	}
}

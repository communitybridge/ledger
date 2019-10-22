package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/imroc/req"
	. "github.com/smartystreets/goconvey/convey"
)

// transaction data
var entityID = "b582a786-48ec-469b-b655-17cf779b9ce1"
var entityType = "project"
var asset = "usd"
var externalTransactionID = "ex1234abcid"
var externalSourceType = "bill.com"
var transactionCategory = "donation"
var externalAccountID = "exaccountid1234"

// transaction line item data
var lineItemAmountOne = 1500
var lineItemDescriptionOne = "donation"

var lineItemAmountTwo = -500
var lineItemDescriptionTwo = "fee"

type LineItem struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type CreateTransaction struct {
	EntityID              string     `json:"entity_id"`
	EntityType            string     `json:"entity_type"`
	ExternalTransactionID string     `json:"external_transaction_id"`
	Asset                 string     `json:"asset"`
	ExternalSourceType    string     `json:"external_source_type"`
	TransactionCategory   string     `json:"transaction_category"`
	ExternalAccountID     string     `json:"external_account_id"`
	LineItems             []LineItem `json:"line_items"`
}

func GetCreateTransactionPayload() CreateTransaction {
	lineItemOne := LineItem{}
	lineItemOne.Amount = lineItemAmountOne
	lineItemOne.Description = lineItemDescriptionOne

	lineItemTwo := LineItem{}
	lineItemTwo.Amount = lineItemAmountTwo
	lineItemTwo.Description = lineItemDescriptionTwo

	createTransaction := CreateTransaction{}
	createTransaction.EntityID = entityID
	createTransaction.EntityType = entityType
	createTransaction.ExternalTransactionID = externalTransactionID
	createTransaction.Asset = asset
	createTransaction.ExternalSourceType = externalSourceType
	createTransaction.TransactionCategory = transactionCategory
	createTransaction.ExternalAccountID = externalAccountID
	createTransaction.LineItems = []LineItem{lineItemOne, lineItemTwo}

	return createTransaction
}

func TestCreateTransactionEndpoint(t *testing.T) {
	Convey("Given API is running", t, func() {

		Convey("When the transactions endpoint is hit with valid POST data to create a new Transaction", func() {

			createTransaction := GetCreateTransactionPayload()
			json, err := json.Marshal(createTransaction)
			if err != nil {
				fmt.Println(err)
				return
			}

			url := fmt.Sprintf("%stransactions", BaseURL)
			header := req.Header{
				"Content-Type": "application/json",
			}
			resp, err := req.Post(url, header, json)
			fmt.Println(resp)
			if err != nil {
				t.Error("Response: ", resp.String())
				t.Fail()
			}

			Convey("It will get 201 status", func() {
				So(resp.Response().StatusCode, ShouldEqual, 201)
			})

			transaction := models.Transaction{}
			err = resp.ToJSON(&transaction)
			if err != nil {
				t.Error("Error: ", err.Error())
			}

			fmt.Println(fmt.Sprintf("#%v", transaction))

			Convey("It will get the specified transaction values", func() {
				So(transaction.ExternalTransactionID, ShouldEqual, externalTransactionID)
				So(transaction.TransactionCategory, ShouldEqual, transactionCategory)
				So(transaction.Asset, ShouldEqual, asset)
			})

			Convey("It will get the specified line_item values", func() {
				So(transaction.LineItems[0].Amount, ShouldEqual, lineItemAmountOne)
				So(transaction.LineItems[0].Description, ShouldEqual, lineItemDescriptionOne)
				So(transaction.LineItems[1].Amount, ShouldEqual, lineItemAmountTwo)
				So(transaction.LineItems[1].Description, ShouldEqual, lineItemDescriptionTwo)
			})
		})

		Convey("When the transactions endpoint is hit with invalid POST data to create a new Transaction", func() {

			createTransaction := GetCreateTransactionPayload()

			// EntityID should be valid UUID
			createTransaction.EntityID = ""

			json, err := json.Marshal(createTransaction)
			if err != nil {
				fmt.Println(err)
				return
			}

			url := fmt.Sprintf("%stransactions", BaseURL)
			header := req.Header{
				"Content-Type": "application/json",
			}
			resp, err := req.Post(url, header, json)
			fmt.Println(resp)
			if err != nil {
				t.Error("Response: ", resp.String())
				t.Fail()
			}

			Convey("It will get 400 status", func() {
				So(resp.Response().StatusCode, ShouldEqual, 400)
			})

		})

	})
}

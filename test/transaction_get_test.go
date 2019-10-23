package test

import (
	"fmt"
	"testing"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/imroc/req"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTransactionEndpoint(t *testing.T) {
	Convey("Given API is running", t, func() {

		// Transactionent with ID fd54832d-d872-428b-a10d-17ddf782b4df is seeded in testDD.
		var transactionID = "fd54832d-d872-428b-a10d-17ddf782b4df"
		var invalidTransactionID = "61b0c143-f1f9-not-real-id"

		var transaction models.Transaction

		Convey("When user requests payments/<payment-id> endpoint with an invalid payment ID", func() {

			url := fmt.Sprintf("%stransactions/%s", BaseURL, invalidTransactionID)
			header := req.Header{
				"Content-Type": "application/json",
			}
			resp, err := req.Get(url, header)
			fmt.Println(resp)
			if err != nil {
				t.Error("Response: ", resp.String())
				t.Fail()
			}

			Convey("It will get 404 status", func() {
				So(resp.Response().StatusCode, ShouldEqual, 404)
			})

		})

		var accountID = "6eae6bb8-f7fb-425a-8af8-64adb616b739"
		var transactionCategory = "random"
		var externalTransactionID = "a04c291f-234567"

		Convey("When user requests transactions/<transaction-id> endpoint with a valid transaction ID", func() {

			url := fmt.Sprintf("%stransactions/%s", BaseURL, transactionID)
			header := req.Header{
				"Content-Type": "application/json",
			}
			resp, err := req.Get(url, header)
			fmt.Println(resp)
			if err != nil {
				t.Error("Response: ", resp.String())
				t.Fail()
			}

			Convey("It will get 404 status", func() {
				So(resp.Response().StatusCode, ShouldEqual, 200)
			})

			err = resp.ToJSON(&transaction)
			if err != nil {
				t.Error("Error: ", err.Error())
			}

			Convey("It will get the referenced payment values", func() {
				So(transaction.ID, ShouldEqual, transactionID)
				So(transaction.AccountID, ShouldEqual, accountID)
				So(transaction.TransactionCategory, ShouldEqual, transactionCategory)
				So(transaction.ExternalTransactionID, ShouldEqual, externalTransactionID)
			})

			Convey("It will get the transaction LineItems", func() {
				lineItems := transaction.LineItems
				So(len(lineItems), ShouldEqual, 3)

				lineItemOne := *lineItems[0]
				So(lineItemOne.Amount, ShouldEqual, -1000)
			})

		})

	})
}

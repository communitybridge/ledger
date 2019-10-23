package test

import (
	"fmt"
	"testing"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/imroc/req"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTransactionsEndpoint(t *testing.T) {
	Convey("Given API is running", t, func() {

		var transactions = models.TransactionList{}

		Convey("When user access list transactions endpoint", func() {

			url := fmt.Sprintf("%stransactions", BaseURL)
			header := req.Header{
				"Content-Type": "application/json",
			}
			resp, err := req.Get(url, header)
			fmt.Println(resp)
			if err != nil {
				t.Error("Response: ", resp.String())
				t.Fail()
			}

			Convey("It will get 200 status", func() {
				So(resp.Response().StatusCode, ShouldEqual, 200)
			})

			Convey("It will get array of transactions with length equals 1", func() {
				err = resp.ToJSON(&transactions)
				if err != nil {
					t.Error("Error: ", err.Error())
				}

				So(len(transactions.Data), ShouldEqual, 1)
			})

		})

	})
}

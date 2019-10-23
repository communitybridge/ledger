package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetEntityBalanceEndpoint(t *testing.T) {
	Convey("Given API is running", t, func() {

		// var entityType = "project"
		// var availableBalance = 850
		// var totalCount = 4
		// var totalCredit = 1200
		// var creditCount = 2
		// var totalDebit = 350
		// var debitCount = 2

		// // entity with external ID is seeded in testDD.
		// var entityID = "b582a786-48ec-469b-b655-17cf779b9ce1"

		// var balance models.Balance

		// Convey("When user requests balance for an entity with external ID", func() {

		// 	url := fmt.Sprintf("%sbalance/%s", BaseURL, entityID)
		// 	header := req.Header{
		// 		"Content-Type": "application/json",
		// 	}
		// 	resp, err := req.Get(url, header)
		// 	fmt.Println(resp)
		// 	if err != nil {
		// 		t.Error("Response: ", resp.String())
		// 		t.Fail()
		// 	}

		// 	err = resp.ToJSON(&balance)
		// 	if err != nil {
		// 		t.Error("Error: ", err.Error())
		// 	}

		// 	Convey("It will get the referenced balance values", func() {
		// 		So(balance.EntityType, ShouldEqual, entityType)
		// 		So(balance.AvailableBalance, ShouldEqual, availableBalance)
		// 		So(balance.TotalCount, ShouldEqual, totalCount)
		// 		So(balance.TotalCredit, ShouldEqual, totalCredit)
		// 		So(balance.CreditCount, ShouldEqual, creditCount)
		// 		So(balance.TotalDebit, ShouldEqual, totalDebit)
		// 		So(balance.DebitCount, ShouldEqual, debitCount)
		// 	})
		// })
	})
}

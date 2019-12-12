package transaction

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/communitybridge/ledger/gen/restapi/operations/transactions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	asset = "usd"
)

func TestTransactionGetByID(t *testing.T) {
	fmt.Println(os.Getenv("DATABASE_URL"))
	conn, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	require.Nil(t, err)
	defer conn.Close()
	prepareTestDatabase()

	transactionID := "61b0c143-f1f9-457d-a889-80570b820348"
	accountID := "5701249e-f33a-45a3-8722-e6917ccff6f0"
	externalTransactionID := "a04c291f-234567"

	transactionRepository := NewRepository(conn)
	transaction, err := transactionRepository.GetTransaction(context.Background(), transactionID)
	assert.Nil(t, err)
	assert.Equal(t, transactionID, transaction.ID)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, asset, transaction.Asset)
	assert.Equal(t, externalTransactionID, transaction.ExternalTransactionID)
}

func TestTransactionGetTransactionCount(t *testing.T) {
	conn, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	require.Nil(t, err)
	defer conn.Close()
	prepareTestDatabase()

	numberOfTransactions := int64(2)

	transactionRepository := NewRepository(conn)
	transactionCount, err := transactionRepository.GetTransactionCount(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, transactionCount, numberOfTransactions)
}

func TestTransactionListTransactions(t *testing.T) {
	conn, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	require.Nil(t, err)
	defer conn.Close()
	prepareTestDatabase()

	numberOfTransactions := 2
	pageSize := int64(10)
	offset := int64(0)
	orderBy := "created_at"
	listTransactionsParams := transactions.ListTransactionsParams{}
	listTransactionsParams.Offset = &offset
	listTransactionsParams.PageSize = &pageSize
	listTransactionsParams.OrderBy = &orderBy

	transactionRepository := NewRepository(conn)
	transactions, err := transactionRepository.ListTransactions(context.Background(), &listTransactionsParams)
	assert.Nil(t, err)
	assert.Equal(t, len(transactions), numberOfTransactions)
}

func TestTransactionCreateTransaction(t *testing.T) {
	conn, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	require.Nil(t, err)
	defer conn.Close()
	prepareTestDatabase()

	// transaction data
	entityID := "b582a786-48ec-469b-b655-17cf779b9ce1"
	entityType := "project"
	externalSourceType := "bill.com"
	transactionCategory := "unknown"
	exAccountID := "exaccountid1234"

	// transaction line item data
	lineItemAmountOne := int64(1500)
	lineItemDescriptionOne := "donation"

	lineItemAmountTwo := int64(-500)
	lineItemDescriptionTwo := "fee"

	lineItemOne := models.CreateLineItem{}
	lineItemOne.Amount = &lineItemAmountOne
	lineItemOne.Description = &lineItemDescriptionOne

	lineItemTwo := models.CreateLineItem{}
	lineItemTwo.Amount = &lineItemAmountTwo
	lineItemTwo.Description = &lineItemDescriptionTwo

	createTransaction := models.CreateTransaction{}
	createTransaction.EntityID = &entityID
	createTransaction.EntityType = &entityType
	createTransaction.ExternalTransactionID = &exAccountID
	createTransaction.Asset = asset
	createTransaction.ExternalSourceType = &externalSourceType
	createTransaction.TransactionCategory = transactionCategory
	createTransaction.ExternalAccountID = &exAccountID
	createTransaction.LineItems = []*models.CreateLineItem{}
	createTransaction.LineItems = append(createTransaction.LineItems, &lineItemOne)
	createTransaction.LineItems = append(createTransaction.LineItems, &lineItemTwo)

	transactionRepository := NewRepository(conn)
	transaction, err := transactionRepository.CreateTransaction(context.Background(), &createTransaction)
	assert.Nil(t, err)
	assert.Equal(t, transaction.Asset, createTransaction.Asset)
	assert.Equal(t, transaction.ExternalTransactionID, exAccountID)
	assert.Equal(t, len(transaction.LineItems), len(createTransaction.LineItems))

	for _, lineItem := range transaction.LineItems {
		if lineItem.Description == lineItemDescriptionOne {
			assert.Equal(t, lineItem.Amount, lineItemAmountOne)
		} else {
			assert.Equal(t, lineItem.Amount, lineItemAmountTwo)
		}
	}

}

package transaction

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/communitybridge/ledger/gen/restapi/operations/transactions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	transactionID         = "61b0c143-f1f9-457d-a889-80570b820348"
	accountID             = "5701249e-f33a-45a3-8722-e6917ccff6f0"
	asset                 = "usd"
	externalTransactionID = "a04c291f-234567"
)

func TestTransactionGetByID(t *testing.T) {
	fmt.Println(os.Getenv("DATABASE_URL"))
	conn, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	require.Nil(t, err)
	defer conn.Close()
	prepareTestDatabase()

	transactionRepository := NewRepository(conn)
	transaction, err := transactionRepository.GetTransaction(context.Background(), transactionID)
	assert.Nil(t, err)
	assert.Equal(t, transactionID, transaction.ID)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, asset, transaction.Asset)
	assert.Equal(t, externalTransactionID, transaction.ExternalTransactionID)
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

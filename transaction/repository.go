package transaction

import (
	"context"
	"fmt"
	"strconv"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/communitybridge/ledger/gen/restapi/operations/transactions"
	"github.com/communitybridge/ledger/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	log "github.com/communitybridge/ledger/logging"
)

type Repository interface {
	ListTransactions(ctx context.Context, params *transactions.ListTransactionsParams) ([]*models.Transaction, error)
	CreateTransaction(ctx context.Context, params *models.CreateTransaction) (*models.Transaction, error)
	GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) GetDB() *sqlx.DB {
	return repo.db
}

// DoesAssetExist checks if a given currency exists
func DoesAssetExist(repo *repository, abbrv *string) (bool, string) {
	log.Info("entered function DoesAssetExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM asset WHERE abbrv=$1", abbrv)
	if err != nil {
		err = fmt.Errorf("asset with abbreviation : `%s` does not exist", *abbrv)
		log.Info(err.Error())
		return false, res
	}

	return true, res
}

// DoesAccountExist checks if a given account exists
func DoesAccountExist(repo *repository, id string) (bool, error) {
	log.Info("entered function DoesAccountExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM accounts WHERE id=$1", id)
	if err != nil {
		err = fmt.Errorf("account with id : `%s` does not exist", id)
		log.Info(err.Error())
		return false, err
	}

	return true, nil
}

// DoesTransactionExist checks if a given transaction exists
func DoesTransactionExist(repo *repository, id string) (bool, string) {
	log.Info("entered function DoesTransactionExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM t Wransactions WHERE id=$1", id)
	if err != nil {
		err = fmt.Errorf("transaction with id : `%s` does not exist", id)
		log.Info(err.Error())
		return false, res
	}

	return true, res
}

// ListTransactions is a function to get a list of transactions
func (repo *repository) ListTransactions(ctx context.Context, params *transactions.ListTransactionsParams) ([]*models.Transaction, error) {
	log.Info("entered function ListTransactions")

	pagesize, err := strconv.Atoi(*params.PageSize)
	if err != nil {
		log.Fatal(log.Trace(), err)
		return nil, errors.Wrap(err, "ListTransactions.convertPageSize")
	}

	offset, err := strconv.Atoi(*params.Offset)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, errors.Wrap(err, "ListTransactions.convertOffset")
	}

	sql := `
		SELECT
			t.id AS ID
		FROM
			transaction t
		limit $1
		offset $2;`

	log.Info(fmt.Sprintf(log.StripSpecialChars(sql), pagesize, offset))

	rows, err := repo.db.Queryx(sql,
		pagesize,
		offset)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, errors.Wrap(err, "ListTransactions.query")
	}

	var transactions []*models.Transaction
	for rows.Next() {

		transaction := &models.Transaction{}

		if err := rows.Scan(
			&transaction.ID,
		); err != nil {

			log.Error(err.Error(), err)
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// GetTransaction is function to get a specific transaction
func (repo *repository) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	log.Info("entered function GetTransaction")

	exist, _ := DoesTransactionExist(repo, transactionID)
	if !exist {
		return nil, swagger.ErrNotFound
	}

	sql := `
	SELECT
		t.id AS ID
	FROM
		transaction t
	WHERE
		id = $1;`

	log.Info(log.StripSpecialChars(sql))

	row := repo.db.QueryRow(sql, transactionID)

	transaction := &models.Transaction{}
	if err := row.Scan(
		&transaction.ID,
	); err != nil {
		log.Error(err.Error(), err)
		return nil, err
	}

	return transaction, nil
}

// CreateTransaction creates a new transaction and any related rows
// in required tables if they don't already exist
func (repo *repository) CreateTransaction(ctx context.Context, params *models.CreateTransaction) (*models.Transaction, error) {

	log.Info("entered function CreateTransaction")

	transactionResp := models.Transaction{}

	return &transactionResp, nil
}

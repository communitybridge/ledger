package transaction

import (
	"context"
	"fmt"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/communitybridge/ledger/gen/restapi/operations/transactions"
	"github.com/communitybridge/ledger/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"

	log "github.com/communitybridge/ledger/logging"
)

// Repository interface for repo calls
type Repository interface {
	ListTransactions(ctx context.Context, params *transactions.ListTransactionsParams) ([]*models.Transaction, error)
	GetTransactionCount(ctx context.Context) (int64, error)
	CreateTransaction(ctx context.Context, params *models.CreateTransaction) (*models.Transaction, error)
	GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository ...
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) GetDB() *sqlx.DB {
	return repo.db
}

// DoesAssetExist checks if a given asset exists
func DoesAssetExist(repo *repository, id string) (bool, error) {
	log.Info("entered function DoesAssetExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM assets WHERE id=$1", id)
	if err != nil {
		err = fmt.Errorf("asset with id : `%s` does not exist", id)
		log.Info(err.Error())
		return false, err
	}

	return true, nil
}

// DoesAccountExist checks if a given account exists
func DoesAccountExist(repo *repository, id string) bool {
	log.Info("entered function DoesAccountExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM accounts WHERE id=$1", id)
	if err != nil {
		err = fmt.Errorf("account with id : `%s` does not exist", id)
		log.Info(err.Error())
		return false
	}

	return true
}

// DoesEntityExist checks if a given account exists
func DoesEntityExist(repo *repository, entityID string, entityType string) bool {
	log.Info("entered function DoesEntityExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM entities WHERE entity_id=$1 AND entity_type=$2", entityID, entityType)
	if err != nil {
		err = fmt.Errorf("entity with id: `%s` and type: '%s' does not exist", entityID, entityType)
		log.Info(err.Error())
		return false
	}

	return true
}

// DoesTransactionExist checks if a given transaction exists
func DoesTransactionExist(repo *repository, id string) (bool, string) {
	log.Info("entered function DoesTransactionExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM transactions WHERE id=$1", id)
	if err != nil {
		err = fmt.Errorf("transaction with id : `%s` does not exist", id)
		log.Info(err.Error())
		return false, res
	}

	return true, res
}

// getTransactionLineItems returns party data from required tables
func getTransactionLineItems(repo *repository, transactionID string) ([]*models.LineItem, error) {

	if transactionID == "" {
		err := fmt.Errorf("account id is empty")
		return []*models.LineItem{}, err
	}

	sql := `
		SELECT
			l.id AS ID,
			l.amount AS Amount,
			l.asset_id AS AssetID,
			l.metadata AS Metadata,
			l.created_at AS CreatedAt,
			l.updated_at AS UpdatedAt
		FROM
			line_items l
		WHERE
			l.transaction_id = $1;`

	log.Info(log.StripSpecialChars(sql))

	rows, err := repo.db.Queryx(sql, transactionID)
	if err != nil {
		log.Error(err.Error(), err)
		return nil, err
	}

	lineItems := []*models.LineItem{}
	for rows.Next() {

		lineItem := &models.LineItem{}
		err := rows.StructScan(&lineItem)
		if err != nil {
			log.Error(err.Error(), err)
			return nil, err
		}

		lineItems = append(lineItems, lineItem)
	}

	return lineItems, nil
}

// GetTransactionCount is a function to get a count of available transactions
func (repo *repository) GetTransactionCount(ctx context.Context) (int64, error) {
	log.Info("entered function GetTransactionCount")

	sql := `SELECT count(*) FROM transactions t;`

	log.Info(log.StripSpecialChars(sql))

	row := repo.db.QueryRow(sql)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ListTransactions is a function to get a list of transactions
func (repo *repository) ListTransactions(ctx context.Context, params *transactions.ListTransactionsParams) ([]*models.Transaction, error) {
	log.Info("entered function ListTransactions")

	pagesize := *params.PageSize
	offset := *params.Offset

	sql := `
		SELECT
			t.id AS ID,
			t.account_id AS AccountID,
			t.external_transaction_id AS ExternalTransactionID,
			t.metadata AS Metadata,
			t.running_balance AS RunningBalance,
			t.transaction_category AS TransactionCategory,
			t.created_at AS CreatedAt,
			t.updated_at AS UpdatedAt
		FROM
			transactions t
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
			&transaction.AccountID,
			&transaction.ExternalTransactionID,
			&transaction.Metadata,
			&transaction.RunningBalance,
			&transaction.TransactionCategory,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		); err != nil {
			log.Error(err.Error(), err)
			return nil, err
		}

		lineItems, err := getTransactionLineItems(repo, transaction.ID)
		if err != nil {
			log.Error(err.Error(), err)
			return nil, err
		}
		transaction.LineItems = lineItems

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
			t.id AS ID,
			t.account_id AS AccountID,
			t.external_transaction_id AS ExternalTransactionID,
			t.metadata AS Metadata,
			t.running_balance AS RunningBalance,
			t.transaction_category AS TransactionCategory,
			t.created_at AS CreatedAt,
			t.updated_at AS UpdatedAt
		FROM
			transactions t
		WHERE
			id = $1;`

	log.Info(log.StripSpecialChars(sql))

	row := repo.db.QueryRowx(sql, transactionID)

	transaction := models.Transaction{}
	err := row.StructScan(&transaction)
	if err != nil {
		log.Error(err.Error(), err)
		return nil, err
	}

	lineItems, err := getTransactionLineItems(repo, transaction.ID)
	if err != nil {
		log.Error(err.Error(), err)
		return nil, err
	}
	transaction.LineItems = lineItems

	return &transaction, nil
}

// CreateTransaction creates a new transaction and any related rows
// in required tables if they don't already exist
func (repo *repository) CreateTransaction(ctx context.Context, params *models.CreateTransaction) (*models.Transaction, error) {

	log.Info("entered function CreateTransaction")

	// Check if account exists
	accountExist := DoesAccountExist(repo, *params.AccountID)
	if !accountExist {
		err := fmt.Errorf("invalid account ID for this transaction %s", *params.AccountID)
		log.Error(log.Trace(), err)
		return nil, err
	}

	// Check if entity exists
	entityExist := DoesEntityExist(repo, *params.EntityID, *params.EntityType)
	if !entityExist {
		err := fmt.Errorf("invalid entity ID: %s or Type: %s", *params.EntityID, *params.EntityType)
		log.Error(log.Trace(), err)
		return nil, err
	}

	// Stub, replace.
	runningBalance := 1000

	metaDataJSON := types.JSONText(string(params.Metadata))
	metaDataJSONValue, err := metaDataJSON.Value()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new transaction entry
	sql := `
		INSERT INTO transactions (
			account_id,
			transaction_category,
			external_transaction_id,
			running_balance,
			metadata
		)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			id AS ID,
			account_id AS AccountID,
			external_transaction_id AS ExternalTransactionID,
			metadata AS Metadata,
			running_balance AS RunningBalance,
			transaction_category AS TransactionCategory,
			created_at AS CreatedAt,
			updated_at AS UpdatedAt;`

	log.Info(fmt.Sprintf(log.StripSpecialChars(sql),
		params.AccountID,
		params.TransactionCategory,
		params.ExternalTransactionID,
		runningBalance,
		params.Metadata,
	))

	// Begin a transaction
	tx, err := repo.db.Beginx()
	if err != nil {
		log.Fatal(log.Trace(), err)
	}
	defer func() {
		if err != nil {
			log.Error(log.Trace(), err)
			tx.Rollback()
			return
		}
	}()

	// Insert Statement
	row := tx.QueryRowx(sql,
		params.AccountID,
		params.TransactionCategory,
		params.ExternalTransactionID,
		runningBalance,
		metaDataJSONValue,
	)

	transaction := models.Transaction{}
	if err = row.StructScan(&transaction); err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	items := []*models.LineItem{}
	for _, item := range params.LineItems {

		// Check if asset exists
		var res = ""
		err := tx.Get(&res, "SELECT id FROM assets WHERE id=$1", *item.AssetID)
		if err != nil {
			err = fmt.Errorf("asset with id : `%s` does not exist", *item.AssetID)
			log.Error(log.Trace(), err)
			return nil, err
		}

		metaDataJSON := types.JSONText(string(item.Metadata))
		metaDataJSONValue, err := metaDataJSON.Value()
		if err != nil {
			log.Fatal(err)
		}

		// Inset Line Items
		sql = `
		INSERT INTO line_items (
			transaction_id,
			amount,
			description,
			asset_id,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			id AS ID,
			transaction_id AS TransactionID,
			amount AS Amount,
			description as Description,
			asset_id as AssetID,
			metadata AS Metadata,
			created_at AS CreatedAt,
			updated_at AS UpdatedAt;`

		log.Info(fmt.Sprintf(log.StripSpecialChars(sql),
			transaction.ID,
			*item.Amount,
			*item.Description,
			*item.AssetID,
			metaDataJSONValue,
		))

		// Insert LineItem
		row := tx.QueryRowx(sql,
			transaction.ID,
			*item.Amount,
			*item.Description,
			*item.AssetID,
			metaDataJSONValue,
		)

		lineItem := models.LineItem{}
		if err = row.StructScan(&lineItem); err != nil {
			log.Error(log.Trace(), err)
			return nil, err
		}

		items = append(items, &lineItem)
	}

	transaction.LineItems = items

	err = tx.Commit()
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	return &transaction, nil
}

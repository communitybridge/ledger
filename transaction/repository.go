package transaction

import (
	"context"
	"database/sql"
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

// CreateAccount creates a new account
func CreateAccount(repo *repository, params *models.CreateTransaction) (string, error) {
	log.Info("entered function CreateAccount")

	// Get entity.id from entity.entity_id
	entityID := GetExistingEntity(repo, *params.EntityID, *params.EntityType)
	if entityID == "" {
		err := errors.New(fmt.Sprintf("could not get entity.id for entity.entity_id: %s entity.entity_type: %s", *params.EntityID, *params.EntityType))
		log.Error(err.Error(), err)
		return "", err
	}

	sql := `
	INSERT INTO accounts (entity_id, external_source_type, external_account_id)
		VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING id;`

	log.Info(fmt.Sprintf(log.StripSpecialChars(sql),
		entityID,
		params.ExternalSourceType,
		params.ExternalAccountID))

	// Insert Statement
	row := repo.db.QueryRowx(sql,
		entityID,
		params.ExternalSourceType,
		params.ExternalAccountID,
	)

	var accountID = ""
	if err := row.Scan(&accountID); err != nil {
		log.Error(log.Trace(), err)
		return "", err
	}

	log.Info(fmt.Sprintf("added account with ID %s to table", accountID))

	return accountID, nil
}

// CreateEntity creates a new entity
func CreateEntity(repo *repository, params *models.CreateTransaction) (string, error) {
	log.Info("entered function CreateEntity")

	sql := `
	INSERT INTO entities (entity_id, entity_type)
		VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	RETURNING id;`

	log.Info(fmt.Sprintf(log.StripSpecialChars(sql),
		params.EntityID,
		params.EntityType))

	// Insert Statement
	row := repo.db.QueryRowx(sql,
		params.EntityID,
		params.EntityType,
	)

	var entityID = ""
	if err := row.Scan(&entityID); err != nil {
		log.Error(log.Trace(), err)
		return "", err
	}

	log.Info(fmt.Sprintf("added entity with ID %s to table", entityID))

	return entityID, nil
}

// GetExistingEntity checks if a given entity exists and returns the entity.id
func GetExistingEntity(repo *repository, entityEntityID string, entityType string) string {
	log.Info("entered function GetExistingEntity")

	var entityID = ""
	err := repo.db.Get(&entityID,
		`SELECT id FROM entities WHERE entity_id=$1 AND entity_type=$2`,
		entityEntityID, entityType)
	if err != nil {
		err = fmt.Errorf(`entity with entity_id: %s, entity_type: %s does not exist`,
			entityEntityID, entityType)
		log.Info(err.Error())
	}

	return entityID
}

// GetExistingAccount checks if a given account exists and returns the account ID
func GetExistingAccount(repo *repository, params *models.CreateTransaction) string {
	log.Info("entered function GetExistingAccount")

	// Check if entity exists
	entityID := GetExistingEntity(repo, *params.EntityID, *params.EntityType)
	if entityID == "" {
		err := errors.New(fmt.Sprintf("could not get entity.id for entity.entity_id: %s entity.entity_type: %s", *params.EntityID, *params.EntityType))
		log.Error(err.Error(), err)
		return ""
	}

	var accountID = ""
	err := repo.db.Get(&accountID,
		`SELECT 
			id FROM accounts 
		WHERE entity_id=$1 
		AND external_source_type=$2
		AND external_account_id=$3`,
		entityID, params.ExternalSourceType, params.ExternalAccountID)
	if err != nil {
		err = fmt.Errorf(`account with entity_id: %s, external_source_type: %s, external_account_id: %s does not exist`,
			entityID, *params.ExternalSourceType, *params.ExternalAccountID)
		log.Info(err.Error())
	}

	return accountID
}

// HandleEntity checks if a given entity exists and returns the entity.id
// If it does not exist it creates the entity and returns the new entity.id
func HandleEntity(repo *repository, params *models.CreateTransaction) (string, error) {

	// Check if entity exists
	entityID := GetExistingEntity(repo, *params.EntityID, *params.EntityType)
	if entityID != "" {
		log.Info(fmt.Sprintf("entity with entity_id: %s and entity_type: %s exists", *params.EntityID, *params.EntityType))
		return entityID, nil
	}

	// Else, create a new one.
	entityID, err := CreateEntity(repo, params)
	if err != nil {
		log.Error(log.Trace(), err)
		return "", err
	}

	return entityID, nil
}

// HandleAccount checks if a given account exists and returns the account ID
// If it does not exist it creates the account and returns the new account ID
func HandleAccount(repo *repository, params *models.CreateTransaction) (string, error) {

	// Check if account exists
	accountID := GetExistingAccount(repo, params)
	if accountID != "" {
		return accountID, nil
	}

	// Else, create a new one.
	accountID, err := CreateAccount(repo, params)
	if err != nil {
		log.Error(log.Trace(), err)
		return "", err
	}

	return accountID, nil
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
			l.description,
			l.metadata AS Metadata,
			l.created_at AS CreatedAt
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
			t.external_transaction_created_at AS ExternalTransactionCreatedAt,
			t.asset AS Asset,
			t.metadata AS Metadata,
			t.running_balance AS RunningBalance,
			t.transaction_category AS TransactionCategory,
			t.created_at AS CreatedAt
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
			&transaction.ExternalTransactionCreatedAt,
			&transaction.Asset,
			&transaction.Metadata,
			&transaction.RunningBalance,
			&transaction.TransactionCategory,
			&transaction.CreatedAt,
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
			t.external_transaction_created_at AS ExternalTransactionCreatedAt,
			t.asset AS Asset,
			t.metadata AS Metadata,
			t.running_balance AS RunningBalance,
			t.transaction_category AS TransactionCategory,
			t.created_at AS CreatedAt
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

// getRunningBalance is a function to get the balance of all accounts associated with
// the entity for the specified source_type (e.g. bill.com)
func getRunningBalance(repo *repository, params *models.CreateTransaction) (int64, error) {
	log.Info("entered function getRunningBalance")

	// Get the entity.id from the entity.entity_id
	entityID := GetExistingEntity(repo, *params.EntityID, *params.EntityType)
	if entityID == "" {
		err := errors.New(fmt.Sprintf("could not get entity.id for entity.entity_id: %s entity.entity_type: %s", *params.EntityID, *params.EntityType))
		log.Error(err.Error(), err)
		return 0, err
	}

	query := `
		SELECT
			t.id,
			t.running_balance AS RunningBalance
		FROM transactions t
		JOIN accounts a on t.account_id = a.id
		WHERE
			a.entity_id = $1 AND a.external_source_type = $2 AND a.external_account_id = $3
		order by t.created_at desc
		limit 1;`

	log.Info(log.StripSpecialChars(query))

	row := repo.db.QueryRowx(query, entityID, params.ExternalSourceType, params.ExternalAccountID)

	balance := RunningBalance{}
	if err := row.Scan(&balance.TransactionID, &balance.CurrentRunningBalance); err != nil {
		log.Error(err.Error(), err)

		// If it's an empty result
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	query = `
	SELECT 
		sum(case when amount < 0 then amount else 0 end)*-1 as TotalDebit,
		sum(case when amount >= 0 then amount else 0 end) as TotalCredit
	FROM line_items l
	WHERE l.transaction_id = $1;`

	log.Info(log.StripSpecialChars(query))

	row = repo.db.QueryRowx(query, balance.TransactionID)
	if err := row.Scan(&balance.TotalDebit, &balance.TotalCredit); err != nil {
		log.Error(err.Error(), err)
		return 0, err
	}

	newRunningBalance := balance.CurrentRunningBalance + (balance.TotalCredit - balance.TotalDebit)

	return newRunningBalance, nil
}

// CreateTransaction creates a new transaction and any related rows
// in required tables if they don't already exist
func (repo *repository) CreateTransaction(ctx context.Context, params *models.CreateTransaction) (*models.Transaction, error) {
	log.Info("entered function CreateTransaction")

	// Check if entity exists
	_, err := HandleEntity(repo, params)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	accountID, err := HandleAccount(repo, params)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	// Stub, replace.
	runningBalanceValue, err := getRunningBalance(repo, params)
	if err != nil {
		log.Fatal(err)
	}

	metaDataJSON := types.JSONText(params.Metadata)
	metaDataJSONValue, err := metaDataJSON.Value()
	if err != nil {
		log.Fatal(err)
	}

	// Set asset to usd default,
	// allow for optional asset provided via param
	// asset := ""
	// if params.Asset != "" {
	// 	asset = params.Asset
	// }

	// Create a new transaction entry
	sql := `
		INSERT INTO transactions (
			transaction_category,
			external_transaction_id,
			asset,
			account_id,
			running_balance,
			metadata
		)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING 
			id AS ID,
			account_id AS AccountID,
			external_transaction_id AS ExternalTransactionID,
			external_transaction_created_at AS ExternalTransactionCreatedAt,
			asset AS Asset,
			metadata AS Metadata,
			running_balance AS RunningBalance,
			transaction_category AS TransactionCategory,
			created_at AS CreatedAt`

	log.Info(fmt.Sprintf(log.StripSpecialChars(sql),
		&params.TransactionCategory,
		params.ExternalTransactionID,
		params.Asset,
		accountID,
		runningBalanceValue,
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
		&params.TransactionCategory,
		params.ExternalTransactionID,
		params.Asset,
		accountID,
		runningBalanceValue,
		metaDataJSONValue,
	)

	transaction := models.Transaction{}
	if err = row.Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.ExternalTransactionID,
		&transaction.ExternalTransactionCreatedAt,
		&transaction.Asset,
		&transaction.Metadata,
		&transaction.RunningBalance,
		&transaction.TransactionCategory,
		&transaction.CreatedAt,
	); err != nil {
		log.Error(err.Error(), err)
		return nil, err
	}

	items := []*models.LineItem{}
	for _, item := range params.LineItems {

		metaDataJSON := types.JSONText(item.Metadata)
		metaDataJSONValue, err = metaDataJSON.Value()
		if err != nil {
			log.Fatal(err)
		}

		// Inset Line Items
		sql = `
		INSERT INTO line_items (
			transaction_id,
			amount,
			description,
			metadata
		)
		VALUES ($1, $2, $3, $4)
		RETURNING 
			id AS ID,
			transaction_id AS TransactionID,
			amount AS Amount,
			description as Description,
			metadata AS Metadata,
			created_at AS CreatedAt`

		log.Info(fmt.Sprintf(log.StripSpecialChars(sql),
			transaction.ID,
			*item.Amount,
			*item.Description,
			metaDataJSONValue,
		))

		// Insert LineItem
		row := tx.QueryRowx(sql,
			transaction.ID,
			*item.Amount,
			*item.Description,
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

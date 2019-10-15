package balance

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/communitybridge/ledger/gen/restapi/operations/balance"
	"github.com/communitybridge/ledger/swagger"
	"github.com/jmoiron/sqlx"

	log "github.com/communitybridge/ledger/logging"
)

// Repository interface for repo calls
type Repository interface {
	GetEntityBalance(ctx context.Context, params *balance.GetBalanceParams) (*models.Balance, error)
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

// DoesEntityExist checks if a given entity exists in the entities table
func DoesEntityExist(repo *repository, id string) (bool, string) {
	log.Info("entered function DoesProjectExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM entities WHERE entity_id=$1", id)
	if err != nil {
		err = fmt.Errorf("project with id : `%s` does not exist", id)
		log.Info(err.Error())
		return false, res
	}

	return true, res
}

// GetBalance is function to get a balance for a specified Project
func (repo *repository) GetEntityBalance(ctx context.Context, params *balance.GetBalanceParams) (*models.Balance, error) {
	log.Info("entered function GetBalance")

	exist, _ := DoesEntityExist(repo, params.EntityID)
	if !exist {
		return nil, swagger.ErrNotFound
	}

	currentTime := time.Now().Unix()

	endDate := int64(0)
	if params.EndDate != nil {
		endDate = *params.EndDate
	}

	if endDate == int64(0) {
		endDate = currentTime
	}

	startDate := int64(0)
	if params.StartDate != nil {
		startDate = *params.StartDate
	}

	query := `
		SELECT
			e.entity_id AS EntityID,
			e.entity_type AS EntityType,
			count(t.id) as TotalCount,
			sum(l.amount) as TotalAmount
		FROM entities e
		LEFT JOIN transactions AS t on t.account_id = e.account_id
		LEFT JOIN line_items AS l on l.transaction_id = t.id
		WHERE
			e.entity_id = $1 AND t.created_at >= $2 AND t.created_at <= $3
		GROUP BY e.entity_id, e.entity_type;`

	log.Info(log.StripSpecialChars(query))

	row := repo.db.QueryRowx(query, params.EntityID, startDate, endDate)

	balance := &models.Balance{}
	if err := row.Scan(
		&balance.EntityID,
		&balance.EntityType,
		&balance.TotalCount,
		&balance.TotalBalance); err != nil {

		log.Error(err.Error(), err)

		// If it's an empty result
		if err == sql.ErrNoRows {
			return nil, swagger.ErrNotFound
		}

		return nil, swagger.ErrNotFound
	}

	return balance, nil
}

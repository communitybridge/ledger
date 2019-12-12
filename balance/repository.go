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
	log.Info("entered function DoesEntityExist")

	var res = ""
	err := repo.db.Get(&res, "SELECT id FROM entities WHERE entity_id=$1", id)
	if err != nil {
		err = fmt.Errorf("entity with id : `%s` does not exist", id)
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
			sum(DebitCount) AS SumDebitCount,
			sum(TotalDebit) AS SumTotalDebit,
			sum(CreditCount) AS SumCreditCount,
			sum(TotalCredit) AS SumTotalCredit
		FROM entities e
		LEFT JOIN accounts on accounts.entity_id = e.id
		LEFT JOIN transactions on transactions.account_id = accounts.id
		LEFT JOIN (
			SELECT transaction_id,
			sum(case when amount < 0 then 1 else 0 end) as DebitCount, 
			sum(case when amount < 0 then amount else 0 end)*-1 as TotalDebit,
			sum(case when amount >= 0 then 1 else 0 end) as CreditCount, 
			sum(case when amount >= 0 then amount else 0 end) as TotalCredit
			FROM line_items
			group by transaction_id) as l on l.transaction_id = transactions.id
		WHERE
			e.entity_id = $1 AND transactions.created_at >= $2 AND transactions.created_at <= $3
		GROUP BY EntityID, EntityType;`

	log.Info(log.StripSpecialChars(query))

	row := repo.db.QueryRowx(query, params.EntityID, startDate, endDate)

	balance := &models.Balance{}
	if err := row.Scan(
		&balance.EntityID,
		&balance.EntityType,
		&balance.DebitCount,
		&balance.TotalDebit,
		&balance.CreditCount,
		&balance.TotalCredit); err != nil {

		log.Error(err.Error(), err)

		// If it's an empty result
		if err == sql.ErrNoRows {
			return nil, swagger.ErrNotFound
		}

		return nil, swagger.ErrNotFound
	}

	balance.AvailableBalance = balance.TotalCredit - balance.TotalDebit
	balance.TotalCount = balance.CreditCount + balance.DebitCount

	fmt.Println(fmt.Sprintf("balance: %+v", balance))

	return balance, nil
}

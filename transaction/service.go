package transaction

import (
	"context"

	log "github.com/communitybridge/ledger/logging"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/communitybridge/ledger/gen/restapi/operations/transactions"
)

// Service ...
type Service interface {
	ListTransactions(ctx context.Context, params *transactions.ListTransactionsParams) (*models.TransactionList, error)
	CreateTransaction(ctx context.Context, params *transactions.CreateTransactionParams) (*models.Transaction, error)
	GetTransaction(ctx context.Context, params *transactions.GetTransactionParams) (*models.Transaction, error)
}

type service struct {
	repo Repository
}

// NewService create a service instance
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// ListTransactions calls on repository to get a list of transactions
func (s *service) ListTransactions(ctx context.Context, params *transactions.ListTransactionsParams) (*models.TransactionList, error) {
	log.Info("entered service ListTransactions")

	transactions, err := s.repo.ListTransactions(ctx, params)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	transactionList := models.TransactionList{}
	transactionList.TotalSize = int64(len(transactions))
	transactionList.PageSize = int64(len(transactions))

	transactionList.Data = transactions

	return &transactionList, nil
}

// GetTransaction returns the expected transaction response.
func (s *service) GetTransaction(ctx context.Context, params *transactions.GetTransactionParams) (*models.Transaction, error) {
	log.Info("entered service GetTransaction")

	transaction, err := s.repo.GetTransaction(ctx, params.TransactionID)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	return transaction, err
}

// CreateTransaction calls on repository to create a new transaction
// and then to return the transaction response expected.
func (s *service) CreateTransaction(ctx context.Context, params *transactions.CreateTransactionParams) (*models.Transaction, error) {
	log.Info("entered service CreateTransaction")

	transaction, err := s.repo.CreateTransaction(ctx, params.Transaction)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	transactionDetails, err := s.repo.GetTransaction(ctx, transaction.ID)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	return transactionDetails, err

}

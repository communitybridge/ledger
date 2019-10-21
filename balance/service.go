package balance

import (
	"context"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/communitybridge/ledger/gen/restapi/operations/balance"
	log "github.com/communitybridge/ledger/logging"
)

// Service ...
type Service interface {
	GetEntityBalance(ctx context.Context, params *balance.GetBalanceParams) (*models.Balance, error)
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

// GetEntityBalance returns the expected balance response for the specified Project.
func (s *service) GetEntityBalance(ctx context.Context, params *balance.GetBalanceParams) (*models.Balance, error) {
	log.Info("entered service GetBalance")

	balance, err := s.repo.GetEntityBalance(ctx, params)
	if err != nil {
		log.Error(log.Trace(), err)
		return nil, err
	}

	return balance, err
}

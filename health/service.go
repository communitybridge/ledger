package health

import (
	"context"
	"fmt"
	"time"

	log "github.com/communitybridge/ledger/logging"

	"github.com/communitybridge/ledger/gen/models"
)

// Service handles async log of audit event
type Service interface {
	GetHealth(ctx context.Context) (*models.Health, error)
}

type service struct {
}

// New returns new Service
func New() Service {
	return &service{}
}

func (s *service) GetHealth(ctx context.Context) (*models.Health, error) {
	log.Info("entered service GetHealth")

	t := time.Now()
	health := models.Health{
		DateTime: t.String(),
	}

	log.Debug(fmt.Sprintf("%#v", health))

	return &health, nil
}

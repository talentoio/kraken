package services

import (
	"context"
	"fmt"
	"kraken/internal/config"
	"kraken/internal/gateway"
	"kraken/internal/services/domain"
	"kraken/pkg/storage"

	"go.uber.org/zap"
)

type (
	ltpService struct {
		logger        *zap.Logger
		cfg           config.Config
		cache         storage.Cache
		krakenGateway gateway.KrakenGateway
	}

	LTPService interface {
		LTP(ctx context.Context) ([]*domain.LTPPair, error)
	}
)

//go:generate mockgen -source=ltp.go -destination=../../tests/mocks/ltp_service.go -package=mocks
func NewLTPService(
	cfg config.Config,
	logger *zap.Logger,
	cache storage.Cache,
	krakenGateway gateway.KrakenGateway,
) LTPService {
	return &ltpService{
		cfg:           cfg,
		logger:        logger,
		cache:         cache,
		krakenGateway: krakenGateway,
	}
}

// LTP fetches all available pairs data and return an  array of domain.LTPPair
func (s *ltpService) LTP(ctx context.Context) ([]*domain.LTPPair, error) {

	var ltpList = make([]*domain.LTPPair, 0, len(s.cfg.AvailablePairs))

	select {
	case <-ctx.Done():
	default:
		for _, pair := range s.cfg.AvailablePairs {
			value, exists := s.cache.Get(pair)
			if !exists {
				return nil, fmt.Errorf("%s value missing in cache", pair)
			}
			ltpList = append(ltpList, &domain.LTPPair{
				Pair:   pair,
				Amount: value,
			})
		}
	}

	return ltpList, nil
}

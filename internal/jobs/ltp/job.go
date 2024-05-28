package ltp

import (
	"context"
	"kraken/internal/config"
	"kraken/internal/gateway"
	"kraken/internal/jobs"
	"strconv"
	"time"

	"kraken/pkg/storage"

	"go.uber.org/zap"
)

type (
	ltpUpdaterJob struct {
		cfg           config.Config
		logger        *zap.Logger
		cache         storage.Cache
		krakenGateway gateway.KrakenGateway
	}
)

func NewLTPUpdaterJob(
	cfg config.Config,
	logger *zap.Logger,
	cache storage.Cache,
	krakenGateway gateway.KrakenGateway,
) jobs.Job {
	return &ltpUpdaterJob{
		cfg:           cfg,
		logger:        logger,
		cache:         cache,
		krakenGateway: krakenGateway,
	}
}

var (
	krakenPairKeys = map[string]string{
		"BTCEUR": "XXBTZEUR",
		"BTCUSD": "XXBTZUSD",
		"BTCCHF": "XBTCHF",
	}
)

func (j *ltpUpdaterJob) Run(ctx context.Context) {

	ticker := time.NewTicker(j.cfg.Kraken.UpdateDuration)

	j.logger.Info("starting ltpUpdaterJob scheduler....")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				j.logger.Info("updating ltp data....")
				for _, pair := range j.cfg.AvailablePairs {
					func(pair string) {

						requestCtx, cnc := context.WithTimeout(ctx, j.cfg.Kraken.RequestTimeOut)
						defer cnc()

						data, err := j.krakenGateway.LTP(requestCtx, pair)
						if err != nil {
							j.logger.Error("kraken.updating pair", zap.String("pair", pair), zap.Error(err))
							return
						}

						if len(data.Error) != 0 {
							j.logger.Error("kraken.updating pair.external error",
								zap.String("pair", pair),
								zap.Any("errors", data.Error))
							return
						}

						pairKey, ok := krakenPairKeys[pair]
						if !ok {
							j.logger.Error("kraken.updating pair.received invalid pair data")
							return
						}

						if data.Result == nil ||
							data.Result[pairKey].Closed == nil ||
							len(data.Result[pairKey].Closed) == 0 {
							j.logger.Error("kraken.updating pair.received wrong data")
							return
						}

						value := data.Result[pairKey].Closed[0]
						if _, err := strconv.ParseFloat(value, 64); err != nil {
							j.logger.Error("kraken.updating pair.received invalid float data")
							return
						}

						j.cache.Set(pair, value, j.cfg.Cache.Expire)

					}(pair)
				}
			}
		}
	}()

}

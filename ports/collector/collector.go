package collector

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/kriuchkov/tonbeacon/core/model"
	"github.com/kriuchkov/tonbeacon/core/ports"
)

const defaultCollectInterval = 10 * time.Millisecond

type Options struct {
	WalletPort      ports.WalletPort          `required:"true"`
	RepositoryPort  ports.AccountDatabasePort `required:"true"`
	CollectInterval time.Duration
}

func (o *Options) SetDefaults() {
	if o.CollectInterval == 0 {
		o.CollectInterval = defaultCollectInterval
	}
}

type collectorService struct {
	walletPort ports.WalletPort
	dbPort     ports.AccountDatabasePort
	interval   time.Duration
}

func NewCollectorService(opts *Options) ports.CollectorServicePort {
	if err := validator.New().Struct(opts); err != nil {
		log.Panic().Err(err).Msg("invalid options")
	}

	opts.SetDefaults()

	s := &collectorService{
		walletPort: opts.WalletPort,
		dbPort:     opts.RepositoryPort,
	}
	return s
}

func (s *collectorService) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := s.CollectFunds(context.Background()); err != nil {
				log.Warn().Err(err).Msg("collect funds")
			}

			<-ticker.C
		}
	}
}

func (s *collectorService) CollectFunds(ctx context.Context) error {
	_, err := s.dbPort.ListAccounts(ctx, model.ListAccountFilter{})
	if err != nil {
		return err
	}

	return nil
}

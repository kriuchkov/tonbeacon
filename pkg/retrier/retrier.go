//nolint:gomnd,gochecknoglobals // default values
package retrier

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type RetryPolicy struct {
	MaxAttempts        int
	StartDelay         time.Duration
	MaxDelay           *time.Duration
	BackoffCoefficient float32
}

type Retrier struct {
	policy         RetryPolicy
	excludedErrors []error
}

type Options func(r *Retrier)

func WithRetryPolicy(rp RetryPolicy) Options {
	return func(r *Retrier) {
		r.policy = rp
	}
}

func WithExcludedErrors(errors ...error) Options {
	return func(r *Retrier) {
		r.excludedErrors = errors
	}
}

var defaultPolicy = RetryPolicy{
	MaxAttempts:        10,
	StartDelay:         100 * time.Millisecond,
	MaxDelay:           lo.ToPtr(10 * time.Second),
	BackoffCoefficient: 2,
}

func NewRetrier(opts ...Options) *Retrier {
	retrier := &Retrier{policy: defaultPolicy}

	for _, opt := range opts {
		opt(retrier)
	}
	return retrier
}

func (r *Retrier) Wrap(ctx context.Context, name string, f func() error) (err error) {
	logger := log.Ctx(ctx).With().Str("name", name).Logger()

	delay := r.policy.StartDelay
	for i := 1; i <= r.policy.MaxAttempts; i++ {
		logger.Debug().Int("attempt", i).Dur("delay", delay).Msg("execution started")

		if err = f(); err == nil || r.checkExcludedErrors(err) {
			break
		}
		logger.Warn().Err(err).Msg("execution failed")

		if i != r.policy.MaxAttempts {
			time.Sleep(delay)
			delay = time.Duration(float32(delay) * r.policy.BackoffCoefficient)
			if r.policy.MaxDelay != nil && delay > *r.policy.MaxDelay {
				delay = *r.policy.MaxDelay
			}
		}
	}
	if err == nil {
		logger.Debug().Msg("execution succeeded")
	}
	return err
}

func (r *Retrier) checkExcludedErrors(err error) bool {
	_, ok := lo.Find(r.excludedErrors, func(item error) bool {
		return errors.Is(err, item)
	})
	return ok
}

package comfig

import (
	"sync"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type Sentrier interface {
	Sentry() *raven.Client
}

type sentrier struct {
	getter kv.Getter
	once   sync.Once
	value  *raven.Client
	err    error
}

func NewSentrier(getter kv.Getter) Sentrier {
	return &sentrier{getter: getter}
}

func (s *sentrier) Sentry() *raven.Client {
	s.once.Do(func() {
		var config struct {
			DSN string `fig:"dsn,required"`
		}

		err := figure.
			Out(&config).
			From(kv.MustGetStringMap(s.getter, "sentry")).
			Please()
		if err != nil {
			s.err = errors.Wrap(err, "failed to figure out")
			return
		}

		client, err := raven.New(config.DSN)
		if err != nil {
			s.err = errors.Wrap(err, "failed to init sentry client")
			return
		}

		// TODO tags

		s.value = client
	})
	if s.err != nil {
		panic(s.err)
	}
	return s.value
}

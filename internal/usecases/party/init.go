package party

import (
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/election/internal/domain/repository"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/txmanager"
)

type Module struct {
	partyRepo repository.PartyRepository
	txMgr     txmanager.TxManager
	publisher event.MessagePublisher
	topics    config.KafkaConfig
}

type Opts struct {
	PartyRepo repository.PartyRepository
	TxMgr     txmanager.TxManager
	Publisher event.MessagePublisher
	Topics    config.KafkaConfig
}

func New(opts *Opts) usecases.PartyUseCases {
	return &Module{
		partyRepo: opts.PartyRepo,
		txMgr:     opts.TxMgr,
		publisher: opts.Publisher,
		topics:    opts.Topics,
	}
}

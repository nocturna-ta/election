package election

import (
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/election/internal/domain/repository"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/txmanager"
)

type Module struct {
	electionRepo      repository.ElectionRepository
	supportingPartyUC usecases.SupportingPartyUseCases
	txMgr             txmanager.TxManager
	publisher         event.MessagePublisher
	topics            config.KafkaTopics
}

type Opts struct {
	ElectionRepo      repository.ElectionRepository
	SupportingPartyUC usecases.SupportingPartyUseCases
	TxMgr             txmanager.TxManager
	Publisher         event.MessagePublisher
	Topics            config.KafkaTopics
}

func New(opts *Opts) usecases.ElectionUseCases {
	return &Module{
		electionRepo:      opts.ElectionRepo,
		supportingPartyUC: opts.SupportingPartyUC,
		txMgr:             opts.TxMgr,
		publisher:         opts.Publisher,
		topics:            opts.Topics,
	}
}

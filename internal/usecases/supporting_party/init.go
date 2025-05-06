package supporting_party

import (
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/election/internal/domain/repository"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/txmanager"
)

type Module struct {
	supportingPartyRepo repository.SupportingPartyRepository
	partyRepo           repository.PartyRepository
	electionRepo        repository.ElectionRepository
	txMgr               txmanager.TxManager
	publisher           event.MessagePublisher
	topics              config.KafkaTopics
}

type Opts struct {
	SupportingPartyRepo repository.SupportingPartyRepository
	PartyRepo           repository.PartyRepository
	ElectionRepo        repository.ElectionRepository
	TxMgr               txmanager.TxManager
	Publisher           event.MessagePublisher
	Topics              config.KafkaTopics
}

func New(opts *Opts) usecases.SupportingPartyUseCases {
	return &Module{
		supportingPartyRepo: opts.SupportingPartyRepo,
		partyRepo:           opts.PartyRepo,
		electionRepo:        opts.ElectionRepo,
		txMgr:               opts.TxMgr,
		publisher:           opts.Publisher,
		topics:              opts.Topics,
	}
}

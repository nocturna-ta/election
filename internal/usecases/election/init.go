package election

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/election/internal/domain/repository"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/txmanager"
	"github.com/nocturna-ta/votechain-contract/binding/electionManager"
	"github.com/nocturna-ta/votechain-contract/interfaces"
)

type Module struct {
	electionRepo      repository.ElectionRepository
	supportingPartyUC usecases.SupportingPartyUseCases
	txMgr             txmanager.TxManager
	publisher         event.MessagePublisher
	topics            config.KafkaTopics
	electionContract  interfaces.ElectionManagerInterface
	client            ethereum.Client
}

type Opts struct {
	ElectionRepo      repository.ElectionRepository
	SupportingPartyUC usecases.SupportingPartyUseCases
	TxMgr             txmanager.TxManager
	Publisher         event.MessagePublisher
	Topics            config.KafkaTopics
	ElectionContract  interfaces.ElectionManagerInterface
	Client            ethereum.Client
	ContractAddress   common.Address
}

func New(opts *Opts) usecases.ElectionUseCases {
	var contractInterface interfaces.ElectionManagerInterface
	contract, err := electionManager.NewElectionManager(opts.ContractAddress, opts.Client.GetEthClient())
	if err != nil {
		return nil
	}
	contractInterface = contract
	return &Module{
		electionRepo:      opts.ElectionRepo,
		supportingPartyUC: opts.SupportingPartyUC,
		txMgr:             opts.TxMgr,
		publisher:         opts.Publisher,
		topics:            opts.Topics,
		electionContract:  contractInterface,
		client:            opts.Client,
	}
}

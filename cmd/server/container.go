package server

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/election/internal/interfaces/dao"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/election/internal/usecases/election"
	"github.com/nocturna-ta/election/internal/usecases/party"
	"github.com/nocturna-ta/election/internal/usecases/supporting_party"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/txmanager"
	txSql "github.com/nocturna-ta/golib/txmanager/sql"
)

type container struct {
	Cfg               config.MainConfig
	ElectionUc        usecases.ElectionUseCases
	PartyUc           usecases.PartyUseCases
	SupportingPartyUc usecases.SupportingPartyUseCases
}

type options struct {
	Cfg       *config.MainConfig
	DB        *sql.Store
	Client    ethereum.Client
	Publisher event.MessagePublisher
}

func newContainer(opts *options) *container {
	electionRepo := dao.NewElectionRepository(&dao.OptsElectionRepository{
		DB:     opts.DB,
		Client: opts.Client,
	})

	partyRepo := dao.NewPartyRepository(&dao.OptsPartyRepository{
		DB: opts.DB,
	})

	supportingRepo := dao.NewSupportingPartyRepository(&dao.OptsSupportingPartyRepository{
		DB: opts.DB,
	})

	txMgr, err := txmanager.New(context.Background(), &txmanager.DriverConfig{
		Type: "sql",
		Config: txSql.Config{
			DB: opts.DB,
		},
	})
	if err != nil {
		log.Fatal("Failed to instantiate transaction manager ")
	}

	partyUc := party.New(&party.Opts{
		PartyRepo: partyRepo,
		TxMgr:     txMgr,
		Publisher: opts.Publisher,
		Topics:    opts.Cfg.Kafka.Topics,
	})

	supportingPartyUc := supporting_party.New(&supporting_party.Opts{
		SupportingPartyRepo: supportingRepo,
		PartyRepo:           partyRepo,
		ElectionRepo:        electionRepo,
		TxMgr:               txMgr,
		Publisher:           opts.Publisher,
		Topics:              opts.Cfg.Kafka.Topics,
	})

	electionUc := election.New(&election.Opts{
		ElectionRepo:      electionRepo,
		SupportingPartyUC: supportingPartyUc,
		TxMgr:             txMgr,
		Publisher:         opts.Publisher,
		Topics:            opts.Cfg.Kafka.Topics,
		ContractAddress:   common.HexToAddress(opts.Cfg.Blockchain.ElectionManagerAddress),
		Client:            opts.Client,
	})

	return &container{
		Cfg:               *opts.Cfg,
		ElectionUc:        electionUc,
		PartyUc:           partyUc,
		SupportingPartyUc: supportingPartyUc,
	}

}

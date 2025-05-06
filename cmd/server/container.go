package server

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/election/internal/interfaces/dao"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/election/internal/usecases/election"
	"github.com/nocturna-ta/election/internal/usecases/party"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/txmanager"
	txSql "github.com/nocturna-ta/golib/txmanager/sql"
)

type container struct {
	Cfg        config.MainConfig
	ElectionUc usecases.ElectionUseCases
	PartyUc    usecases.PartyUseCases
}

type options struct {
	Cfg       *config.MainConfig
	DB        *sql.Store
	Client    ethereum.Client
	Publisher event.Publisher
}

func newContainer(opts *options) *container {
	electionRepo := dao.NewElectionRepository(&dao.OptsElectionRepository{
		DB:              opts.DB,
		ContractAddress: common.HexToAddress(opts.Cfg.Blockchain.ElectionManagerAddress),
		Client:          opts.Client,
	})

	partyRepo := dao.NewPartyRepository(&dao.OptsPartyRepository{
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

	electionUc := election.New(&election.Opts{
		ElectionRepo: electionRepo,
		TxMgr:        txMgr,
	})

	partyUc := party.New(&party.Opts{
		PartyRepo: partyRepo,
		TxMgr:     txMgr,
	})

	return &container{
		Cfg:        *opts.Cfg,
		ElectionUc: electionUc,
		PartyUc:    partyUc,
	}

}

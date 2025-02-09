package election

import (
	"github.com/nocturna-ta/election/internal/domain/repository"
	"github.com/nocturna-ta/election/internal/usecases"
	"github.com/nocturna-ta/golib/txmanager"
)

type Module struct {
	electionRepo repository.ElectionRepository
	txMgr        txmanager.TxManager
}

type Opts struct {
	ElectionRepo repository.ElectionRepository
	TxMgr        txmanager.TxManager
}

func New(opts *Opts) usecases.ElectionUseCases {
	return &Module{
		electionRepo: opts.ElectionRepo,
		txMgr:        opts.TxMgr,
	}
}

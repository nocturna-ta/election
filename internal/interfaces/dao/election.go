package dao

import (
	"github.com/nocturna-ta/election/internal/domain/repository"
	"github.com/nocturna-ta/golib/database/sql"
)

type ElectionRepository struct {
	db *sql.Store
}
type OptsElectionRepository struct {
	DB *sql.Store
}

func NewElectionRepository(opts *OptsElectionRepository) repository.ElectionRepository {
	return &ElectionRepository{
		db: opts.DB,
	}
}

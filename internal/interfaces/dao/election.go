package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lib/pq"
	"github.com/nocturna-ta/election/internal/domain/model"
	"github.com/nocturna-ta/election/internal/domain/repository"
	utils2 "github.com/nocturna-ta/election/pkg/utils"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/golib/txmanager/utils"
	"github.com/nocturna-ta/votechain-contract/binding"
)

type ElectionRepository struct {
	db       *sql.Store
	contract *binding.Votechain
	client   *ethclient.Client
}

type OptsElectionRepository struct {
	DB              *sql.Store
	ContractAddress common.Address
	Client          *ethclient.Client
}

func NewElectionRepository(opts *OptsElectionRepository) repository.ElectionRepository {
	contract, err := binding.NewVotechain(opts.ContractAddress, opts.Client)
	if err != nil {
		return nil
	}

	return &ElectionRepository{
		db:       opts.DB,
		contract: contract,
		client:   opts.Client,
	}
}

const (
	insertCandidate = `INSERT INTO candidates (id, name, election_no, is_active, created_at, updated_at)`
	selectCandidate = `SELECT %s FROM candidates %s WHERE TRUE %s`
	updateCandidate = `UPDATE candidates SET %s = WHERE TRUE %s`
)

func (e *ElectionRepository) InsertCandidate(ctx context.Context, candidate *model.Candidate, signedTransaction string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.InsertCandidate")
	defer span.End()

	var (
		err error
	)

	sqlTrx := utils.GetSqlTx(ctx)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, insertCandidate, candidate.ID, candidate.NameCandidate, candidate.ElectionNo, candidate.IsActive, candidate.CreatedAt, candidate.UpdatedAt)
	} else {
		_, err = e.db.GetMaster().ExecContext(ctx, insertCandidate, candidate.ID, candidate.NameCandidate, candidate.ElectionNo, candidate.IsActive, candidate.CreatedAt, candidate.UpdatedAt)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				log.WithFields(log.Fields{
					"error":     err,
					"candidate": candidate,
				}).ErrorWithCtx(ctx, "[ElectionRepository.InsertCandidate] Duplicate candidate")
				return ErrDuplicate
			}
		}
		log.WithFields(log.Fields{
			"error":     err,
			"candidate": candidate,
		}).ErrorWithCtx(ctx, "[ElectionRepository.InsertCandidate] Failed to insert candidate")
		return err
	}

	tx, err := utils2.StringToTx(signedTransaction)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.InsertCandidate] Failed to convert string to transaction")
		return err
	}

	err = e.client.SendTransaction(ctx, tx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.InsertCandidate] Failed to send transaction")
		return err
	}

	return nil
}

func (e *ElectionRepository) GetAllCandidate(ctx context.Context) ([]model.Candidate, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.GetAllCandidate")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		candidatesModels []model.Candidate
		err              error
	)

	candidates, err := e.contract.GetAllCandidates(nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.GetAllCandidate] Failed to get all candidates")
		return nil, err
	}

	selectQuery := `elections.id, elections.name_candidate, elections.election_no elections.is_active, elections.created_at, elections.updated_at`
	whereQuery := ` AND elections.is_deleted = false`
	joinQuery := ``

	query := fmt.Sprintf(selectCandidate, selectQuery, joinQuery, whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &candidatesModels, query)
	} else {
		err = e.db.GetMaster().SelectContext(ctx, &candidatesModels, query)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.GetAllCandidate] Failed to get all candidates")
		return nil, err
	}

	var matcherCandidate []model.Candidate
	for _, candidate := range candidates {
		for _, candidateModel := range candidatesModels {
			if candidateModel.ID.String() == candidate.Id {
				candidateModel.VoteCount = int(candidate.VoteCount.Int64())
				matcherCandidate = append(matcherCandidate, candidateModel)
			}
		}
	}

	if len(matcherCandidate) == 0 {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.GetAllCandidate] No candidate found")
		return nil, err
	}

	return matcherCandidate, nil
}

func (e *ElectionRepository) GetCandidateByNo(ctx context.Context, no string) (*model.Candidate, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.GetCandidateByNo")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		candidateModel model.Candidate
		err            error
		args           []any
	)

	candidates, err := e.contract.GetCandidateByNo(nil, no)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.GetCandidateByNo] Failed to get candidate by no")
		return nil, err
	}

	selectQuery := `elections.id, elections.name_candidate, elections.election_no elections.is_active, elections.created_at, elections.updated_at`
	whereQuery := ` AND elections.election_no = $1 AND elections.is_deleted = false`
	joinQuery := ``

	args = append(args, no)

	query := fmt.Sprintf(selectCandidate, selectQuery, joinQuery, whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &candidateModel, query, args...)
	} else {
		err = e.db.GetMaster().GetContext(ctx, &candidateModel, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.GetCandidateByNo] Failed to get candidate by no")
		return nil, err
	}

	if candidates.CandidateNo != candidateModel.ElectionNo {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.GetCandidateByNo] No candidate found")
		return nil, ErrNoResult
	}

	return &candidateModel, nil
}

func (e *ElectionRepository) CandidateActivate(ctx context.Context, id string, signedTransaction string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.CandidateActivate")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		err  error
		args []any
	)

	candidate, err := e.contract.GetCandidate(nil, id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.CandidateActivate] Failed to get candidate by no")
		return err
	}

	if candidate.IsActive {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.CandidateActivate] Candidate already activated")
		return ErrNoUpdateHappened
	}

	setQuery := "is_active = true"
	whereQuery := " AND id = $1 AND is_deleted = false"
	args = append(args, id)

	query := fmt.Sprintf(updateCandidate, setQuery, whereQuery)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		_, err = e.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"args":  args,
		}).ErrorWithCtx(ctx, "[ElectionRepository.CandidateActivate] Failed to activate candidate")
		return err
	}
	tx, err := utils2.StringToTx(signedTransaction)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.CandidateActivate] Failed to convert string to transaction")
		return err
	}

	err = e.client.SendTransaction(ctx, tx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.CandidateActivate] Failed to send transaction")
		return err
	}
	return nil
}

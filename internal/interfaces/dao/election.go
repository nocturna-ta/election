package dao

import (
	"context"
	sql2 "database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nocturna-ta/election/internal/domain/model"
	"github.com/nocturna-ta/election/internal/domain/repository"
	utils2 "github.com/nocturna-ta/election/pkg/utils"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/golib/txmanager/utils"
	"github.com/nocturna-ta/votechain-contract/binding/electionManager"
	"github.com/nocturna-ta/votechain-contract/interfaces"
)

type ElectionRepository struct {
	db       *sql.Store
	contract interfaces.ElectionManagerInterface
	client   ethereum.Client
}

type OptsElectionRepository struct {
	DB              *sql.Store
	ContractAddress common.Address
	Contract        interfaces.ElectionManagerInterface
	Client          ethereum.Client
}

func NewElectionRepository(opts *OptsElectionRepository) repository.ElectionRepository {
	var contractInterface interfaces.ElectionManagerInterface
	contract, err := electionManager.NewElectionManager(opts.ContractAddress, opts.Client.GetEthClient())
	if err != nil {
		return nil
	}
	contractInterface = contract
	return &ElectionRepository{
		db:       opts.DB,
		contract: contractInterface,
		client:   opts.Client,
	}
}

const (
	insertElectionPair = `
		INSERT INTO election_pairs (
		id, election_no, vote_count, is_active, pair_photo_path, president_full_name,
		president_education_history, president_work_experience, president_gender, president_birth_place,
		president_birth_date, president_last_education, president_job, president_photo_path, vice_president_full_name,
		vice_president_education_history, vice_president_work_experience, vice_president_gender, vice_president_birth_place,
		vice_president_birth_date, vice_president_last_education, vice_president_job, vice_president_photo_path, created_at,
		updated_at, is_deleted
		) VALUES (
			$1, $2, $3, $4, 
			$5, $6, $7, $8, $9, 
			$10, $11, $12, $13, $14,
			$15, $16, $17,$18,$19,$20,$21,$22,$23,$24,$25,$26
		)
	`
	selectElectionPair = `SELECT %s FROM election_pairs %s WHERE TRUE %s`
	updateElectionPair = `UPDATE election_pairs SET %s = WHERE TRUE %s`

	insertPairDetail = `
		INSERT INTO election_pair_details (
			id, election_pair_id, vision, mission, work_program,
			created_at, updated_at, is_deleted
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`
	selectPairDetail = `SELECT %s FROM election_pair_details %s WHERE TRUE %s`
	updatePairDetail = `UPDATE election_pair_details SET %s = WHERE TRUE %s`
)

func (e *ElectionRepository) InsertElectionPair(ctx context.Context, pair *model.ElectionPair, signedTransaction string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.InsertElectionPair")
	defer span.End()

	presidentEducationJSON, err := json.Marshal(pair.President.EducationHistory)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to marshal president education history")
		return err
	}

	presidentWorkJSON, err := json.Marshal(pair.President.WorkExperience)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to marshal president work experience")
		return err
	}

	vicePresidentEducationJSON, err := json.Marshal(pair.VicePresident.EducationHistory)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to marshal vice president work experience")
		return err
	}

	vicePresidentWorkJSON, err := json.Marshal(pair.VicePresident.WorkExperience)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to marshal vice president work experience")
		return err
	}

	tx, err := utils2.StringToTx(signedTransaction)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to convert signed transaction")
		return err
	}

	sqlTrx := utils.GetSqlTx(ctx)

	var ownTransaction bool
	if sqlTrx == nil {
		var err error
		sqlTrx, err = e.db.GetMaster().BeginTxx(ctx, nil)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to begin transaction")
			return err
		}
		ownTransaction = true

		defer func() {
			if err != nil && ownTransaction {
				rollbackErr := sqlTrx.Rollback()
				if rollbackErr != nil {
					log.WithFields(log.Fields{
						"error": rollbackErr,
					}).ErrorWithCtx(ctx, "[ElectionRepository.InsertElectionPair] Failed to rollback transaction")
				}
			}
		}()
	}
	_, err = sqlTrx.ExecContext(
		ctx,
		insertElectionPair,
		pair.ID,
		pair.ElectionNo,
		pair.VoteCount,
		pair.IsActive,
		pair.PairPhotoPath,
		pair.President.FullName,
		presidentEducationJSON,
		presidentWorkJSON,
		pair.President.Gender,
		pair.President.BirthPlace,
		pair.President.BirthDate,
		pair.President.LastEducation,
		pair.President.Job,
		pair.President.PhotoPath,
		pair.VicePresident.FullName,
		vicePresidentEducationJSON,
		vicePresidentWorkJSON,
		pair.VicePresident.Gender,
		pair.VicePresident.BirthPlace,
		pair.VicePresident.BirthDate,
		pair.VicePresident.LastEducation,
		pair.VicePresident.Job,
		pair.VicePresident.PhotoPath,
		pair.CreatedAt,
		pair.UpdatedAt,
		pair.IsDeleted,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				log.WithFields(log.Fields{
					"error": err,
					"pair":  pair,
				}).ErrorWithCtx(ctx, "[ElectionRepository.InsertElectionPair] Duplicate election pair")
				return ErrDuplicate
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"pair":  pair,
		}).ErrorWithCtx(ctx, "[ElectionRepository.InsertElectionPair] Failed to insert election pair")
		return err
	}

	err = e.client.SendTransaction(ctx, tx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.InsertElectionPair] Failed to send transaction")

		if ownTransaction {
			rollbackErr := sqlTrx.Rollback()
			if rollbackErr != nil {
				log.WithFields(log.Fields{
					"error": rollbackErr,
				}).ErrorWithCtx(ctx, "[ElectionRepository.InsertElectionPair] Failed to rollback transaction")
			}
		}
		return err
	}

	if ownTransaction {
		if err := sqlTrx.Commit(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionRepository.InsertElectionPair] Failed to commit transaction")
			return err
		}
	}

	return nil
}

func (e *ElectionRepository) GetElectionPairByID(ctx context.Context, id uuid.UUID) (*model.ElectionPair, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.GetElectionPairByID")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		electionPairModelDTO model.ElectionPairDTO
		err                  error
		args                 []any
	)

	electionPair, err := e.contract.GetElection(nil, id.String())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to get election pair from contract")
		return nil, err
	}

	selectQuery := `id, election_no, vote_count, is_active, pair_photo_path, 
		president_full_name, president_education_history, president_work_experience, 
		president_gender, president_birth_place, president_birth_date, president_last_education, 
		president_job, president_photo_path, vice_president_full_name, vice_president_education_history, 
		vice_president_work_experience, vice_president_gender, vice_president_birth_place, 
		vice_president_birth_date, vice_president_last_education, vice_president_job, 
		vice_president_photo_path, created_at, updated_at, is_deleted`
	whereClause := ` AND id = $1 AND is_deleted = false`
	joinQuery := ``
	args = append(args, id)

	query := fmt.Sprintf(selectElectionPair, selectQuery, joinQuery, whereClause)
	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &electionPairModelDTO, query, args...)
	} else {
		err = e.db.GetMaster().GetContext(ctx, &electionPairModelDTO, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to get election pair by id")
		return nil, err
	}
	if errors.Is(err, sql2.ErrNoRows) {
		return nil, ErrNoResult
	}

	if electionPair.Id != electionPairModelDTO.ID.String() && electionPair.ElectionNo != electionPairModelDTO.ElectionNo {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Election pair id not match")
		return nil, ErrNoResult
	}

	dtoToDomainElection, err := electionPairModelDTO.ToDomain()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to convert election pair model to domain")
		return nil, err
	}

	return dtoToDomainElection, err
}

func (e *ElectionRepository) GetElectionPairByNo(ctx context.Context, no string) (*model.ElectionPair, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.GetElectionPairByNo")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		electionPairModelDTO model.ElectionPairDTO
		err                  error
		args                 []any
	)

	electionPair, err := e.contract.GetElectionByNo(nil, no)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to get election pair from contract")
		return nil, err
	}

	selectQuery := `id, election_no, vote_count, is_active, pair_photo_path, 
		president_full_name, president_education_history, president_work_experience, 
		president_gender, president_birth_place, president_birth_date, president_last_education, 
		president_job, president_photo_path, vice_president_full_name, vice_president_education_history, 
		vice_president_work_experience, vice_president_gender, vice_president_birth_place, 
		vice_president_birth_date, vice_president_last_education, vice_president_job, 
		vice_president_photo_path, created_at, updated_at, is_deleted`
	whereClause := ` AND election_no = $1 AND is_deleted = false`
	joinQuery := ``
	args = append(args, no)

	query := fmt.Sprintf(selectElectionPair, selectQuery, joinQuery, whereClause)
	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &electionPairModelDTO, query, args...)
	} else {
		err = e.db.GetMaster().GetContext(ctx, &electionPairModelDTO, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"no":    no,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to get election pair by no")
		return nil, err
	}
	if errors.Is(err, sql2.ErrNoRows) {
		return nil, ErrNoResult
	}

	dtoToDomainElection, err := electionPairModelDTO.ToDomain()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to convert election pair model to domain")
		return nil, err
	}

	if electionPair.ElectionNo != dtoToDomainElection.ElectionNo {
		log.WithFields(log.Fields{
			"error": err,
			"no":    no,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Election pair no not match")
		return nil, ErrNoResult
	}

	return dtoToDomainElection, nil
}

func (e *ElectionRepository) GetAllElectionPairs(ctx context.Context) ([]model.ElectionPair, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.GetAllElectionPairs")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		electionPairModelsDTO []model.ElectionPairDTO
		err                   error
	)

	electionPairs, err := e.contract.GetAllElection(nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to get all election pairs from contract")
		return nil, err
	}

	selectQuery := `id, election_no, vote_count, is_active, pair_photo_path, 
		president_full_name, president_education_history, president_work_experience, 
		president_gender, president_birth_place, president_birth_date, president_last_education, 
		president_job, president_photo_path, vice_president_full_name, vice_president_education_history, 
		vice_president_work_experience, vice_president_gender, vice_president_birth_place, 
		vice_president_birth_date, vice_president_last_education, vice_president_job, 
		vice_president_photo_path, created_at, updated_at, is_deleted`

	whereClause := ` AND is_deleted = false ORDER BY election_no`
	joinQuery := ``

	query := fmt.Sprintf(selectElectionPair, selectQuery, joinQuery, whereClause)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &electionPairModelsDTO, query)
	} else {
		err = e.db.GetMaster().SelectContext(ctx, &electionPairModelsDTO, query)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to get all election pairs")
		return nil, err
	}

	pairs := make([]model.ElectionPair, len(electionPairModelsDTO))
	for i, dto := range electionPairModelsDTO {
		pairModel, err := dto.ToDomain()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to convert election pair model to domain")
			return nil, err
		}
		pairs[i] = *pairModel
	}

	var matchedPairs []model.ElectionPair
	for _, electionPair := range electionPairs {
		for _, pairModel := range pairs {
			if electionPair.Id == pairModel.ID.String() {
				matchedPairs = append(matchedPairs, pairModel)
				break
			}
		}
	}

	if len(matchedPairs) == 0 {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] No election pairs found")
		return nil, ErrNoResult
	}

	return pairs, nil
}

func (e *ElectionRepository) ActivateElectionPair(ctx context.Context, id uuid.UUID, signedTransaction string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.ActivateElectionPair")
	defer span.End()

	tx, err := utils2.StringToTx(signedTransaction)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to convert signed transaction")
		return err
	}

	sqlTrx := utils.GetSqlTx(ctx)

	var ownTransaction bool
	if sqlTrx == nil {
		var err error
		sqlTrx, err = e.db.GetMaster().BeginTxx(ctx, nil)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to begin transaction")
			return err
		}
		ownTransaction = true

		defer func() {
			if err != nil && ownTransaction {
				rollbackErr := sqlTrx.Rollback()
				if rollbackErr != nil {
					log.WithFields(log.Fields{
						"error": rollbackErr,
					}).ErrorWithCtx(ctx, "[ElectionRepository.ActivateElectionPair] Failed to rollback transaction")
				}
			}
		}()
	}

	var args []any
	setQuery := `is_active = true`
	whereQuery := ` AND id = $1 AND is_deleted = false`
	args = append(args, id)

	query := fmt.Sprintf(updateElectionPair, setQuery, whereQuery)

	result, err := sqlTrx.ExecContext(ctx, query, args...)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to activate election pair")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionRepository] Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		return ErrNoUpdateHappened
	}

	err = e.client.SendTransaction(ctx, tx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.ActivateElectionPair] Failed to send transaction")

		if ownTransaction {
			rollbackErr := sqlTrx.Rollback()
			if rollbackErr != nil {
				log.WithFields(log.Fields{
					"error": rollbackErr,
				}).ErrorWithCtx(ctx, "[ElectionRepository.ActivateElectionPair] Failed to rollback transaction")
			}
		}
		return err
	}

	if ownTransaction {
		if err := sqlTrx.Commit(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionRepository.ActivateElectionPair] Failed to commit transaction")
			return err
		}
	}

	return nil
}

func (e *ElectionRepository) UpsertPairDetail(ctx context.Context, detail *model.PairDetail) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.UpsertPairDetail")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var ownTransaction bool
	if sqlTrx == nil {
		var err error
		sqlTrx, err = e.db.GetMaster().BeginTxx(ctx, nil)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionRepository.UpsertPairDetail] Failed to begin transaction")
			return err
		}
		ownTransaction = true

		defer func() {
			if err != nil && ownTransaction {
				rollbackErr := sqlTrx.Rollback()
				if rollbackErr != nil {
					log.WithFields(log.Fields{
						"error": rollbackErr,
					}).ErrorWithCtx(ctx, "[ElectionRepository.UpsertPairDetail] Failed to rollback transaction")
				}
			}
		}()
	}

	var (
		count      int
		err        error
		updateArgs []any
	)
	countQuery := `SELECT COUNT(*) FROM election_pair_details WHERE id = $1`

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &count, countQuery, detail.ID)
	} else {
		err = e.db.GetMaster().GetContext(ctx, &count, countQuery, detail.ID)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionRepository.UpsertPairDetail] Failed to get count")
		return err
	}

	setQuery := `vision = $1, mission = $2, work_program = $3, updated_at = $4`
	whereQuery := ` AND election_pair_id = $5 AND is_deleted = false`
	updateArgs = append(updateArgs, detail.Vision, detail.Mission, detail.WorkProgram, detail.UpdatedAt, detail.ElectionPairID)

	queryUpdate := fmt.Sprintf(updatePairDetail, setQuery, whereQuery)

	if count > 0 {
		_, err = sqlTrx.ExecContext(ctx, queryUpdate, updateArgs...)
	} else {
		_, err = sqlTrx.ExecContext(ctx, insertPairDetail,
			detail.ID,
			detail.ElectionPairID,
			detail.Vision,
			detail.Mission,
			detail.WorkProgram,
			detail.CreatedAt,
			detail.UpdatedAt,
			detail.IsDeleted,
		)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"detail": detail,
		}).ErrorWithCtx(ctx, "[ElectionRepository.UpsertPairDetail] Failed to upsert pair detail")
	}

	if ownTransaction {
		if err := sqlTrx.Commit(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionRepository.UpsertPairDetail] Failed to commit transaction")
			return err
		}
	}

	return nil
}

func (e *ElectionRepository) GetPairDetailByPairID(ctx context.Context, pairID uuid.UUID) (*model.PairDetail, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.GetPairDetailByPairID")
	defer span.End()

	var (
		detail model.PairDetail
		err    error
		args   []any
	)

	sqlTrx := utils.GetSqlTx(ctx)

	selectQuery := `id, election_pair_id, vision, mission, work_program, created_at, updated_at`
	whereClause := ` AND election_pair_id = $1 AND is_deleted = false`
	joinQuery := ``
	args = append(args, pairID)

	query := fmt.Sprintf(selectPairDetail, selectQuery, joinQuery, whereClause)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &detail, query, args...)
	} else {
		err = e.db.GetMaster().GetContext(ctx, &detail, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    pairID,
		}).ErrorWithCtx(ctx, "[ElectionRepository.GetPairDetailByPairID] Failed to get pair detail by pair id")

		if errors.Is(err, sql2.ErrNoRows) {
			return nil, ErrNoResult
		}

		return nil, err
	}

	var docsPath []string
	queryDocs := `SELECT document_path FROM election_pair_details WHERE election_pair_id = $1 AND is_deleted = false`
	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &docsPath, queryDocs, pairID)
	} else {
		err = e.db.GetMaster().SelectContext(ctx, &docsPath, queryDocs, pairID)
	}

	if err != nil && !errors.Is(err, sql2.ErrNoRows) {
		log.WithFields(log.Fields{
			"error":    err,
			"detailID": detail.ID,
		}).ErrorWithCtx(ctx, "[ElectionRepository.GetPairDetailByPairID] Failed to get program documents")
		return nil, err
	}

	detail.ProgramDocs = docsPath

	return &detail, nil

}

func (e *ElectionRepository) UpdateElectionPairPhoto(ctx context.Context, id uuid.UUID, photoPath string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionRepository.UpdateElectionPairPhoto")
	defer span.End()

	panic("implement me")
}

func (e *ElectionRepository) UpdatePresidentPhoto(ctx context.Context, id uuid.UUID, photoPath string) error {
	//TODO implement me
	panic("implement me")
}

func (e *ElectionRepository) UpdateVicePresidentPhoto(ctx context.Context, id uuid.UUID, photoPath string) error {
	//TODO implement me
	panic("implement me")
}

func (e *ElectionRepository) GetPresidentPhotoPath(ctx context.Context, id uuid.UUID) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (e *ElectionRepository) GetVicePresidentPhotoPath(ctx context.Context, id uuid.UUID) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (e *ElectionRepository) GetElectionPairPhotoPath(ctx context.Context, id uuid.UUID) (string, error) {
	//TODO implement me
	panic("implement me")
}

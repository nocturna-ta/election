package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nocturna-ta/election/internal/domain/model"
	"github.com/nocturna-ta/election/internal/domain/repository"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/golib/txmanager/utils"
)

type SupportingPartyRepository struct {
	db *sql.Store
}

type OptsSupportingPartyRepository struct {
	DB *sql.Store
}

func NewSupportingPartyRepository(opts *OptsSupportingPartyRepository) repository.SupportingPartyRepository {
	return &SupportingPartyRepository{
		db: opts.DB,
	}
}

const (
	insertSupportingPartyQuery = `INSERT INTO supporting_parties (id, election_pair_id, party_id, created_at, updated_at, is_deleted) 
									VALUES ($1, $2, $3, $4, $5, $6)`
	selectSupportingQuery = `SELECT %s FROM supporting_parties %s WHERE TRUE %s`
	deleteSupportingQuery = `UPDATE supporting_parties SET %s WHERE TRUE %s`
)

func (s *SupportingPartyRepository) AddSupportingParty(ctx context.Context, supportingParty *model.SupportingParty) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "SupportingPartyRepository.AddSupportingParty")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		err error
	)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, insertSupportingPartyQuery, supportingParty.ID,
			supportingParty.ElectionPairID,
			supportingParty.PartyID,
			supportingParty.CreatedAt,
			supportingParty.UpdatedAt,
			supportingParty.IsDeleted)
	} else {
		_, err = s.db.GetMaster().ExecContext(ctx, insertSupportingPartyQuery, supportingParty.ID, supportingParty.ElectionPairID,
			supportingParty.PartyID,
			supportingParty.CreatedAt,
			supportingParty.UpdatedAt,
			supportingParty.IsDeleted)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				log.WithFields(log.Fields{
					"error":           err,
					"supportingParty": supportingParty,
				}).ErrorWithCtx(ctx, "[SupportingParty.ddSupportingParty] duplicate key value violates unique constraint")
				return ErrDuplicate
			}
		}

		log.WithFields(log.Fields{
			"error":           err,
			"supportingParty": *supportingParty,
		}).ErrorWithCtx(ctx, "[SupportingPartyRepository.AddSupportingParty] failed to add supporting party")

		return err
	}

	return nil
}

func (s *SupportingPartyRepository) GetSupportingPartiesByPairID(ctx context.Context, pairID uuid.UUID) ([]model.SupportingParty, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "SupportingPartyRepository.GetSupportingPartiesByPairID")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		err                error
		supportingPartyDTO []model.SupportingPartyDTO
		args               []any
	)

	selectQuery := `
		sp.id, sp.election_pair_id, sp.party_id, sp.created_at, sp.updated_at, sp.is_deleted,
		p.id as party_id2, p.name as party_name, p.logo_path as party_logo_path,
		p.created_at as party_created_at, p.updated_at as party_updated_at, p.is_deleted as party_is_deleted
	`
	joinQuery := ` as sp JOIN parties as p ON sp.party_id = p.id AND p.is_deleted = false`
	whereClause := ` AND sp.election_pair_id = $1 AND sp.is_deleted = false`
	args = append(args, pairID)

	query := fmt.Sprintf(selectSupportingQuery, selectQuery, joinQuery, whereClause)

	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &supportingPartyDTO, query, args...)
	} else {
		err = s.db.GetMaster().SelectContext(ctx, &supportingPartyDTO, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"pairID": pairID,
		}).ErrorWithCtx(ctx, "[SupportingPartyRepository.GetSupportingPartiesByPairID] failed to get supporting parties by pair ID")
		return nil, err
	}

	supportingParties := make([]model.SupportingParty, len(supportingPartyDTO))
	for i, dto := range supportingPartyDTO {
		supportingParty := dto.ToDomain()
		supportingParties[i] = *supportingParty
	}

	return supportingParties, nil

}

func (s *SupportingPartyRepository) RemoveSupportingParty(ctx context.Context, pairID uuid.UUID, partyID uuid.UUID) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "SupportingPartyRepository.RemoveSupportingPart")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		err  error
		args []any
	)

	setQuery := `is_deleted = true`
	whereClause := ` AND election_pair_id = $1 AND party_id = $2`
	args = append(args, pairID, partyID)
	query := fmt.Sprintf(deleteSupportingQuery, setQuery, whereClause)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		_, err = s.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"pairID":  pairID,
			"partyID": partyID,
		}).ErrorWithCtx(ctx, "[SupportingPartyRepository.RemoveSupportingPart] failed to remove supporting party")
		return err
	}

	return nil
}

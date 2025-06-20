package dao

import (
	"context"
	sql2 "database/sql"
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

type PartyRepository struct {
	db *sql.Store
}

type OptsPartyRepository struct {
	DB *sql.Store
}

func NewPartyRepository(opts *OptsPartyRepository) repository.PartyRepository {
	return &PartyRepository{
		db: opts.DB,
	}
}

const (
	insertPartyQuery = `INSERT INTO parties (id, name, logo_path) VALUES ($1, $2, $3)`
	selectPartyQuery = `SELECT %s FROM parties %s WHERE TRUE %s`
	updatePartyQuery = `UPDATE parties SET %s WHERE TRUE %s`
)

func (p *PartyRepository) GetPartyByID(ctx context.Context, id uuid.UUID) (*model.Party, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyRepository.GetPartyByID")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		party model.Party
		err   error
		args  []any
	)

	selectQuery := `id, name, logo_path, created_at, updated_at, is_deleted`
	whereClause := ` AND id = $1 AND is_deleted = false`
	joinQuery := ``
	args = append(args, id)

	query := fmt.Sprintf(selectPartyQuery, selectQuery, joinQuery, whereClause)
	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &party, query, args...)
	} else {
		err = p.db.GetMaster().GetContext(ctx, &party, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).ErrorWithCtx(ctx, "failed to get party by id")
		return nil, err
	}

	if errors.Is(err, sql2.ErrNoRows) {
		return nil, ErrNoResult
	}

	return &party, nil
}

func (p *PartyRepository) GetAllParties(ctx context.Context) ([]model.Party, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyRepository.GetAllParties")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		parties []model.Party
		err     error
		args    []any
	)

	selectQuery := `id, name, logo_path, created_at, updated_at, is_deleted`
	whereClause := ` AND is_deleted = false`
	joinQuery := ``

	query := fmt.Sprintf(selectPartyQuery, selectQuery, joinQuery, whereClause)
	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &parties, query, args...)
	} else {
		err = p.db.GetMaster().SelectContext(ctx, &parties, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "failed to get all parties")
		return nil, err
	}

	return parties, nil
}

func (p *PartyRepository) InsertParty(ctx context.Context, party *model.Party) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyRepository.InsertParty")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err error
	)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, insertPartyQuery, party.ID, party.Name, party.LogoPath)

	} else {
		_, err = p.db.GetMaster().ExecContext(ctx, insertPartyQuery, party.ID, party.Name, party.LogoPath)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				log.WithFields(log.Fields{
					"error": err,
					"party": party,
				}).ErrorWithCtx(ctx, "[SupportingParty.ddSupportingParty] duplicate key value violates unique constraint")
				return ErrDuplicate
			}
		}

		log.WithFields(log.Fields{
			"party": *party,
			"error": err,
		}).ErrorWithCtx(ctx, "[PartyRepository.InsertParty] failed to insert party")
		return err
	}

	return nil

}

func (p *PartyRepository) UpdateParty(ctx context.Context, party *model.Party) (*model.Party, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyRepository.UpdateParty")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err  error
		args []any
	)

	updateQuery := `name = $1, logo_path = $2`
	whereClause := ` AND id = $3 AND is_deleted = false`
	args = append(args, party.Name, party.LogoPath, party.ID)
	updateQuery = fmt.Sprintf(updatePartyQuery, updateQuery, whereClause)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, updateQuery, args...)
	} else {
		_, err = p.db.GetMaster().ExecContext(ctx, updateQuery, args...)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"party": party,
			"error": err,
		}).ErrorWithCtx(ctx, "failed to update party")
		return nil, err
	}

	return party, nil
}

func (p *PartyRepository) DeleteParty(ctx context.Context, id uuid.UUID) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyRepository.DeleteParty")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		err  error
		args []any
	)

	updateQuery := `is_deleted = true`
	whereClause := ` AND id = $1 AND is_deleted = false`
	args = append(args, id)
	updateQuery = fmt.Sprintf(updatePartyQuery, updateQuery, whereClause)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, updateQuery, args...)
	} else {
		_, err = p.db.GetMaster().ExecContext(ctx, updateQuery, args...)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).ErrorWithCtx(ctx, "failed to delete party")
		return err
	}

	return nil
}

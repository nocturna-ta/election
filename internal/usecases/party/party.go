package party

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/domain/model"
	"github.com/nocturna-ta/election/internal/interfaces/dao"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
	"github.com/nocturna-ta/election/pkg/constants"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/fileutils"
	"github.com/nocturna-ta/golib/http"
	"github.com/nocturna-ta/golib/http/filehandler"
	"github.com/nocturna-ta/golib/log"
	response2 "github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/tracing"
)

func (m *Module) RegisterParty(ctx context.Context, req *request.PartyRegisterRequest) (*response.PartyResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyUseCases.RegisterParty")
	defer span.End()

	var (
		party *model.Party
	)

	transaction := func(txCtx context.Context) (any, error) {
		party = model.ConstructPartyRegistration(req)

		if req.LogoFile != nil {
			fileConfig := fileutils.DefaultConfig()
			fileConfig.SetAllowedImageExtension()
			fileConfig.EntityType = "party"

			logoPath, err := fileutils.StoreFile(txCtx, req.LogoFile, req.LogoName, fileConfig)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).ErrorWithCtx(txCtx, "[PartyUseCases.RegisterParty] failed to store logo file")
				return nil, err
			}

			party.LogoPath = logoPath
		}

		errTx := m.partyRepo.InsertParty(txCtx, party)
		if errTx != nil {
			_ = fileutils.DeleteFile(txCtx, party.LogoPath)
			if errors.Is(errTx, dao.ErrDuplicate) {
				return nil, &custerr.ErrChain{
					Message: "party already exists",
					Cause:   errTx,
					Code:    400,
					Type:    response2.ErrBadRequest,
				}
			}
			return nil, errTx
		}

		errTx = m.publisher.Publish(txCtx, m.topics.MasterDataParty.Value, party.ID.String(), party.ToMessageModel(), map[string]any{
			constants.MetaDataOperation: constants.Create,
		})

		if errTx != nil {
			return nil, errTx
		}
		return nil, nil
	}
	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return nil, err
	}

	return &response.PartyResponse{
		ID:       party.ID.String(),
		Name:     party.Name,
		LogoPath: party.LogoPath,
	}, nil
}

func (m *Module) GetPartyByID(ctx context.Context, id string) (*response.PartyResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyUseCases.GetPartyByID")
	defer span.End()

	partyID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	party, err := m.partyRepo.GetPartyByID(ctx, partyID)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    id,
			"error": err,
		}).ErrorWithCtx(ctx, "[PartyUseCases.GetPartyByID] failed to get party by id")
		return nil, err
	}

	return &response.PartyResponse{
		ID:       party.ID.String(),
		Name:     party.Name,
		LogoPath: party.LogoPath,
	}, nil
}

func (m *Module) GetAllParties(ctx context.Context) ([]response.PartyResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyUseCases.GetAllParties")
	defer span.End()

	parties, err := m.partyRepo.GetAllParties(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[PartyUseCases.GetAllParties] failed to get all parties")
		return nil, err
	}

	var partyResponses []response.PartyResponse
	for _, party := range parties {
		partyResponses = append(partyResponses, response.PartyResponse{
			ID:       party.ID.String(),
			Name:     party.Name,
			LogoPath: party.LogoPath,
		})
	}

	return partyResponses, nil
}

func (m *Module) UpdateParty(ctx context.Context, req *request.PartyUpdateRequest) (*response.PartyResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyUseCases.UpdateParty")
	defer span.End()

	partyID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	var (
		updatedParty *model.Party
	)

	transaction := func(txCtx context.Context) (any, error) {
		existing, errTx := m.partyRepo.GetPartyByID(txCtx, partyID)
		if errTx != nil {
			log.WithFields(log.Fields{
				"id":    req.ID,
				"error": errTx,
			}).ErrorWithCtx(txCtx, "[PartyUseCases.UpdateParty] failed to get party by id")
			return nil, errTx
		}

		oldLogoPath := existing.LogoPath
		existing.Name = req.Name

		if req.LogoFile != nil {
			fileConfig := fileutils.DefaultConfig()
			fileConfig.SetAllowedImageExtension()
			fileConfig.EntityType = "party"

			logoPath, errTx := fileutils.StoreFile(txCtx, req.LogoFile, req.LogoName, fileConfig)
			if errTx != nil {
				log.WithFields(log.Fields{
					"error": errTx,
				}).ErrorWithCtx(txCtx, "[PartyUseCases.UpdateParty] failed to store logo file")
				return nil, errTx
			}

			existing.LogoPath = logoPath
		}

		updatedParty, errTx = m.partyRepo.UpdateParty(txCtx, existing)
		if errTx != nil {
			if req.LogoFile != nil && existing.LogoPath != "" {
				_ = fileutils.DeleteFile(txCtx, existing.LogoPath)
			}
			log.WithFields(log.Fields{
				"id":    req.ID,
				"error": errTx,
			}).ErrorWithCtx(txCtx, "[PartyUseCases.UpdateParty] failed to update party")
			return nil, errTx
		}

		if req.LogoFile != nil && oldLogoPath != "" && oldLogoPath != existing.LogoPath {
			_ = fileutils.DeleteFile(txCtx, oldLogoPath)
		}

		errTx = m.publisher.Publish(txCtx, m.topics.MasterDataParty.Value, updatedParty.ID.String(), updatedParty.ToMessageModel(), map[string]any{
			constants.MetaDataOperation: constants.Update,
		})

		return nil, nil
	}

	_, err = m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return nil, err
	}

	return &response.PartyResponse{
		ID:       updatedParty.ID.String(),
		Name:     updatedParty.Name,
		LogoPath: updatedParty.LogoPath,
	}, nil
}

func (m *Module) DeleteParty(ctx context.Context, id string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyUseCases.DeleteParty")
	defer span.End()

	partyID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	transaction := func(txCtx context.Context) (any, error) {
		party, err := m.partyRepo.GetPartyByID(txCtx, partyID)
		if err != nil {
			if errors.Is(err, dao.ErrNoResult) {
				return nil, &custerr.ErrChain{
					Message: "party not found",
					Cause:   err,
					Code:    404,
					Type:    response2.ErrNotFound,
				}
			}
			log.WithFields(log.Fields{
				"id":    id,
				"error": err,
			}).ErrorWithCtx(txCtx, "[PartyUseCases.DeleteParty] failed to get party by id")
			return nil, err
		}

		err = m.partyRepo.DeleteParty(txCtx, partyID)
		if err != nil {
			log.WithFields(log.Fields{
				"id":    id,
				"error": err,
			}).ErrorWithCtx(txCtx, "[PartyUseCases.DeleteParty] failed to delete party")
			return nil, err
		}
		if party.LogoPath != "" {
			err = fileutils.DeleteFile(txCtx, party.LogoPath)
			if err != nil {
				log.WithFields(log.Fields{
					"id":    id,
					"error": err,
				}).ErrorWithCtx(txCtx, "[PartyUseCases.DeleteParty] failed to delete party logo file")
			}
		}
		return nil, nil
	}

	_, err = m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *Module) GetPartyPhoto(ctx context.Context, id uuid.UUID) (*http.File, string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "PartyUseCases.GetPartyPhoto")
	defer span.End()

	party, err := m.partyRepo.GetPartyByID(ctx, id)
	if err != nil {
		return nil, "", &custerr.ErrChain{
			Message: "party not found",
			Code:    404,
			Type:    response2.ErrNotFound,
			Cause:   err,
		}
	}

	if party.LogoPath == "" {
		return nil, "", &custerr.ErrChain{
			Message: "no logo available for this party",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	file, contentType, err := filehandler.GetFileFromPath(ctx, party.LogoPath, filehandler.DisplayModeInline)
	if err != nil {
		return nil, "", &custerr.ErrChain{
			Message: "failed to retrieve logo",
			Code:    500,
			Type:    response2.ErrInternalServerError,
			Cause:   err,
		}
	}

	return file, contentType, nil
}

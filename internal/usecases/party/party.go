package party

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/domain/model"
	"github.com/nocturna-ta/election/internal/interfaces/dao"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
	"github.com/nocturna-ta/golib/custerr"
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

		errTx := m.partyRepo.InsertParty(txCtx, party)
		if errTx != nil {
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

		//publisher
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
	existing, err := m.partyRepo.GetPartyByID(ctx, partyID)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    req.ID,
			"error": err,
		}).ErrorWithCtx(ctx, "[PartyUseCases.UpdateParty] failed to get party by id")
		return nil, err
	}

	existing.Name = req.Name
	existing.LogoPath = req.LogoPath

	err = m.partyRepo.UpdateParty(ctx, party)
	if err != nil {
		log.WithFields(log.Fields{
			"id":    req.ID,
			"error": err,
		}).ErrorWithCtx(ctx, "[PartyUseCases.UpdateParty] failed to update party")
		return nil, err
	}

	return &response.PartyResponse{
		ID:       party.ID.String(),
		Name:     party.Name,
		LogoPath: party.LogoPath,
	}, nil
}

func (m *Module) DeleteParty(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

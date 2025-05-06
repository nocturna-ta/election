package supporting_party

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

func (m *Module) AddSupportingParty(ctx context.Context, req *request.AddSupportingPartyRequest) (*response.SupportingPartyResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "SupportingPartyUseCases.AddSupportingParty")
	defer span.End()

	var (
		supportingParty *model.SupportingParty
	)

	pairID, err := uuid.Parse(req.ElectionPairID)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid election pair ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	partyID, err := uuid.Parse(req.PartyID)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid party ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	_, err = m.electionRepo.GetElectionPairByID(ctx, pairID)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "Election pair not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"pairID": req.ElectionPairID,
			"error":  err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.AddSupportingParty] failed to get election pair")
		return nil, err
	}

	party, err := m.partyRepo.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "Party not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"partyID": req.PartyID,
			"error":   err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.AddSupportingParty] failed to get party")
		return nil, err
	}

	transaction := func(txCtx context.Context) (any, error) {
		supportingParty, err = model.ConstructSupportingParty(req)
		if err != nil {
			return nil, err
		}

		err = m.supportingPartyRepo.AddSupportingParty(txCtx, supportingParty)
		if err != nil {
			if errors.Is(err, dao.ErrDuplicate) {
				return nil, &custerr.ErrChain{
					Message: "Supporting party already exists",
					Code:    400,
					Type:    response2.ErrBadRequest,
					Cause:   err,
				}
			}

			log.WithFields(log.Fields{
				"error":           err,
				"supportingParty": *supportingParty,
			}).ErrorWithCtx(ctx, "[SupportingPartyUseCases.AddSupportingParty] failed to add supporting party")
			return nil, err
		}

		//publisher
		return nil, nil
	}

	_, err = m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return nil, err
	}

	return &response.SupportingPartyResponse{
		ID:             supportingParty.ID.String(),
		ElectionPairID: supportingParty.ElectionPairID.String(),
		PartyID:        supportingParty.PartyID.String(),
		Party: response.PartyResponse{
			ID:       party.ID.String(),
			Name:     party.Name,
			LogoPath: party.LogoPath,
		},
	}, nil
}

func (m *Module) GetSupportingPartiesByPairID(ctx context.Context, pairID string) (*response.SupportingPartiesResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "SupportingPartyUseCases.GetSupportingPartiesByPairID")
	defer span.End()

	id, err := uuid.Parse(pairID)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid election pair ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	_, err = m.electionRepo.GetElectionPairByID(ctx, id)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "Election pair not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"pairID": pairID,
			"error":  err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetSupportingPartiesByPairID] failed to get election pair")
		return nil, err
	}

	supportingParties, err := m.supportingPartyRepo.GetSupportingPartiesByPairID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"pairID": pairID,
			"error":  err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetSupportingPartiesByPairID] failed to get supporting parties")
		return nil, err
	}

	parties := make([]response.SupportingPartyResponse, len(supportingParties))
	for i, sp := range supportingParties {
		parties[i] = response.SupportingPartyResponse{
			ID:             sp.ID.String(),
			ElectionPairID: sp.ElectionPairID.String(),
			PartyID:        sp.PartyID.String(),
			Party: response.PartyResponse{
				ID:       sp.Party.ID.String(),
				Name:     sp.Party.Name,
				LogoPath: sp.Party.LogoPath,
			},
		}
	}

	return &response.SupportingPartiesResponse{
		Parties: parties,
		Total:   len(parties),
	}, nil
}

func (m *Module) RemoveSupportingParty(ctx context.Context, req *request.RemoveSupportingPartyRequest) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "SupportingPartyUseCases.RemoveSupportingParty")
	defer span.End()

	pairID, err := uuid.Parse(req.ElectionPairID)
	if err != nil {
		return &custerr.ErrChain{
			Message: "Invalid election pair ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	partyID, err := uuid.Parse(req.PartyID)
	if err != nil {
		return &custerr.ErrChain{
			Message: "Invalid party ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	_, err = m.electionRepo.GetElectionPairByID(ctx, pairID)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return &custerr.ErrChain{
				Message: "Election pair not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"pairID": req.ElectionPairID,
			"error":  err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.RemoveSupportingParty] failed to get election pair")
		return err
	}

	_, err = m.partyRepo.GetPartyByID(ctx, partyID)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return &custerr.ErrChain{
				Message: "Party not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"partyID": req.PartyID,
			"error":   err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.RemoveSupportingParty] failed to get party")
		return err
	}

	transaction := func(txCtx context.Context) (any, error) {
		err := m.supportingPartyRepo.RemoveSupportingParty(txCtx, pairID, partyID)
		if err != nil {
			log.WithFields(log.Fields{
				"pairID":  req.ElectionPairID,
				"partyID": req.PartyID,
				"error":   err,
			}).ErrorWithCtx(txCtx, "[ElectionUseCases.RemoveSupportingParty] failed to remove supporting party")
			return nil, err
		}

		//publisher

		return nil, nil
	}

	_, err = m.txMgr.Execute(ctx, transaction, nil)
	return err

}

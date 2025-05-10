package supporting_party

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
		supportingParty, errTx := model.ConstructSupportingParty(req)
		if err != nil {
			return nil, err
		}

		errTx = m.supportingPartyRepo.AddSupportingParty(txCtx, supportingParty)
		if errTx != nil {
			if errors.Is(errTx, dao.ErrDuplicate) {
				return nil, &custerr.ErrChain{
					Message: "Supporting party already exists",
					Code:    400,
					Type:    response2.ErrBadRequest,
					Cause:   errTx,
				}
			}

			log.WithFields(log.Fields{
				"error":           errTx,
				"supportingParty": *supportingParty,
			}).ErrorWithCtx(ctx, "[SupportingPartyUseCases.AddSupportingParty] failed to add supporting party")
			return nil, errTx
		}

		errTx = m.publisher.Publish(txCtx, m.topics.MasterDataSupportingParty.Value, supportingParty.ID.String(), supportingParty.ToMessageModel(), map[string]any{
			constants.MetaDataOperation: constants.Create,
		})
		if errTx != nil {
			return nil, errTx
		}
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

func (m *Module) GetSupportingPartiesByPairID(ctx context.Context, pairID uuid.UUID) (*response.SupportingPartiesResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "SupportingPartyUseCases.GetSupportingPartiesByPairID")
	defer span.End()

	_, err := m.electionRepo.GetElectionPairByID(ctx, pairID)
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

	supportingParties, err := m.supportingPartyRepo.GetSupportingPartiesByPairID(ctx, pairID)
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
		errTx := m.supportingPartyRepo.RemoveSupportingParty(txCtx, pairID, partyID)
		if errTx != nil {
			log.WithFields(log.Fields{
				"pairID":  req.ElectionPairID,
				"partyID": req.PartyID,
				"error":   errTx,
			}).ErrorWithCtx(txCtx, "[ElectionUseCases.RemoveSupportingParty] failed to remove supporting party")
			return nil, errTx
		}

		errTx = m.publisher.Publish(txCtx, m.topics.MasterDataSupportingParty.Value, partyID.String(), nil, map[string]any{
			constants.MetaDataOperation: constants.Delete,
		})

		return nil, nil
	}

	_, err = m.txMgr.Execute(ctx, transaction, nil)
	return err

}

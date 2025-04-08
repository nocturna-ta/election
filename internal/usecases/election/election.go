package election

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

func (m *Module) RegisterCandidate(ctx context.Context, req *request.CandidateRegistrationRequest) (*response.CandidateResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.RegisterCandidate")
	defer span.End()

	var (
		err       error
		candidate *model.Candidate
	)

	transaction := func(txCtx context.Context) (any, error) {
		candidate = model.ConstructRegistration(req)

		errTx := m.electionRepo.InsertCandidate(txCtx, candidate, req.SignedTransaction)
		if errTx != nil {
			if errors.Is(errTx, dao.ErrDuplicate) {
				return nil, &custerr.ErrChain{
					Message: "Election already exists",
					Code:    400,
					Cause:   errTx,
					Type:    response2.ErrBadRequest,
				}
			}
			return nil, errTx
		}

		errTx = m.publisher.Publish(txCtx, m.topics.MasterDataElection.Value, candidate.ID.String(), candidate.ToMessageModel(), map[string]any{
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

	return &response.CandidateResponse{
		ID:            candidate.ID.String(),
		NameCandidate: candidate.NameCandidate,
		ElectionNo:    candidate.ElectionNo,
		VoteCount:     candidate.VoteCount,
		IsActive:      candidate.IsActive,
	}, nil
}

func (m *Module) GetAllCandidate(ctx context.Context) (*[]response.CandidateResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetAllCandidate")
	defer span.End()

	candidates, err := m.electionRepo.GetAllCandidate(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "Failed to get all candidate")
		return nil, err
	}

	var candidateResponses []response.CandidateResponse
	for _, candidate := range candidates {
		candidateResponses = append(candidateResponses, response.CandidateResponse{
			ID:            candidate.ID.String(),
			NameCandidate: candidate.NameCandidate,
			ElectionNo:    candidate.ElectionNo,
			VoteCount:     candidate.VoteCount,
			IsActive:      candidate.IsActive,
		})
	}

	return &candidateResponses, nil
}

func (m *Module) GetCandidateByNo(ctx context.Context, no string) (*response.CandidateResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetCandidateByNo")
	defer span.End()

	candidate, err := m.electionRepo.GetCandidateByNo(ctx, no)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"no":    no,
		}).ErrorWithCtx(ctx, "Failed to get candidate by no")
		return nil, err
	}

	return &response.CandidateResponse{
		ID:            candidate.ID.String(),
		NameCandidate: candidate.NameCandidate,
		ElectionNo:    candidate.ElectionNo,
		VoteCount:     candidate.VoteCount,
		IsActive:      candidate.IsActive,
	}, nil
}

func (m *Module) ActivateCandidate(ctx context.Context, req *request.CandidateActivationRequest) (*response.CandidateActivation, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.ActivateCandidate")
	defer span.End()

	if err := m.electionRepo.CandidateActivate(ctx, req.ID, req.SignedTransaction); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "Failed to activate candidate")
		return nil, err
	}

	return &response.CandidateActivation{
		IsActive: true,
	}, nil
}

func (m *Module) UpsertCandidateDetail(ctx context.Context, req *request.CandidateDetailRequest) (*response.CandidateDetailResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.UpsertCandidateDetail")
	defer span.End()

	var (
		detail *model.CandidateDetail
	)

	transaction := func(txCtx context.Context) (any, error) {
		detail = model.ConstructUpsertCandidateDetail(req)

		errTx := m.electionRepo.UpsertCandidateDetail(txCtx, detail, uuid.MustParse(req.ID), uuid.MustParse(req.CandidateID))
		if errTx != nil {
			if errors.Is(errTx, dao.ErrDuplicate) {
				return nil, &custerr.ErrChain{
					Message: "Detail for this candidate already exists",
					Cause:   errTx,
					Code:    400,
					Type:    response2.ErrBadRequest,
				}
			}
			return nil, errTx
		}

		errTx = m.publisher.Publish(txCtx, m.topics.MasterDataElectionDetail.Value, detail.ID.String(), detail.ToMessageModel(), map[string]any{
			constants.MetaDataOperation: constants.Update,
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

	return &response.CandidateDetailResponse{
		ID:           detail.ID.String(),
		CandidateID:  detail.CandidateID.String(),
		Biodata:      detail.Biodata,
		Visi:         detail.Visi,
		Misi:         detail.Misi,
		ProgramKerja: detail.ProgramKerja,
	}, err
}

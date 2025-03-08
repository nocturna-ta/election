package election

import (
	"context"
	"errors"
	"github.com/nocturna-ta/election/internal/domain/model"
	"github.com/nocturna-ta/election/internal/interfaces/dao"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/log"
	response2 "github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/tracing"
)

func (m *Module) RegisterCandidate(ctx context.Context, req *request.CandidateRegistrationRequest) (*response.CandidateResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "UseCases.election.RegisterCandidate")
	defer span.End()

	var (
		err       error
		candidate *model.Candidate
	)

	if err := m.electionRepo.InsertCandidate(ctx, candidate, req.SignedTransaction); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "Failed to insert candidate")
	}

	if errors.Is(err, dao.ErrDuplicate) {
		return nil, &custerr.ErrChain{
			Message: "Duplicate candidate",
			Code:    400,
			Type:    response2.ErrBadRequest,
			Cause:   err,
		}
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
	span, ctx := tracing.StartSpanFromContext(ctx, "UseCases.election.GetAllCandidate")
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
	span, ctx := tracing.StartSpanFromContext(ctx, "UseCases.election.GetCandidateByNo")
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
	span, ctx := tracing.StartSpanFromContext(ctx, "UseCases.election.ActivateCandidate")
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

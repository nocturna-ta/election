package election

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

func (m *Module) RegisterElectionPair(ctx context.Context, req *request.ElectionPairRegistrationRequest) (*response.ElectionPairResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.RegisterElectionPair")
	defer span.End()

	var (
		electionPair *model.ElectionPair
		err          error
	)

	//reqCtx, err := libCtx.GetRequestContext(ctx)
	//if err != nil {
	//	return nil, err
	//}

	transaction := func(txCtx context.Context) (any, error) {
		electionPair = model.ConstructElectionPair(req)

		if errTx := m.electionRepo.InsertElectionPair(txCtx, electionPair, req.SignedTransaction); err != nil {
			if errors.Is(errTx, dao.ErrDuplicate) {
				return nil, &custerr.ErrChain{
					Message: "Election Pair already exists",
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

	_, err = m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return nil, err
	}

	return &response.ElectionPairResponse{
		ID:            electionPair.ID.String(),
		ElectionNo:    electionPair.ElectionNo,
		VoteCount:     electionPair.VoteCount,
		IsActive:      electionPair.IsActive,
		PairPhotoPath: electionPair.PairPhotoPath,
		President: response.CandidateInfoResponse{
			FullName:           electionPair.President.FullName,
			EducationHistory:   electionPair.President.EducationHistory,
			WorkExperience:     electionPair.President.WorkExperience,
			LegalRecordHistory: electionPair.President.LegalRecordHistory,
			PhotoPath:          electionPair.President.PhotoPath,
		},
		VicePresident: response.CandidateInfoResponse{
			FullName:           electionPair.VicePresident.FullName,
			EducationHistory:   electionPair.VicePresident.EducationHistory,
			WorkExperience:     electionPair.VicePresident.WorkExperience,
			LegalRecordHistory: electionPair.VicePresident.LegalRecordHistory,
			PhotoPath:          electionPair.VicePresident.PhotoPath,
		},
	}, err
}

func (m *Module) GetElectionPairByID(ctx context.Context, id string) (*response.ElectionPairResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetElectionPairByID")
	defer span.End()

	pairID, err := uuid.Parse(id)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	electionPair, err := m.electionRepo.GetElectionPairByID(ctx, pairID)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "Election Pair not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionUseCases] failed to get election pair by id")
		return nil, err
	}

	return &response.ElectionPairResponse{
		ID:            electionPair.ID.String(),
		ElectionNo:    electionPair.ElectionNo,
		VoteCount:     electionPair.VoteCount,
		IsActive:      electionPair.IsActive,
		PairPhotoPath: electionPair.PairPhotoPath,
		President: response.CandidateInfoResponse{
			FullName:           electionPair.President.FullName,
			EducationHistory:   electionPair.President.EducationHistory,
			WorkExperience:     electionPair.President.WorkExperience,
			LegalRecordHistory: electionPair.President.LegalRecordHistory,
			PhotoPath:          electionPair.President.PhotoPath,
		},
		VicePresident: response.CandidateInfoResponse{
			FullName:           electionPair.VicePresident.FullName,
			EducationHistory:   electionPair.VicePresident.EducationHistory,
			WorkExperience:     electionPair.VicePresident.WorkExperience,
			LegalRecordHistory: electionPair.VicePresident.LegalRecordHistory,
			PhotoPath:          electionPair.VicePresident.PhotoPath,
		},
	}, nil

}

func (m *Module) GetElectionPairByNo(ctx context.Context, no string) (*response.ElectionPairResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetElectionPairByNo")
	defer span.End()

	electionPair, err := m.electionRepo.GetElectionPairByNo(ctx, no)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "Election Pair not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"no":    no,
		}).ErrorWithCtx(ctx, "[ElectionUseCases] failed to get election pair by no")
		return nil, err
	}

	return &response.ElectionPairResponse{
		ID:            electionPair.ID.String(),
		ElectionNo:    electionPair.ElectionNo,
		VoteCount:     electionPair.VoteCount,
		IsActive:      electionPair.IsActive,
		PairPhotoPath: electionPair.PairPhotoPath,
		President: response.CandidateInfoResponse{
			FullName:           electionPair.President.FullName,
			EducationHistory:   electionPair.President.EducationHistory,
			WorkExperience:     electionPair.President.WorkExperience,
			LegalRecordHistory: electionPair.President.LegalRecordHistory,
			PhotoPath:          electionPair.President.PhotoPath,
		},
		VicePresident: response.CandidateInfoResponse{
			FullName:           electionPair.VicePresident.FullName,
			EducationHistory:   electionPair.VicePresident.EducationHistory,
			WorkExperience:     electionPair.VicePresident.WorkExperience,
			LegalRecordHistory: electionPair.VicePresident.LegalRecordHistory,
			PhotoPath:          electionPair.VicePresident.PhotoPath,
		},
	}, nil
}

func (m *Module) GetAllElectionPairs(ctx context.Context) (*response.ElectionPairListResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetAllElectionPairs")
	defer span.End()

	electionPairs, err := m.electionRepo.GetAllElectionPairs(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases] failed to get all election pairs")
		return nil, err
	}

	pairResponse := make([]response.ElectionPairResponse, len(electionPairs))
	for i, pair := range electionPairs {
		pairResponse[i] = response.ElectionPairResponse{
			ID:            pair.ID.String(),
			ElectionNo:    pair.ElectionNo,
			VoteCount:     pair.VoteCount,
			IsActive:      pair.IsActive,
			PairPhotoPath: pair.PairPhotoPath,
			President: response.CandidateInfoResponse{
				FullName:           pair.President.FullName,
				EducationHistory:   pair.President.EducationHistory,
				WorkExperience:     pair.President.WorkExperience,
				LegalRecordHistory: pair.President.LegalRecordHistory,
				PhotoPath:          pair.President.PhotoPath,
			},
			VicePresident: response.CandidateInfoResponse{
				FullName:           pair.VicePresident.FullName,
				EducationHistory:   pair.VicePresident.EducationHistory,
				WorkExperience:     pair.VicePresident.WorkExperience,
				LegalRecordHistory: pair.VicePresident.LegalRecordHistory,
				PhotoPath:          pair.VicePresident.PhotoPath,
			},
		}
	}

	return &response.ElectionPairListResponse{
		Pairs: pairResponse,
		Total: len(electionPairs),
	}, nil

}

func (m *Module) ActivateElectionPair(ctx context.Context, req *request.ElectionPairActivationRequest) (*response.ElectionPairActivationResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.ActivateElectionPair")
	defer span.End()

	pairID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	err = m.electionRepo.ActivateElectionPair(ctx, pairID, req.SignedTransaction)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "Election Pair not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    req.ID,
		}).ErrorWithCtx(ctx, "[ElectionUseCases] failed to activate election pair")
		return nil, err
	}

	return &response.ElectionPairActivationResponse{
		ID:       req.ID,
		IsActive: true,
	}, nil
}

func (m *Module) UpsertElectionPairDetail(ctx context.Context, req *request.ElectionPairDetailRequest) (*response.ElectionPairDetailResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.UpsertElectionPairDetail")
	defer span.End()

	var (
		detail *model.PairDetail
	)

	transaction := func(txCtx context.Context) (any, error) {
		detail = model.ConstructPairDetail(req)

		if err := m.electionRepo.UpsertPairDetail(txCtx, detail); err != nil {
			if errors.Is(err, dao.ErrDuplicate) {
				return nil, &custerr.ErrChain{
					Message: "Election Pair Detail already exists",
					Cause:   err,
					Code:    400,
					Type:    response2.ErrBadRequest,
				}
			}
			return nil, err
		}

		//publisher

		return nil, nil
	}
	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return nil, err
	}

	return &response.ElectionPairDetailResponse{
		ID:             detail.ID.String(),
		ElectionPairID: detail.ElectionPairID.String(),
		Vision:         detail.Vision,
		Mission:        detail.Mission,
		WorkProgram:    detail.WorkProgram,
		ProgramDocs:    detail.ProgramDocs,
	}, nil
}

func (m *Module) GetElectionPairDetail(ctx context.Context, pairID string) (*response.ElectionPairDetailResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetElectionPairDetail")
	defer span.End()

	id, err := uuid.Parse(pairID)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	detail, err := m.electionRepo.GetPairDetailByPairID(ctx, id)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "Election Pair Detail not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    pairID,
		}).ErrorWithCtx(ctx, "[ElectionUseCases] failed to get election pair detail")
		return nil, err
	}

	return &response.ElectionPairDetailResponse{
		ID:             detail.ID.String(),
		ElectionPairID: detail.ElectionPairID.String(),
		Vision:         detail.Vision,
		Mission:        detail.Mission,
		WorkProgram:    detail.WorkProgram,
		ProgramDocs:    detail.ProgramDocs,
	}, nil
}

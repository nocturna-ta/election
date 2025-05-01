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
	"github.com/nocturna-ta/golib/fileutils"
	"github.com/nocturna-ta/golib/http"
	"github.com/nocturna-ta/golib/http/filehandler"
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

		if electionPair.PairPhotoPath != "" {
			_ = fileutils.DeleteFile(txCtx, electionPair.PairPhotoPath)
		}
		if req.President.PhotoPath != "" {
			_ = fileutils.DeleteFile(txCtx, req.President.PhotoPath)
		}
		if req.VicePresident.PhotoPath != "" {
			_ = fileutils.DeleteFile(txCtx, req.VicePresident.PhotoPath)
		}

		fileConfigPairPhoto := fileutils.DefaultConfig()
		fileConfigPairPhoto.SetAllowedImageExtension()
		fileConfigPairPhoto.EntityType = "election_pair"

		photoPathPair, err := fileutils.StoreFile(txCtx, req.PairPhotoFile, req.PairPhotoName, fileConfigPairPhoto)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionUseCases] failed to store pair photo")
			return nil, &custerr.ErrChain{
				Message: "Failed to store pair photo",
				Cause:   err,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		fileConfigPresident := fileutils.DefaultConfig()
		fileConfigPresident.SetAllowedImageExtension()
		fileConfigPresident.EntityType = "election_pair"

		presidentPhoto, err := fileutils.StoreFile(txCtx, req.President.PhotoFile, req.President.PhotoName, fileConfigPresident)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionUseCases] failed to store president photo")
			return nil, &custerr.ErrChain{
				Message: "Failed to store president photo",
				Cause:   err,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		fileConfigVicePresident := fileutils.DefaultConfig()
		fileConfigVicePresident.SetAllowedImageExtension()
		fileConfigVicePresident.EntityType = "election_pair"

		vicePresidentPhoto, err := fileutils.StoreFile(txCtx, req.VicePresident.PhotoFile, req.VicePresident.PhotoName, fileConfigVicePresident)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).ErrorWithCtx(ctx, "[ElectionUseCases] failed to store vice president photo")
			return nil, &custerr.ErrChain{
				Message: "Failed to store vice president photo",
				Cause:   err,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		electionPair.PairPhotoPath = photoPathPair
		electionPair.President.PhotoPath = presidentPhoto
		electionPair.VicePresident.PhotoPath = vicePresidentPhoto

		if errTx := m.electionRepo.InsertElectionPair(txCtx, electionPair, req.SignedTransaction); err != nil {
			if errors.Is(errTx, dao.ErrDuplicate) {
				_ = fileutils.DeleteFile(txCtx, photoPathPair)
				_ = fileutils.DeleteFile(txCtx, presidentPhoto)
				_ = fileutils.DeleteFile(txCtx, vicePresidentPhoto)
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

	presidentEducationHistory := make([]response.EducationHistoryResponse, len(electionPair.President.EducationHistory))
	for i, history := range electionPair.President.EducationHistory {
		presidentEducationHistory[i] = response.EducationHistoryResponse{
			InstituteName: history.InstituteName,
			Year:          history.Year,
		}
	}

	presidentWorkExperience := make([]response.WorkHistoryResponse, len(electionPair.President.WorkExperience))
	for i, history := range electionPair.President.WorkExperience {
		presidentWorkExperience[i] = response.WorkHistoryResponse{
			InstituteName: history.InstituteName,
			Position:      history.Position,
			Year:          history.Year,
		}
	}

	vicePresidentEducationHistory := make([]response.EducationHistoryResponse, len(electionPair.VicePresident.EducationHistory))
	for i, history := range electionPair.VicePresident.EducationHistory {
		vicePresidentEducationHistory[i] = response.EducationHistoryResponse{
			InstituteName: history.InstituteName,
			Year:          history.Year,
		}
	}

	vicePresidentWorkExperience := make([]response.WorkHistoryResponse, len(electionPair.VicePresident.WorkExperience))
	for i, history := range electionPair.VicePresident.WorkExperience {
		vicePresidentWorkExperience[i] = response.WorkHistoryResponse{
			InstituteName: history.InstituteName,
			Position:      history.Position,
			Year:          history.Year,
		}
	}

	return &response.ElectionPairResponse{
		ID:            electionPair.ID.String(),
		ElectionNo:    electionPair.ElectionNo,
		VoteCount:     electionPair.VoteCount,
		IsActive:      electionPair.IsActive,
		PairPhotoPath: electionPair.PairPhotoPath,
		President: response.CandidateInfoResponse{
			FullName:         electionPair.President.FullName,
			EducationHistory: presidentEducationHistory,
			WorkExperience:   vicePresidentWorkExperience,
			Gender:           electionPair.President.Gender,
			BirthPlace:       electionPair.President.BirthPlace,
			BirthDate:        electionPair.President.BirthDate,
			Religion:         electionPair.President.Religion,
			LastEducation:    electionPair.President.LastEducation,
			Job:              electionPair.President.Job,
			PhotoPath:        electionPair.President.PhotoPath,
		},
		VicePresident: response.CandidateInfoResponse{
			FullName:         electionPair.VicePresident.FullName,
			EducationHistory: vicePresidentEducationHistory,
			WorkExperience:   vicePresidentWorkExperience,
			Gender:           electionPair.VicePresident.Gender,
			BirthPlace:       electionPair.VicePresident.BirthPlace,
			BirthDate:        electionPair.VicePresident.BirthDate,
			Religion:         electionPair.VicePresident.Religion,
			LastEducation:    electionPair.VicePresident.LastEducation,
			Job:              electionPair.VicePresident.Job,
			PhotoPath:        electionPair.VicePresident.PhotoPath,
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

	presidentEducationHistory := make([]response.EducationHistoryResponse, len(electionPair.President.EducationHistory))
	for i, history := range electionPair.President.EducationHistory {
		presidentEducationHistory[i] = response.EducationHistoryResponse{
			InstituteName: history.InstituteName,
			Year:          history.Year,
		}
	}

	presidentWorkExperience := make([]response.WorkHistoryResponse, len(electionPair.President.WorkExperience))
	for i, history := range electionPair.President.WorkExperience {
		presidentWorkExperience[i] = response.WorkHistoryResponse{
			InstituteName: history.InstituteName,
			Position:      history.Position,
			Year:          history.Year,
		}
	}

	vicePresidentEducationHistory := make([]response.EducationHistoryResponse, len(electionPair.VicePresident.EducationHistory))
	for i, history := range electionPair.VicePresident.EducationHistory {
		vicePresidentEducationHistory[i] = response.EducationHistoryResponse{
			InstituteName: history.InstituteName,
			Year:          history.Year,
		}
	}

	vicePresidentWorkExperience := make([]response.WorkHistoryResponse, len(electionPair.VicePresident.WorkExperience))
	for i, history := range electionPair.VicePresident.WorkExperience {
		vicePresidentWorkExperience[i] = response.WorkHistoryResponse{
			InstituteName: history.InstituteName,
			Position:      history.Position,
			Year:          history.Year,
		}
	}

	return &response.ElectionPairResponse{
		ID:            electionPair.ID.String(),
		ElectionNo:    electionPair.ElectionNo,
		VoteCount:     electionPair.VoteCount,
		IsActive:      electionPair.IsActive,
		PairPhotoPath: electionPair.PairPhotoPath,
		President: response.CandidateInfoResponse{
			FullName:         electionPair.President.FullName,
			EducationHistory: presidentEducationHistory,
			WorkExperience:   vicePresidentWorkExperience,
			Gender:           electionPair.President.Gender,
			BirthPlace:       electionPair.President.BirthPlace,
			BirthDate:        electionPair.President.BirthDate,
			Religion:         electionPair.President.Religion,
			LastEducation:    electionPair.President.LastEducation,
			Job:              electionPair.President.Job,
			PhotoPath:        electionPair.President.PhotoPath,
		},
		VicePresident: response.CandidateInfoResponse{
			FullName:         electionPair.VicePresident.FullName,
			EducationHistory: vicePresidentEducationHistory,
			WorkExperience:   vicePresidentWorkExperience,
			Gender:           electionPair.VicePresident.Gender,
			BirthPlace:       electionPair.VicePresident.BirthPlace,
			BirthDate:        electionPair.VicePresident.BirthDate,
			Religion:         electionPair.VicePresident.Religion,
			LastEducation:    electionPair.VicePresident.LastEducation,
			Job:              electionPair.VicePresident.Job,
			PhotoPath:        electionPair.VicePresident.PhotoPath,
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

	presidentEducationHistory := make([]response.EducationHistoryResponse, len(electionPair.President.EducationHistory))
	for i, history := range electionPair.President.EducationHistory {
		presidentEducationHistory[i] = response.EducationHistoryResponse{
			InstituteName: history.InstituteName,
			Year:          history.Year,
		}
	}

	presidentWorkExperience := make([]response.WorkHistoryResponse, len(electionPair.President.WorkExperience))
	for i, history := range electionPair.President.WorkExperience {
		presidentWorkExperience[i] = response.WorkHistoryResponse{
			InstituteName: history.InstituteName,
			Position:      history.Position,
			Year:          history.Year,
		}
	}

	vicePresidentEducationHistory := make([]response.EducationHistoryResponse, len(electionPair.VicePresident.EducationHistory))
	for i, history := range electionPair.VicePresident.EducationHistory {
		vicePresidentEducationHistory[i] = response.EducationHistoryResponse{
			InstituteName: history.InstituteName,
			Year:          history.Year,
		}
	}

	vicePresidentWorkExperience := make([]response.WorkHistoryResponse, len(electionPair.VicePresident.WorkExperience))
	for i, history := range electionPair.VicePresident.WorkExperience {
		vicePresidentWorkExperience[i] = response.WorkHistoryResponse{
			InstituteName: history.InstituteName,
			Position:      history.Position,
			Year:          history.Year,
		}
	}

	return &response.ElectionPairResponse{
		ID:            electionPair.ID.String(),
		ElectionNo:    electionPair.ElectionNo,
		VoteCount:     electionPair.VoteCount,
		IsActive:      electionPair.IsActive,
		PairPhotoPath: electionPair.PairPhotoPath,
		President: response.CandidateInfoResponse{
			FullName:         electionPair.President.FullName,
			EducationHistory: presidentEducationHistory,
			WorkExperience:   vicePresidentWorkExperience,
			Gender:           electionPair.President.Gender,
			BirthPlace:       electionPair.President.BirthPlace,
			BirthDate:        electionPair.President.BirthDate,
			Religion:         electionPair.President.Religion,
			LastEducation:    electionPair.President.LastEducation,
			Job:              electionPair.President.Job,
			PhotoPath:        electionPair.President.PhotoPath,
		},
		VicePresident: response.CandidateInfoResponse{
			FullName:         electionPair.VicePresident.FullName,
			EducationHistory: vicePresidentEducationHistory,
			WorkExperience:   vicePresidentWorkExperience,
			Gender:           electionPair.VicePresident.Gender,
			BirthPlace:       electionPair.VicePresident.BirthPlace,
			BirthDate:        electionPair.VicePresident.BirthDate,
			Religion:         electionPair.VicePresident.Religion,
			LastEducation:    electionPair.VicePresident.LastEducation,
			Job:              electionPair.VicePresident.Job,
			PhotoPath:        electionPair.VicePresident.PhotoPath,
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
		presidentEducationHistory := make([]response.EducationHistoryResponse, len(pair.President.EducationHistory))
		for j, history := range pair.President.EducationHistory {
			presidentEducationHistory[j] = response.EducationHistoryResponse{
				InstituteName: history.InstituteName,
				Year:          history.Year,
			}
		}

		presidentWorkExperience := make([]response.WorkHistoryResponse, len(pair.President.WorkExperience))
		for j, history := range pair.President.WorkExperience {
			presidentWorkExperience[j] = response.WorkHistoryResponse{
				InstituteName: history.InstituteName,
				Position:      history.Position,
				Year:          history.Year,
			}
		}

		vicePresidentEducationHistory := make([]response.EducationHistoryResponse, len(pair.VicePresident.EducationHistory))
		for j, history := range pair.VicePresident.EducationHistory {
			vicePresidentEducationHistory[j] = response.EducationHistoryResponse{
				InstituteName: history.InstituteName,
				Year:          history.Year,
			}
		}

		vicePresidentWorkExperience := make([]response.WorkHistoryResponse, len(pair.VicePresident.WorkExperience))
		for j, history := range pair.VicePresident.WorkExperience {
			vicePresidentWorkExperience[j] = response.WorkHistoryResponse{
				InstituteName: history.InstituteName,
				Position:      history.Position,
				Year:          history.Year,
			}
		}
		pairResponse[i] = response.ElectionPairResponse{
			ID:            pair.ID.String(),
			ElectionNo:    pair.ElectionNo,
			VoteCount:     pair.VoteCount,
			IsActive:      pair.IsActive,
			PairPhotoPath: pair.PairPhotoPath,
			President: response.CandidateInfoResponse{
				FullName:         pair.President.FullName,
				EducationHistory: presidentEducationHistory,
				WorkExperience:   presidentWorkExperience,
				Gender:           pair.President.Gender,
				BirthPlace:       pair.President.BirthPlace,
				BirthDate:        pair.President.BirthDate,
				Religion:         pair.President.Religion,
				LastEducation:    pair.President.LastEducation,
				Job:              pair.President.Job,
				PhotoPath:        pair.President.PhotoPath,
			},
			VicePresident: response.CandidateInfoResponse{
				FullName:         pair.VicePresident.FullName,
				EducationHistory: vicePresidentEducationHistory,
				WorkExperience:   vicePresidentWorkExperience,
				Gender:           pair.VicePresident.Gender,
				BirthPlace:       pair.VicePresident.BirthPlace,
				BirthDate:        pair.VicePresident.BirthDate,
				Religion:         pair.VicePresident.Religion,
				LastEducation:    pair.VicePresident.LastEducation,
				Job:              pair.VicePresident.Job,
				PhotoPath:        pair.VicePresident.PhotoPath,
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

func (m *Module) GetElectionPairPhoto(ctx context.Context, id uuid.UUID) (*http.File, string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetElectionPairPhotoPath")
	defer span.End()

	election, err := m.electionRepo.GetElectionPairByID(ctx, id)
	if err != nil {
		return nil, "", custerr.ErrChain{
			Message: "Election not found",
			Cause:   err,
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	if election.PairPhotoPath == "" {
		return nil, "", &custerr.ErrChain{
			Message: "Election Pair photo not found",
			Cause:   err,
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	file, contentType, err := filehandler.GetFileFromPath(ctx, election.PairPhotoPath, filehandler.DisplayModeAttachment)
	if err != nil {
		return nil, "", &custerr.ErrChain{
			Message: "Failed to get file",
			Cause:   err,
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	return file, contentType, nil
}

func (m *Module) GetPresidentPhoto(ctx context.Context, id uuid.UUID) (*http.File, string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetPresidentPhoto")
	defer span.End()
	election, err := m.electionRepo.GetElectionPairByID(ctx, id)
	if err != nil {
		return nil, "", custerr.ErrChain{
			Message: "Election not found",
			Cause:   err,
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	if election.President.PhotoPath == "" {
		return nil, "", &custerr.ErrChain{
			Message: "President photo not found",
			Cause:   err,
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	file, contentType, err := filehandler.GetFileFromPath(ctx, election.President.PhotoPath, filehandler.DisplayModeInline)
	if err != nil {
		return nil, "", &custerr.ErrChain{
			Message: "Failed to get file",
			Cause:   err,
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	return file, contentType, nil
}

func (m *Module) GetVicePresidentPhoto(ctx context.Context, id uuid.UUID) (*http.File, string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetVicePresidentPhoto")
	defer span.End()
	election, err := m.electionRepo.GetElectionPairByID(ctx, id)
	if err != nil {
		return nil, "", custerr.ErrChain{
			Message: "Election not found",
			Cause:   err,
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	if election.President.PhotoPath == "" {
		return nil, "", &custerr.ErrChain{
			Message: "President photo not found",
			Cause:   err,
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	file, contentType, err := filehandler.GetFileFromPath(ctx, election.President.PhotoPath, filehandler.DisplayModeAttachment)
	if err != nil {
		return nil, "", &custerr.ErrChain{
			Message: "Failed to get file",
			Cause:   err,
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	return file, contentType, nil
}

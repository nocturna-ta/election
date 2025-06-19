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
		txHash       string
	)

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

		photoPathPair, errTx := fileutils.StoreFile(txCtx, req.PairPhotoFile, req.PairPhotoName, fileConfigPairPhoto)
		if errTx != nil {
			log.WithFields(log.Fields{
				"error": errTx,
			}).ErrorWithCtx(ctx, "[ElectionUseCases.RegisterElectionPair] failed to store pair photo")
			return nil, &custerr.ErrChain{
				Message: "Failed to store pair photo",
				Cause:   errTx,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		fileConfigPresident := fileutils.DefaultConfig()
		fileConfigPresident.SetAllowedImageExtension()
		fileConfigPresident.EntityType = "election_pair"

		presidentPhoto, errTx := fileutils.StoreFile(txCtx, req.President.PhotoFile, req.President.PhotoName, fileConfigPresident)
		if errTx != nil {
			log.WithFields(log.Fields{
				"error": errTx,
			}).ErrorWithCtx(ctx, "[ElectionUseCases.RegisterElectionPair] failed to store president photo")
			return nil, &custerr.ErrChain{
				Message: "Failed to store president photo",
				Cause:   errTx,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		fileConfigVicePresident := fileutils.DefaultConfig()
		fileConfigVicePresident.SetAllowedImageExtension()
		fileConfigVicePresident.EntityType = "election_pair"

		vicePresidentPhoto, errTx := fileutils.StoreFile(txCtx, req.VicePresident.PhotoFile, req.VicePresident.PhotoName, fileConfigVicePresident)
		if errTx != nil {
			log.WithFields(log.Fields{
				"error": errTx,
			}).ErrorWithCtx(ctx, "[ElectionUseCases.RegisterElectionPair] failed to store vice president photo")
			return nil, &custerr.ErrChain{
				Message: "Failed to store vice president photo",
				Cause:   errTx,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		electionPair.PairPhotoPath = photoPathPair
		electionPair.President.PhotoPath = presidentPhoto
		electionPair.VicePresident.PhotoPath = vicePresidentPhoto

		if errTx = m.electionRepo.InsertElectionPair(txCtx, electionPair); err != nil {
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

		txHash, errTx = m.electionRepo.SendTxToBlockchain(txCtx, req.SignedTransaction)
		if errTx != nil {
			log.WithFields(log.Fields{
				"error": errTx,
			}).ErrorWithCtx(ctx, "[ElectionUseCases.RegisterElectionPair] failed to send transaction to blockchain")
			return nil, &custerr.ErrChain{
				Message: "Failed to send transaction to blockchain",
				Cause:   errTx,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		errTx = m.publisher.Publish(txCtx, m.topics.MasterDataElection.Value, electionPair.ID.String(), electionPair.ToMessageModel(txHash), map[string]any{
			constants.MetaDataOperation: constants.Create,
		})

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
		PairName:      electionPair.PairName,
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

func (m *Module) GetElectionPairByID(ctx context.Context, id uuid.UUID) (*response.ElectionPairResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetElectionPairByID")
	defer span.End()

	electionPair, err := m.electionRepo.GetElectionPairByID(ctx, id)
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
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetElectionPairByID] failed to get election pair by id")
		return nil, custerr.ErrChain{
			Message: "Failed to get election pair by ID",
			Cause:   err,
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	electionPairContract, err := m.electionContract.GetElection(nil, electionPair.ID.String())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetElectionPairByID] failed to get election pair contract")
		return nil, custerr.ErrChain{
			Message: "Failed to get election pair contract",
			Cause:   err,
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	if electionPair.ID.String() != electionPairContract.Id {
		log.WithFields(log.Fields{
			"id":          id,
			"contract_id": electionPairContract.Id,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetElectionPairByID] election pair ID mismatch with contract ID")
		return nil, &custerr.ErrChain{
			Message: "Election Pair ID mismatch with contract ID",
			Cause:   errors.New("ID mismatch"),
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
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
		PairName:      electionPair.PairName,
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
		PairName:      electionPair.PairName,
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
func (m *Module) GetAllElectionPairs(ctx context.Context) (*[]response.ElectionPairFullResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetAllElectionPairs")
	defer span.End()

	electionPairs, err := m.electionRepo.GetAllElectionPairs(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetAllElectionPairs] failed to get all election pairs")
		return nil, &custerr.ErrChain{
			Message: "Failed to get all election pairs",
			Cause:   err,
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	electionPairsContract, err := m.electionContract.GetAllElection(nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetAllElectionPairs] failed to get all election pairs contract")
		return nil, &custerr.ErrChain{ // Fixed: Added missing pointer
			Message: "Failed to get all election pairs contract",
			Cause:   err,
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	// Create map for O(1) lookup
	contractId := make(map[string]bool)
	for _, contract := range electionPairsContract {
		contractId[contract.Id] = true
	}

	// Use dynamic slices instead of pre-allocated arrays
	var fullResponse []response.ElectionPairFullResponse

	for _, pair := range electionPairs {
		// Only process pairs that exist in contract
		if !contractId[pair.ID.String()] {
			continue
		}

		// Process president education history
		presidentEducationHistory := make([]response.EducationHistoryResponse, len(pair.President.EducationHistory))
		for j, history := range pair.President.EducationHistory {
			presidentEducationHistory[j] = response.EducationHistoryResponse{
				InstituteName: history.InstituteName,
				Year:          history.Year,
			}
		}

		// Process president work experience
		presidentWorkExperience := make([]response.WorkHistoryResponse, len(pair.President.WorkExperience))
		for j, history := range pair.President.WorkExperience {
			presidentWorkExperience[j] = response.WorkHistoryResponse{
				InstituteName: history.InstituteName,
				Position:      history.Position,
				Year:          history.Year,
			}
		}

		// Process vice president education history
		vicePresidentEducationHistory := make([]response.EducationHistoryResponse, len(pair.VicePresident.EducationHistory))
		for j, history := range pair.VicePresident.EducationHistory {
			vicePresidentEducationHistory[j] = response.EducationHistoryResponse{
				InstituteName: history.InstituteName,
				Year:          history.Year,
			}
		}

		// Process vice president work experience
		vicePresidentWorkExperience := make([]response.WorkHistoryResponse, len(pair.VicePresident.WorkExperience))
		for j, history := range pair.VicePresident.WorkExperience {
			vicePresidentWorkExperience[j] = response.WorkHistoryResponse{
				InstituteName: history.InstituteName,
				Position:      history.Position,
				Year:          history.Year,
			}
		}

		detail, err := m.electionRepo.GetPairDetailByPairID(ctx, pair.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"pair_id": pair.ID.String(),
			}).WarnWithCtx(ctx, "[ElectionUseCases.GetAllElectionPairs] failed to get pair detail, using default values")
		}

		// Build detail response
		var detailResponse response.ElectionPairDetailResponse
		if detail == nil {
			detailResponse = response.ElectionPairDetailResponse{
				ID:             "",
				ElectionPairID: pair.ID.String(),
				Vision:         "",
				Mission:        "",
				WorkProgram:    nil,
				ProgramDocs:    "",
			}
		} else {
			workProgram := make([]response.WorkProgramResponse, len(detail.WorkProgram))
			for i, program := range detail.WorkProgram {
				workProgram[i] = response.WorkProgramResponse{
					ProgramName:  program.ProgramName,
					ProgramPhoto: program.ProgramPhoto,
					ProgramDesc:  program.ProgramDesc,
				}
			}
			detailResponse = response.ElectionPairDetailResponse{
				ID:             detail.ID.String(),
				ElectionPairID: pair.ID.String(),
				Vision:         detail.Vision,
				Mission:        detail.Mission,
				WorkProgram:    workProgram,
				ProgramDocs:    detail.ProgramDocs,
			}
		}

		pairResponse := response.ElectionPairResponse{
			ID:            pair.ID.String(),
			ElectionNo:    pair.ElectionNo,
			VoteCount:     pair.VoteCount,
			IsActive:      pair.IsActive,
			PairName:      pair.PairName,
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

		fullPairResponse := response.ElectionPairFullResponse{
			ElectionPairResponse: pairResponse,
			Detail:               detailResponse,
			SupportingParties:    nil,
		}

		fullResponse = append(fullResponse, fullPairResponse)
	}

	return &fullResponse, nil
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

	transaction := func(txCtx context.Context) (any, error) {
		electionPair, errTx := m.electionRepo.GetElectionPairByID(txCtx, pairID)
		if errTx != nil {
			if errors.Is(errTx, dao.ErrNoResult) {
				return nil, &custerr.ErrChain{
					Message: "Election Pair not found",
					Cause:   errTx,
					Code:    404,
					Type:    response2.ErrNotFound,
				}
			}
			return nil, errTx
		}

		if electionPair.IsActive {
			return nil, &custerr.ErrChain{
				Message: "Election Pair already active",
				Cause:   errTx,
				Code:    400,
				Type:    response2.ErrBadRequest,
			}
		}

		errTx = m.electionRepo.ActivateElectionPair(txCtx, electionPair.ID)
		if errTx != nil {
			if errors.Is(errTx, dao.ErrNoResult) {
				return nil, &custerr.ErrChain{
					Message: "Election Pair not found",
					Cause:   errTx,
					Code:    400,
					Type:    response2.ErrBadRequest,
				}
			}
			log.WithFields(log.Fields{
				"error": errTx,
				"id":    req.ID,
			}).ErrorWithCtx(ctx, "[ElectionUseCases.ActivateElectionPair] failed to activate election pair")
			return nil, &custerr.ErrChain{
				Message: "Failed to activate election pair",
				Cause:   errTx,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		txHash, errTx := m.electionRepo.SendTxToBlockchain(txCtx, req.SignedTransaction)
		if errTx != nil {
			log.WithFields(log.Fields{
				"error": errTx,
			}).ErrorWithCtx(ctx, "[ElectionUseCases.ActivateElectionPair] failed to send transaction to blockchain")
			return nil, &custerr.ErrChain{
				Message: "Failed to send transaction to blockchain",
				Cause:   errTx,
				Code:    500,
				Type:    response2.ErrInternalServerError,
			}
		}

		message := model.CreateElectionActivationMessage(electionPair.ID.String(), electionPair.IsActive, txHash)
		errTx = m.publisher.Publish(txCtx, m.topics.MasterDataElection.Value, electionPair.ID.String(), message, map[string]any{
			constants.MetaDataOperation: constants.Update,
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
	return &response.ElectionPairActivationResponse{
		ID:       req.ID,
		IsActive: true,
	}, nil
}

func (m *Module) UpsertElectionPairDetail(ctx context.Context, req *request.ElectionPairDetailRequest) (*response.ElectionPairDetailResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.UpsertElectionPairDetail")
	defer span.End()

	var (
		detail   *model.PairDetail
		existing *model.PairDetail
		err      error
	)

	pairID, err := uuid.Parse(req.ElectionPairID)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	existing, err = m.electionRepo.GetPairDetailByPairID(ctx, pairID)
	if err != nil && !errors.Is(err, dao.ErrNoResult) {
		log.WithFields(log.Fields{
			"error": err,
			"id":    req.ElectionPairID,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.UpsertElectionPairDetail] failed to get existing election pair detail")
		return nil, err
	}

	transaction := func(txCtx context.Context) (any, error) {
		detail = model.ConstructPairDetail(req)

		var oldProgramDocsPath string
		oldWorkProgramPhotos := make(map[string]string)

		if existing != nil {
			detail.ID = existing.ID
			oldProgramDocsPath = existing.ProgramDocs

			for _, program := range existing.WorkProgram {
				if program.ProgramPhoto != "" {
					oldWorkProgramPhotos[program.ProgramName] = program.ProgramPhoto
				}
			}
		}

		newlyStoredFiles := make([]string, 0)

		if req.ProgramDocsFile != nil {
			fileConfig := fileutils.DefaultConfig()
			fileConfig.SetAllowedDocumentExtensions()
			fileConfig.EntityType = "election_pair"

			programDocs, err := fileutils.StoreFile(txCtx, req.ProgramDocsFile, req.ProgramDocsName, fileConfig)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).ErrorWithCtx(ctx, "[ElectionUseCases.UpsertElectionPairDetail] failed to store program docs")
				return nil, &custerr.ErrChain{
					Message: "Failed to store program docs",
					Cause:   err,
					Code:    500,
					Type:    response2.ErrInternalServerError,
				}
			}

			detail.ProgramDocs = programDocs
			newlyStoredFiles = append(newlyStoredFiles, programDocs)
		} else if existing != nil {
			detail.ProgramDocs = existing.ProgramDocs
			_ = fileutils.DeleteFile(txCtx, oldProgramDocsPath)
			oldProgramDocsPath = ""
		}

		for i, program := range req.WorkProgram {
			if program.ProgramPhotoFile != nil {
				fileConfigPhoto := fileutils.DefaultConfig()
				fileConfigPhoto.SetAllowedImageExtension()
				fileConfigPhoto.EntityType = "election_pair"

				photoPath, err := fileutils.StoreFile(txCtx, program.ProgramPhotoFile, program.ProgramPhotoName, fileConfigPhoto)
				if err != nil {
					for _, path := range newlyStoredFiles {
						_ = fileutils.DeleteFile(txCtx, path)
					}
					log.WithFields(log.Fields{
						"error": err,
						"index": i,
					}).ErrorWithCtx(ctx, "[ElectionUseCases.UpsertElectionPairDetail] failed to store work program photo")
					return nil, &custerr.ErrChain{
						Message: "Failed to store work program photo",
						Cause:   err,
						Code:    500,
						Type:    response2.ErrInternalServerError,
					}
				}

				detail.WorkProgram[i].ProgramPhoto = photoPath
				newlyStoredFiles = append(newlyStoredFiles, photoPath)

				if _, exists := oldWorkProgramPhotos[program.ProgramName]; exists {
					_ = fileutils.DeleteFile(txCtx, oldWorkProgramPhotos[program.ProgramName])
					delete(oldWorkProgramPhotos, program.ProgramName)
				}
			} else if existing != nil {
				for _, existingProgram := range existing.WorkProgram {
					if existingProgram.ProgramName == program.ProgramName && existingProgram.ProgramPhoto != "" {
						detail.WorkProgram[i].ProgramPhoto = existingProgram.ProgramPhoto
						delete(oldWorkProgramPhotos, program.ProgramName)
						break
					}
				}
			}
		}

		if err := m.electionRepo.UpsertPairDetail(txCtx, detail); err != nil {
			for _, path := range newlyStoredFiles {
				_ = fileutils.DeleteFile(txCtx, path)
			}

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

		if oldProgramDocsPath != "" && oldProgramDocsPath != detail.ProgramDocs {
			_ = fileutils.DeleteFile(txCtx, oldProgramDocsPath)
		}

		for _, oldPath := range oldWorkProgramPhotos {
			_ = fileutils.DeleteFile(txCtx, oldPath)
		}

		return nil, nil
	}

	_, err = m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		return nil, err
	}

	workProgram := make([]response.WorkProgramResponse, len(detail.WorkProgram))
	for i, program := range detail.WorkProgram {
		workProgram[i] = response.WorkProgramResponse{
			ProgramName:  program.ProgramName,
			ProgramPhoto: program.ProgramPhoto,
			ProgramDesc:  program.ProgramDesc,
		}
	}

	return &response.ElectionPairDetailResponse{
		ID:             detail.ID.String(),
		ElectionPairID: detail.ElectionPairID.String(),
		Vision:         detail.Vision,
		Mission:        detail.Mission,
		WorkProgram:    workProgram,
		ProgramDocs:    detail.ProgramDocs,
	}, nil
}

func (m *Module) GetElectionPairDetail(ctx context.Context, pairID uuid.UUID) (*response.ElectionPairDetailResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetElectionPairDetail")
	defer span.End()

	detail, err := m.electionRepo.GetPairDetailByPairID(ctx, pairID)
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
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetElectionPairDetail] failed to get election pair detail")
		return nil, err
	}

	workProgram := make([]response.WorkProgramResponse, len(detail.WorkProgram))
	for i, program := range detail.WorkProgram {
		workProgram[i] = response.WorkProgramResponse{
			ProgramName:  program.ProgramName,
			ProgramPhoto: program.ProgramPhoto,
			ProgramDesc:  program.ProgramDesc,
		}
	}

	return &response.ElectionPairDetailResponse{
		ID:             detail.ID.String(),
		ElectionPairID: detail.ElectionPairID.String(),
		Vision:         detail.Vision,
		Mission:        detail.Mission,
		WorkProgram:    workProgram,
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

	file, contentType, err := filehandler.GetFileFromPath(ctx, election.PairPhotoPath, filehandler.DisplayModeInline)
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

	file, contentType, err := filehandler.GetFileFromPath(ctx, election.VicePresident.PhotoPath, filehandler.DisplayModeInline)
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

func (m *Module) GetElectionPairFull(ctx context.Context, id uuid.UUID) (*response.ElectionPairFullResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetElectionPairFull")
	defer span.End()

	pairResponse, err := m.GetElectionPairByID(ctx, id)
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
			"id":    id.String(),
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetElectionPairFull] failed to get election pair by id")
		return nil, err
	}

	fullResp := &response.ElectionPairFullResponse{
		ElectionPairResponse: *pairResponse,
	}

	detailResp, err := m.GetElectionPairDetail(ctx, id)
	if err == nil {
		fullResp.Detail = *detailResp
	} else if !errors.Is(err, dao.ErrNoResult) {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id.String(),
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetElectionPairFull] failed to get election pair detail")
		return nil, err
	}

	partyResponses, err := m.supportingPartyUC.GetSupportingPartiesByPairID(ctx, id)
	if err == nil {
		fullResp.SupportingParties = partyResponses.Parties
	} else {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id.String(),
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetElectionPairFull] failed to get supporting parties")
		return nil, err
	}

	return fullResp, nil
}

func (m *Module) GetWorkProgramFile(ctx context.Context, id uuid.UUID) (*http.File, string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionUseCases.GetWorkProgramFile")
	defer span.End()

	detail, err := m.electionRepo.GetPairDetailByPairID(ctx, id)
	if err != nil {
		if errors.Is(err, dao.ErrNoResult) {
			return nil, "", &custerr.ErrChain{
				Message: "Election Pair Detail not found",
				Cause:   err,
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ElectionUseCases.GetWorkProgramFile] failed to get election pair detail")
		return nil, "", err
	}

	if detail.ProgramDocs == "" {
		return nil, "", &custerr.ErrChain{
			Message: "Work program file not found",
			Cause:   errors.New("file not found"),
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	file, contentType, err := filehandler.GetFileFromPath(ctx, detail.ProgramDocs, filehandler.DisplayModeAttachment)
	if err != nil {
		return nil, "", &custerr.ErrChain{
			Message: "Failed to get work program file",
			Cause:   err,
			Code:    500,
			Type:    response2.ErrInternalServerError,
		}
	}

	return file, contentType, nil
}

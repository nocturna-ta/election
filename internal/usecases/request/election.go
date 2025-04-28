package request

import (
	"github.com/nocturna-ta/election/pkg/utils"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
)

type CandidateInfoRequest struct {
	FullName           string `json:"full_name"`
	EducationHistory   string `json:"education_history"`
	WorkExperience     string `json:"work_experience"`
	LegalRecordHistory string `json:"legal_record_history"`
	PhotoPath          string `json:"photo_path"`
}

type ElectionPairRegistrationRequest struct {
	ID                string               `json:"id"`
	ElectionNo        string               `json:"election_no"`
	PairPhotoPath     string               `json:"pair_photo_path"`
	President         CandidateInfoRequest `json:"president"`
	VicePresident     CandidateInfoRequest `json:"vice_president"`
	SignedTransaction string               `json:"signed_transaction"`
}

type ElectionPairActivationRequest struct {
	ID                string `json:"id"`
	SignedTransaction string `json:"signed_transaction"`
}

type ElectionPairDetailRequest struct {
	ID             string   `json:"id"`
	ElectionPairID string   `json:"election_pair_id"`
	Vision         string   `json:"vision"`
	Mission        string   `json:"mission"`
	WorkProgram    string   `json:"work_program"`
	ProgramDocs    []string `json:"program_docs"`
}

func (req *ElectionPairRegistrationRequest) ValidateRegistrationRequest() error {
	if req == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if utils.IsNotUUID(req.ID) {
		return &custerr.ErrChain{
			Message: "ID is not a valid UUID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if req.ElectionNo == "" {
		return &custerr.ErrChain{
			Message: "Election number cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if req.President.FullName == "" {
		return &custerr.ErrChain{
			Message: "President candidate name cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if req.VicePresident.FullName == "" {
		return &custerr.ErrChain{
			Message: "Vice President candidate name cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if req.SignedTransaction == "" {
		return &custerr.ErrChain{
			Message: "Signed transaction cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

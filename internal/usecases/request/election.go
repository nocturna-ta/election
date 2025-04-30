package request

import (
	"github.com/nocturna-ta/election/pkg/utils"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"io"
)

type CandidateInfoRequest struct {
	FullName         string                    `json:"full_name"`
	EducationHistory []EducationHistoryRequest `json:"education_history"`
	WorkExperience   []WorkHistoryRequest      `json:"work_experience"`
	Gender           string                    `json:"gender"`
	BirthPlace       string                    `json:"birth_place"`
	BirthDate        string                    `json:"birth_date"`
	Religion         string                    `json:"religion"`
	LastEducation    string                    `json:"last_education"`
	Job              string                    `json:"job"`
	PhotoPath        string                    `json:"photo_path"`
	PhotoFile        io.Reader                 `json:"photo_file"`
	PhotoName        string                    `json:"photo_name"`
}

type EducationHistoryRequest struct {
	InstituteName string `json:"institute_name"`
	Year          string `json:"year"`
}

type WorkHistoryRequest struct {
	InstituteName string `json:"institute_name"`
	Position      string `json:"position"`
	Year          string `json:"year"`
}

type ElectionPairRegistrationRequest struct {
	ID                string               `json:"id"`
	ElectionNo        string               `json:"election_no"`
	PairPhotoPath     string               `json:"pair_photo_path"`
	PairPhotoFile     io.Reader            `json:"pair_photo_file"`
	PairPhotoName     string               `json:"pair_photo_name"`
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

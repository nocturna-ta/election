package request

import (
	"github.com/nocturna-ta/election/pkg/common"
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
	PhotoFile        io.Reader                 `json:"-" swaggerignore:"true"`
	PhotoName        string                    `json:"-" swaggerignore:"true"`
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
	PairName          string               `json:"pair_name"`
	PairPhotoPath     string               `json:"pair_photo_path"`
	PairPhotoFile     io.Reader            `json:"-" swaggerignore:"true"`
	PairPhotoName     string               `json:"-" swaggerignore:"true"`
	President         CandidateInfoRequest `json:"president"`
	VicePresident     CandidateInfoRequest `json:"vice_president"`
	SignedTransaction string               `json:"signed_transaction"`
}

type ElectionPairActivationRequest struct {
	ID                string `json:"id"`
	SignedTransaction string `json:"signed_transaction"`
}

type WorkProgramRequest struct {
	ProgramName      string    `json:"program_name"`
	ProgramPhoto     string    `json:"program_photo"`
	ProgramPhotoFile io.Reader `json:"-" swaggerignore:"true"`
	ProgramPhotoName string    `json:"-" swaggerignore:"true"`
	ProgramDesc      []string  `json:"program_desc"`
}

type ElectionPairDetailRequest struct {
	ElectionPairID  string               `json:"election_pair_id"`
	Vision          string               `json:"vision"`
	Mission         string               `json:"mission"`
	WorkProgram     []WorkProgramRequest `json:"work_program"`
	ProgramDocs     string               `json:"program_docs"`
	ProgramDocsFile io.Reader            `json:"-" swaggerignore:"true"`
	ProgramDocsName string               `json:"-" swaggerignore:"true"`
}

func (req *ElectionPairRegistrationRequest) ValidateRegistrationRequest() error {
	if req == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if common.IsNotUUID(req.ID) {
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

func (req *ElectionPairDetailRequest) ValidateDetailRequest() error {
	if req == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

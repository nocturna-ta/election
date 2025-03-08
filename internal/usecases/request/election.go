package request

import (
	"github.com/nocturna-ta/election/pkg/utils"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
)

type CandidateRegistrationRequest struct {
	ID                string   `json:"id"`
	NameCandidate     []string `json:"name_candidate"`
	ElectionNo        string   `json:"election_no"`
	SignedTransaction string   `json:"signed_transaction"`
}

type CandidateActivationRequest struct {
	ID                string `json:"id"`
	SignedTransaction string `json:"signed_transaction"`
}

func (req *CandidateRegistrationRequest) ValidateRegistrationRequest() error {
	if req == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if utils.IsNotUUID(req.ID) {
		return &custerr.ErrChain{
			Message: "ID is not valid",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

func (req *CandidateActivationRequest) ValidateActivationRequest() error {
	if req == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if utils.IsNotUUID(req.ID) {
		return &custerr.ErrChain{
			Message: "ID is not valid",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

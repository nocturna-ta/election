package request

import (
	"github.com/nocturna-ta/election/pkg/common"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
)

type AddSupportingPartyRequest struct {
	ElectionPairID string `json:"election_pair_id"`
	PartyID        string `json:"party_id"`
}

type RemoveSupportingPartyRequest struct {
	ElectionPairID string `json:"election_pair_id"`
	PartyID        string `json:"party_id"`
}

func (r *AddSupportingPartyRequest) Validate() error {
	if r == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if r.ElectionPairID == "" {
		return &custerr.ErrChain{
			Message: "Election Pair ID cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if common.IsNotUUID(r.ElectionPairID) {
		return &custerr.ErrChain{
			Message: "Election Pair ID is not a valid UUID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if r.PartyID == "" {
		return &custerr.ErrChain{
			Message: "Party ID cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if common.IsNotUUID(r.PartyID) {
		return &custerr.ErrChain{
			Message: "Party ID is not a valid UUID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

func (r *RemoveSupportingPartyRequest) Validate() error {
	if r == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if r.ElectionPairID == "" {
		return &custerr.ErrChain{
			Message: "Election Pair ID cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if common.IsNotUUID(r.ElectionPairID) {
		return &custerr.ErrChain{
			Message: "Election Pair ID is not a valid UUID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if r.PartyID == "" {
		return &custerr.ErrChain{
			Message: "Party ID cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if common.IsNotUUID(r.PartyID) {
		return &custerr.ErrChain{
			Message: "Party ID is not a valid UUID",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

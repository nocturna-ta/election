package model

import (
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/golib/custerr"
	response2 "github.com/nocturna-ta/golib/response"
	"time"
)

type SupportingParty struct {
	BaseModel
	ID             uuid.UUID `db:"id"`
	ElectionPairID uuid.UUID `db:"election_pair_id"`
	PartyID        uuid.UUID `db:"party_id"`
	Party          *Party    `db:"-"`
}

type Party struct {
	BaseModel
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	LogoPath string    `db:"logo_path"`
}

func ConstructPartyRegistration(req *request.PartyRegisterRequest) *Party {
	now := time.Now()
	id := uuid.New()

	return &Party{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:       id,
		Name:     req.Name,
		LogoPath: req.LogoPath,
	}

}

func ConstructSupportingParty(req *request.AddSupportingPartyRequest) (*SupportingParty, error) {
	now := time.Now()

	pairID, err := uuid.Parse(req.ElectionPairID)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid election pair ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	partyID, err := uuid.Parse(req.PartyID)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid party ID format",
			Cause:   err,
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	return &SupportingParty{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:             uuid.New(),
		ElectionPairID: pairID,
		PartyID:        partyID,
	}, nil
}

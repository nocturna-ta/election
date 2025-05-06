package model

import (
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/usecases/request"
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

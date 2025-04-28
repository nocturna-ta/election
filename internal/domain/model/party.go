package model

import "github.com/google/uuid"

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

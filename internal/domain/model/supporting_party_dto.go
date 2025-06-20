package model

import (
	"github.com/google/uuid"
	"time"
)

type SupportingPartyDTO struct {
	ID             uuid.UUID `db:"id"`
	ElectionPairID uuid.UUID `db:"election_pair_id"`
	PartyID        uuid.UUID `db:"party_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	IsDeleted      bool      `db:"is_deleted"`
	PartyID2       uuid.UUID `db:"party_id2"`
	PartyName      string    `db:"party_name"`
	PartyLogoPath  string    `db:"party_logo_path"`
	PartyCreatedAt time.Time `db:"party_created_at"`
	PartyUpdatedAt time.Time `db:"party_updated_at"`
	PartyIsDeleted bool      `db:"party_is_deleted"`
}

func (dto *SupportingPartyDTO) ToDomain() *SupportingParty {
	return &SupportingParty{
		BaseModel: BaseModel{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
			IsDeleted: dto.IsDeleted,
		},
		ID:             dto.ID,
		ElectionPairID: dto.ElectionPairID,
		PartyID:        dto.PartyID,
		Party: &Party{
			BaseModel: BaseModel{
				CreatedAt: dto.PartyCreatedAt,
				UpdatedAt: dto.PartyUpdatedAt,
				IsDeleted: dto.PartyIsDeleted,
			},
			ID:       dto.PartyID2,
			Name:     dto.PartyName,
			LogoPath: dto.PartyLogoPath,
		},
	}
}

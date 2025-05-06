package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/domain/model"
)

type SupportingPartyRepository interface {
	AddSupportingParty(ctx context.Context, supportingParty *model.SupportingParty) error
	GetSupportingPartiesByPairID(ctx context.Context, pairID uuid.UUID) ([]model.SupportingParty, error)
	RemoveSupportingParty(ctx context.Context, pairID uuid.UUID, partyID uuid.UUID) error
}

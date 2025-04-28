package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/domain/model"
)

type PartyRepository interface {
	AddSupportingParty(ctx context.Context, supportingParty *model.SupportingParty) error
	GetSupportingPartiesByPairID(ctx context.Context, pairID uuid.UUID) ([]model.SupportingParty, error)
	GetPartyByID(ctx context.Context, id uuid.UUID) (*model.Party, error)
	GetAllParties(ctx context.Context) ([]model.Party, error)
}

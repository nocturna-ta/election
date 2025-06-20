package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/domain/model"
)

type PartyRepository interface {
	GetPartyByID(ctx context.Context, id uuid.UUID) (*model.Party, error)
	GetAllParties(ctx context.Context) ([]model.Party, error)
	InsertParty(ctx context.Context, party *model.Party) error
	UpdateParty(ctx context.Context, party *model.Party) (*model.Party, error)
	DeleteParty(ctx context.Context, id uuid.UUID) error
}

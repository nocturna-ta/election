package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/domain/model"
)

type ElectionRepository interface {
	InsertElectionPair(ctx context.Context, pair *model.ElectionPair) error
	UpdateElectionPair(ctx context.Context, pair *model.ElectionPair) error
	GetElectionPairByID(ctx context.Context, id uuid.UUID) (*model.ElectionPair, error)
	GetElectionPairByNo(ctx context.Context, no string) (*model.ElectionPair, error)
	GetAllElectionPairs(ctx context.Context) ([]model.ElectionPair, error)
	ActivateElectionPair(ctx context.Context, id uuid.UUID) error

	UpdateElectionPairPhoto(ctx context.Context, id uuid.UUID, photoPath string) error
	UpdatePresidentPhoto(ctx context.Context, id uuid.UUID, photoPath string) error
	UpdateVicePresidentPhoto(ctx context.Context, id uuid.UUID, photoPath string) error

	UpsertPairDetail(ctx context.Context, detail *model.PairDetail) error
	GetPairDetailByPairID(ctx context.Context, pairID uuid.UUID) (*model.PairDetail, error)

	SendTxToBlockchain(ctx context.Context, signedTransaction string) (string, error)
}

package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/domain/model"
)

type ElectionRepository interface {
	InsertCandidate(ctx context.Context, candidate *model.Candidate, signedTransaction string) error
	UpsertCandidateDetail(ctx context.Context, detail *model.CandidateDetail, id uuid.UUID, candidateId uuid.UUID) error
	GetAllCandidate(ctx context.Context) ([]model.Candidate, error)
	GetCandidateByNo(ctx context.Context, no string) (*model.Candidate, error)
	CandidateActivate(ctx context.Context, id string, signedTransaction string) error
}

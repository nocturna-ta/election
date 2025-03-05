package repository

import (
	"context"
	"github.com/nocturna-ta/election/internal/domain/model"
)

type ElectionRepository interface {
	InsertCandidate(ctx context.Context, candidate *model.Candidate, signedTransaction string) error
	GetAllCandidate(ctx context.Context) ([]model.Candidate, error)
	GetCandidateByNo(ctx context.Context, no string) (*model.Candidate, error)
	CandidateActivate(ctx context.Context, id string) error
}

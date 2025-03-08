package usecases

import (
	"context"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
)

type ElectionUseCases interface {
	RegisterCandidate(ctx context.Context, req *request.CandidateRegistrationRequest) (*response.CandidateResponse, error)
	GetAllCandidate(ctx context.Context) (*[]response.CandidateResponse, error)
	GetCandidateByNo(ctx context.Context, no string) (*response.CandidateResponse, error)
	ActivateCandidate(ctx context.Context, req *request.CandidateActivationRequest) (*response.CandidateActivation, error)
}

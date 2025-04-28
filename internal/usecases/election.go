package usecases

import (
	"context"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
)

type ElectionUseCases interface {
	RegisterElectionPair(ctx context.Context, req *request.ElectionPairRegistrationRequest) (*response.ElectionPairResponse, error)
	GetElectionPairByID(ctx context.Context, id string) (*response.ElectionPairResponse, error)
	GetElectionPairByNo(ctx context.Context, no string) (*response.ElectionPairResponse, error)
	GetAllElectionPairs(ctx context.Context) (*response.ElectionPairListResponse, error)
	ActivateElectionPair(ctx context.Context, req *request.ElectionPairActivationRequest) (*response.ElectionPairActivationResponse, error)
	UpsertElectionPairDetail(ctx context.Context, req *request.ElectionPairDetailRequest) (*response.ElectionPairDetailResponse, error)
	GetElectionPairDetail(ctx context.Context, pairID string) (*response.ElectionPairDetailResponse, error)
}

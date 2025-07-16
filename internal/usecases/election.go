package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
	"github.com/nocturna-ta/golib/http"
)

type ElectionUseCases interface {
	RegisterElectionPair(ctx context.Context, req *request.ElectionPairRegistrationRequest) (*response.ElectionPairResponse, error)
	UpdateElectionPair(ctx context.Context, req *request.ElectionPairUpdateRequest) (*response.ElectionPairResponse, error)
	GetElectionPairByID(ctx context.Context, id uuid.UUID) (*response.ElectionPairResponse, error)
	GetElectionPairByNo(ctx context.Context, no string) (*response.ElectionPairResponse, error)
	GetAllElectionPairs(ctx context.Context) (*[]response.ElectionPairFullResponse, error)
	ActivateElectionPair(ctx context.Context, req *request.ElectionPairActivationRequest) (*response.ElectionPairActivationResponse, error)
	UpsertElectionPairDetail(ctx context.Context, req *request.ElectionPairDetailRequest) (*response.ElectionPairDetailResponse, error)
	GetElectionPairDetail(ctx context.Context, pairID uuid.UUID) (*response.ElectionPairDetailResponse, error)
	GetElectionPairPhoto(ctx context.Context, id uuid.UUID) (*http.File, string, error)
	GetPresidentPhoto(ctx context.Context, id uuid.UUID) (*http.File, string, error)
	GetWorkProgramFile(ctx context.Context, id uuid.UUID) (*http.File, string, error)
	GetVicePresidentPhoto(ctx context.Context, id uuid.UUID) (*http.File, string, error)
	GetElectionPairFull(ctx context.Context, id uuid.UUID) (*response.ElectionPairFullResponse, error)
}

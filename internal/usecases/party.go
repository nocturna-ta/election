package usecases

import (
	"context"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
)

type PartyUseCases interface {
	RegisterParty(ctx context.Context, req *request.PartyRegisterRequest) (*response.PartyResponse, error)
	GetPartyByID(ctx context.Context, id string) (*response.PartyResponse, error)
	GetAllParties(ctx context.Context) ([]response.PartyResponse, error)
	UpdateParty(ctx context.Context, req *request.PartyUpdateRequest) (*response.PartyResponse, error)
	DeleteParty(ctx context.Context, id string) error
}

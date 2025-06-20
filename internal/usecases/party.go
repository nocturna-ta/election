package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
	"github.com/nocturna-ta/golib/http"
)

type PartyUseCases interface {
	RegisterParty(ctx context.Context, req *request.PartyRegisterRequest) (*response.PartyResponse, error)
	GetPartyByID(ctx context.Context, id string) (*response.PartyResponse, error)
	GetAllParties(ctx context.Context) ([]response.PartyResponse, error)
	UpdateParty(ctx context.Context, req *request.PartyUpdateRequest) (*response.PartyResponse, error)
	DeleteParty(ctx context.Context, id string) error
	GetPartyPhoto(ctx context.Context, id uuid.UUID) (*http.File, string, error)
}

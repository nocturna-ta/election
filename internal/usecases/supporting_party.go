package usecases

import (
	"context"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/election/internal/usecases/response"
)

type SupportingPartyUseCases interface {
	AddSupportingParty(ctx context.Context, req *request.AddSupportingPartyRequest) (*response.SupportingPartyResponse, error)
	GetSupportingPartiesByPairID(ctx context.Context, pairID string) (*response.SupportingPartiesResponse, error)
	RemoveSupportingParty(ctx context.Context, req *request.RemoveSupportingPartyRequest) error
}

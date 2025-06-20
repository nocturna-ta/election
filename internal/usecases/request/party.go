package request

import (
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"io"
)

type PartyRegisterRequest struct {
	Name     string    `json:"name" `
	LogoPath string    `json:"logo_path" `
	LogoFile io.Reader `json:"-" swaggerignore:"true"`
	LogoName string    `json:"-" swaggerignore:"true"`
}

type PartyUpdateRequest struct {
	ID       string    `json:"id" `
	Name     string    `json:"name" `
	LogoPath string    `json:"logo_path" `
	LogoFile io.Reader `json:"-" swaggerignore:"true"`
	LogoName string    `json:"-" swaggerignore:"true"`
}

func (r *PartyRegisterRequest) Validate() error {
	if r == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if r.Name == "" {
		return &custerr.ErrChain{
			Message: "Name cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

func (r *PartyUpdateRequest) Validate() error {
	if r == nil {
		return &custerr.ErrChain{
			Message: "Request cannot be nil",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	if r.Name == "" {
		return &custerr.ErrChain{
			Message: "Name cannot be empty",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	return nil
}

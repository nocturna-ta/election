package utils

import (
	"encoding/json"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"mime/multipart"
)

// ParseRegistrationRequest parses raw request body into ElectionPairRegistrationRequest
// This is in a separate file to avoid circular dependencies
func ParseRegistrationRequest(form *multipart.Form, files map[string]UploadedFile) (*request.ElectionPairRegistrationRequest, error) {

	pairValues := form.Value["pair"]
	if len(pairValues) == 0 {
		return nil, &custerr.ErrChain{
			Message: "Missing registration request data in 'pair' field",
			Code:    400,
			Type:    response.ErrBadRequest,
		}
	}

	pairJSON := pairValues[0]
	var regReq request.ElectionPairRegistrationRequest
	if err := json.Unmarshal([]byte(pairJSON), &regReq); err != nil {
		return nil, &custerr.ErrChain{
			Message: "Invalid JSON in registration request",
			Code:    400,
			Type:    response.ErrBadRequest,
			Cause:   err,
		}
	}

	if err := regReq.ValidateRegistrationRequest(); err != nil {
		return nil, err
	}

	if pairPhoto, exists := files["pair_photo"]; exists {
		regReq.PairPhotoName = pairPhoto.OriginalFilename
		regReq.PairPhotoFile = pairPhoto.File
	}

	if presidentPhoto, exists := files["president_photo"]; exists {
		regReq.President.PhotoPath = presidentPhoto.OriginalFilename
		regReq.President.PhotoFile = presidentPhoto.File
	}

	if vicePresidentPhoto, exists := files["vice_president_photo"]; exists {
		regReq.VicePresident.PhotoPath = vicePresidentPhoto.OriginalFilename
		regReq.VicePresident.PhotoFile = vicePresidentPhoto.File
	}

	return &regReq, nil
}

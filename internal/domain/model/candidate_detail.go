package model

import (
	"github.com/google/uuid"
	"github.com/nocturna-ta/common-model/models/event"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"time"
)

type CandidateDetail struct {
	BaseModel
	ID              uuid.UUID `db:"id"`
	CandidateID     uuid.UUID `db:"candidate_id"`
	Biodata         string    `db:"biodata"`
	Visi            string    `db:"visi"`
	Misi            string    `db:"misi"`
	ProgramKerja    string    `db:"program_kerja"`
	PartaiPendukung Partai    `db:"partai_pendukung"`
}

type Partai struct {
	ID         uuid.UUID `db:"id"`
	NamaPartai string    `db:"nama_partai"`
}

func (c *CandidateDetail) ToMessageModel() *event.ElectionDetailMessage {
	msg := &event.ElectionDetailMessage{
		BaseModelMessage: event.BaseModelMessage{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			IsDeleted: false,
		},
		ID:           c.ID.String(),
		CandidateID:  c.CandidateID.String(),
		Biodata:      c.Biodata,
		Visi:         c.Visi,
		Misi:         c.Misi,
		ProgramKerja: c.ProgramKerja,
	}

	return msg
}

func ConstructUpsertCandidateDetail(req *request.CandidateDetailRequest) *CandidateDetail {
	now := time.Now()
	candidateDetail := &CandidateDetail{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:              uuid.New(),
		CandidateID:     uuid.MustParse(req.CandidateID),
		Biodata:         req.Biodata,
		Visi:            req.Visi,
		Misi:            req.Misi,
		ProgramKerja:    req.ProgramKerja,
		PartaiPendukung: Partai{},
	}

	return candidateDetail
}

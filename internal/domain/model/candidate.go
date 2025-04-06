package model

import (
	"github.com/google/uuid"
	"github.com/nocturna-ta/common-model/models/event"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"strings"
	"time"
)

type Candidate struct {
	BaseModel
	ID              uuid.UUID `db:"id"`
	NameCandidate   []string  `db:"name"`
	ElectionNo      string    `db:"election_no"`
	VoteCount       int
	IsActive        bool            `db:"is_active"`
	CandidateDetail CandidateDetail `db:"candidate_detail"`
}

func (c *Candidate) ToMessageModel() *event.ElectionMessage {
	msg := &event.ElectionMessage{
		BaseModelMessage: event.BaseModelMessage{
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			IsDeleted: c.IsDeleted,
		},
		ID:            c.ID.String(),
		CandidateName: strings.Join(c.NameCandidate, ","),
		CandidateNo:   c.ElectionNo,
		VoteCount:     c.VoteCount,
		IsActive:      c.IsActive,
	}

	return msg
}
func ConstructRegistration(req *request.CandidateRegistrationRequest) *Candidate {
	now := time.Now()
	candidateId := uuid.New()
	candidate := &Candidate{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:            candidateId,
		NameCandidate: req.NameCandidate,
		ElectionNo:    req.ElectionNo,
		VoteCount:     0,
		IsActive:      false,
	}
	return candidate
}

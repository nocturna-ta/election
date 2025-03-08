package model

import "github.com/google/uuid"

type Candidate struct {
	BaseModel
	ID            uuid.UUID `db:"id"`
	NameCandidate []string  `db:"name"`
	ElectionNo    string    `db:"election_no"`
	VoteCount     int
	IsActive      bool `db:"is_active"`
}

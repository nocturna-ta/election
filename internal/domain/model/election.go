package model

import (
	"github.com/google/uuid"
	"github.com/nocturna-ta/election/internal/usecases/request"
	"time"
)

type ElectionPair struct {
	BaseModel
	ID                uuid.UUID         `db:"id"`
	ElectionNo        string            `db:"election_no"`
	VoteCount         int               `db:"vote_count"`
	IsActive          bool              `db:"is_active"`
	PairPhotoPath     string            `db:"pair_photo_path"`
	President         *CandidateInfo    `db:"president"`
	VicePresident     *CandidateInfo    `db:"vice_president"`
	SupportingParties []SupportingParty `db:"-"`
	PairDetail        *PairDetail       `db:"-"`
}

type PairDetail struct {
	BaseModel
	ID             uuid.UUID `db:"id"`
	ElectionPairID uuid.UUID `db:"election_pair_id"`
	Vision         string    `db:"vision"`
	Mission        string    `db:"mission"`
	WorkProgram    string    `db:"work_program"`
	ProgramDocs    []string  `db:"program_docs"`
}

type CandidateInfo struct {
	FullName           string `db:"full_name"`
	EducationHistory   string `db:"education_history"`
	WorkExperience     string `db:"work_experience"`
	LegalRecordHistory string `db:"legal_record_history"`
	PhotoPath          string `db:"photo_path"`
}

type ProgramDocument struct {
	BaseModel
	ID                   uuid.UUID `db:"id"`
	ElectionPairDetailID uuid.UUID `db:"election_pair_detail_id"`
	DocumentPath         string    `db:"document_path"`
	DocumentType         string    `db:"document_type"`
	OriginalFilename     string    `db:"original_filename"`
}

func ConstructElectionPair(req *request.ElectionPairRegistrationRequest) *ElectionPair {
	now := time.Now()
	pairID := uuid.New()

	president := &CandidateInfo{
		FullName:           req.President.FullName,
		EducationHistory:   req.President.EducationHistory,
		WorkExperience:     req.President.WorkExperience,
		LegalRecordHistory: req.President.LegalRecordHistory,
		PhotoPath:          req.President.PhotoPath,
	}

	vicePresident := &CandidateInfo{
		FullName:           req.VicePresident.FullName,
		EducationHistory:   req.VicePresident.EducationHistory,
		WorkExperience:     req.VicePresident.WorkExperience,
		LegalRecordHistory: req.VicePresident.LegalRecordHistory,
		PhotoPath:          req.VicePresident.PhotoPath,
	}

	pair := &ElectionPair{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:            pairID,
		ElectionNo:    req.ElectionNo,
		VoteCount:     0,
		IsActive:      false,
		PairPhotoPath: req.PairPhotoPath,
		President:     president,
		VicePresident: vicePresident,
	}

	return pair
}

func ConstructPairDetail(req *request.ElectionPairDetailRequest) *PairDetail {
	now := time.Now()
	detailID := uuid.New()

	detail := &PairDetail{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:             detailID,
		ElectionPairID: uuid.MustParse(req.ElectionPairID),
		Vision:         req.Vision,
		Mission:        req.Mission,
		WorkProgram:    req.WorkProgram,
		ProgramDocs:    req.ProgramDocs,
	}

	return detail
}

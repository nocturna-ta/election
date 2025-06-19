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
	PairName          string            `db:"pair_name"`
	President         *CandidateInfo    `db:"president"`
	VicePresident     *CandidateInfo    `db:"vice_president"`
	SupportingParties []SupportingParty `db:"-"`
	PairDetail        *PairDetail       `db:"-"`
}

type PairDetail struct {
	BaseModel
	ID             uuid.UUID     `db:"id"`
	ElectionPairID uuid.UUID     `db:"election_pair_id"`
	Vision         string        `db:"vision"`
	Mission        string        `db:"mission"`
	WorkProgram    []WorkProgram `db:"work_program"`
	ProgramDocs    string        `db:"work_program_docs"`
}

type CandidateInfo struct {
	FullName         string             `db:"full_name"`
	EducationHistory []EducationHistory `db:"education_history"`
	WorkExperience   []WorkHistory      `db:"work_experience"`
	Gender           string             `db:"gender"`
	BirthPlace       string             `db:"birth_place"`
	BirthDate        string             `db:"birth_date"`
	Religion         string             `db:"religion"`
	LastEducation    string             `db:"last_education"`
	Job              string             `db:"job"`
	PhotoPath        string             `db:"photo_path"`
}

func ConstructElectionPair(req *request.ElectionPairRegistrationRequest) *ElectionPair {
	now := time.Now()

	presidentEducationHistory := make([]EducationHistory, len(req.President.EducationHistory))
	for i, eh := range req.President.EducationHistory {
		presidentEducationHistory[i] = EducationHistory{
			InstituteName: eh.InstituteName,
			Year:          eh.Year,
		}
	}

	presidentWorkExperience := make([]WorkHistory, len(req.President.WorkExperience))
	for i, wh := range req.President.WorkExperience {
		presidentWorkExperience[i] = WorkHistory{
			InstituteName: wh.InstituteName,
			Position:      wh.Position,
			Year:          wh.Year,
		}
	}

	vicePresidentEducationHistory := make([]EducationHistory, len(req.VicePresident.EducationHistory))
	for i, eh := range req.VicePresident.EducationHistory {
		vicePresidentEducationHistory[i] = EducationHistory{
			InstituteName: eh.InstituteName,
			Year:          eh.Year,
		}
	}

	vicePresidentWorkExperience := make([]WorkHistory, len(req.VicePresident.WorkExperience))
	for i, wh := range req.VicePresident.WorkExperience {
		vicePresidentWorkExperience[i] = WorkHistory{
			InstituteName: wh.InstituteName,
			Position:      wh.Position,
			Year:          wh.Year,
		}
	}

	president := &CandidateInfo{
		FullName:         req.President.FullName,
		EducationHistory: presidentEducationHistory,
		WorkExperience:   presidentWorkExperience,
		Gender:           req.President.Gender,
		BirthPlace:       req.President.BirthPlace,
		BirthDate:        req.President.BirthDate,
		Religion:         req.President.Religion,
		LastEducation:    req.President.LastEducation,
		Job:              req.President.Job,
		PhotoPath:        req.President.PhotoPath,
	}

	vicePresident := &CandidateInfo{
		FullName:         req.VicePresident.FullName,
		EducationHistory: vicePresidentEducationHistory,
		WorkExperience:   vicePresidentWorkExperience,
		Gender:           req.VicePresident.Gender,
		BirthPlace:       req.VicePresident.BirthPlace,
		BirthDate:        req.VicePresident.BirthDate,
		Religion:         req.VicePresident.Religion,
		LastEducation:    req.VicePresident.LastEducation,
		Job:              req.VicePresident.Job,
		PhotoPath:        req.VicePresident.PhotoPath,
	}

	pair := &ElectionPair{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:            uuid.MustParse(req.ID),
		ElectionNo:    req.ElectionNo,
		VoteCount:     0,
		IsActive:      false,
		PairName:      req.PairName,
		PairPhotoPath: req.PairPhotoPath,
		President:     president,
		VicePresident: vicePresident,
	}

	return pair
}

func ConstructPairDetail(req *request.ElectionPairDetailRequest) *PairDetail {
	now := time.Now()
	detailID := uuid.New()

	workProgram := make([]WorkProgram, len(req.WorkProgram))
	for i, wp := range req.WorkProgram {
		workProgram[i] = WorkProgram{
			ProgramName:  wp.ProgramName,
			ProgramPhoto: wp.ProgramPhoto,
			ProgramDesc:  wp.ProgramDesc,
		}
	}

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
		WorkProgram:    workProgram,
		ProgramDocs:    req.ProgramDocs,
	}

	return detail
}

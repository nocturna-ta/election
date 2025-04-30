package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type ElectionPairDTO struct {
	ID            uuid.UUID `db:"id"`
	ElectionNo    string    `db:"election_no"`
	VoteCount     int       `db:"vote_count"`
	IsActive      bool      `db:"is_active"`
	PairPhotoPath string    `db:"pair_photo_path"`

	PresidentFullName         string          `db:"president_full_name"`
	PresidentEducationHistory json.RawMessage `db:"president_education_history"`
	PresidentWorkExperience   json.RawMessage `db:"president_work_experience"`
	PresidentGender           string          `db:"president_gender"`
	PresidentBirthPlace       string          `db:"president_birth_place"`
	PresidentBirthDate        string          `db:"president_birth_date"`
	PresidentReligion         string          `db:"president_religion"`
	PresidentLastEducation    string          `db:"president_last_education"`
	PresidentJob              string          `db:"president_job"`
	PresidentPhotoPath        string          `db:"president_photo_path"`

	VicePresidentFullName         string          `db:"vice_president_full_name"`
	VicePresidentEducationHistory json.RawMessage `db:"vice_president_education_history"`
	VicePresidentWorkExperience   json.RawMessage `db:"vice_president_work_experience"`
	VicePresidentGender           string          `db:"vice_president_gender"`
	VicePresidentBirthPlace       string          `db:"vice_president_birth_place"`
	VicePresidentBirthDate        string          `db:"vice_president_birth_date"`
	VicePresidentReligion         string          `db:"vice_president_religion"`
	VicePresidentLastEducation    string          `db:"vice_president_last_education"`
	VicePresidentJob              string          `db:"vice_president_job"`
	VicePresidentPhotoPath        string          `db:"vice_president_photo_path"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	IsDeleted bool      `db:"is_deleted"`
}

func (dto *ElectionPairDTO) ToDomain() (*ElectionPair, error) {

	var (
		presidentEducationHistory     []EducationHistory
		presidentWorkExperience       []WorkHistory
		vicePresidentEducationHistory []EducationHistory
		vicePresidentWorkExperience   []WorkHistory
	)

	if len(dto.PresidentEducationHistory) > 0 {
		if err := json.Unmarshal(dto.PresidentEducationHistory, &presidentEducationHistory); err != nil {
			return nil, err
		}
	}

	if len(dto.PresidentWorkExperience) > 0 {
		if err := json.Unmarshal(dto.PresidentWorkExperience, &presidentWorkExperience); err != nil {
			return nil, err
		}
	}

	if len(dto.VicePresidentEducationHistory) > 0 {
		if err := json.Unmarshal(dto.VicePresidentEducationHistory, &vicePresidentEducationHistory); err != nil {
			return nil, err
		}
	}

	if len(dto.VicePresidentWorkExperience) > 0 {
		if err := json.Unmarshal(dto.VicePresidentWorkExperience, &vicePresidentWorkExperience); err != nil {
			return nil, err
		}
	}

	return &ElectionPair{
		BaseModel: BaseModel{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
			IsDeleted: dto.IsDeleted,
		},
		ID:            dto.ID,
		ElectionNo:    dto.ElectionNo,
		VoteCount:     dto.VoteCount,
		IsActive:      dto.IsActive,
		PairPhotoPath: dto.PairPhotoPath,
		President: &CandidateInfo{
			FullName:         dto.PresidentFullName,
			EducationHistory: presidentEducationHistory,
			WorkExperience:   presidentWorkExperience,
			Gender:           dto.PresidentGender,
			BirthPlace:       dto.PresidentBirthPlace,
			BirthDate:        dto.PresidentBirthDate,
			Religion:         dto.PresidentReligion,
			LastEducation:    dto.PresidentLastEducation,
			Job:              dto.PresidentJob,
			PhotoPath:        dto.PresidentPhotoPath,
		},
		VicePresident: &CandidateInfo{
			FullName:         dto.VicePresidentFullName,
			EducationHistory: vicePresidentEducationHistory,
			WorkExperience:   vicePresidentWorkExperience,
			Gender:           dto.VicePresidentGender,
			BirthPlace:       dto.VicePresidentBirthPlace,
			BirthDate:        dto.VicePresidentBirthDate,
			Religion:         dto.VicePresidentReligion,
			LastEducation:    dto.VicePresidentLastEducation,
			Job:              dto.VicePresidentJob,
			PhotoPath:        dto.VicePresidentPhotoPath,
		},
	}, nil
}

package model

import (
	"github.com/nocturna-ta/common-model/models/event"
	"time"
)

func (e *ElectionPair) ToMessageModel(txHash string) *event.ElectionPairMessage {
	msg := &event.ElectionPairMessage{
		BaseModelMessage: event.BaseModelMessage{
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
			IsDeleted: e.IsDeleted,
		},
		ID:              e.ID.String(),
		ElectionNo:      e.ElectionNo,
		VoteCount:       e.VoteCount,
		IsActive:        e.IsActive,
		PairPhotoPath:   e.PairPhotoPath,
		TransactionHash: txHash,
	}

	if e.President != nil {
		presidentEduHistory := make([]event.EducationHistoryMessage, len(e.President.EducationHistory))
		for i, edu := range e.President.EducationHistory {
			presidentEduHistory[i] = event.EducationHistoryMessage{
				InstituteName: edu.InstituteName,
				Year:          edu.Year,
			}
		}

		presidentWorkHistory := make([]event.WorkHistoryMessage, len(e.President.WorkExperience))
		for i, work := range e.President.WorkExperience {
			presidentWorkHistory[i] = event.WorkHistoryMessage{
				InstituteName: work.InstituteName,
				Position:      work.Position,
				Year:          work.Year,
			}
		}

		msg.President = event.CandidateInfoMessage{
			FullName:         e.President.FullName,
			EducationHistory: presidentEduHistory,
			WorkExperience:   presidentWorkHistory,
			Gender:           e.President.Gender,
			BirthPlace:       e.President.BirthPlace,
			BirthDate:        e.President.BirthDate,
			Religion:         e.President.Religion,
			LastEducation:    e.President.LastEducation,
			Job:              e.President.Job,
			PhotoPath:        e.President.PhotoPath,
		}
	}

	if e.VicePresident != nil {
		vicePresidentEduHistory := make([]event.EducationHistoryMessage, len(e.VicePresident.EducationHistory))
		for i, edu := range e.VicePresident.EducationHistory {
			vicePresidentEduHistory[i] = event.EducationHistoryMessage{
				InstituteName: edu.InstituteName,
				Year:          edu.Year,
			}
		}

		vicePresidentWorkHistory := make([]event.WorkHistoryMessage, len(e.VicePresident.WorkExperience))
		for i, work := range e.VicePresident.WorkExperience {
			vicePresidentWorkHistory[i] = event.WorkHistoryMessage{
				InstituteName: work.InstituteName,
				Position:      work.Position,
				Year:          work.Year,
			}
		}

		msg.VicePresident = event.CandidateInfoMessage{
			FullName:         e.VicePresident.FullName,
			EducationHistory: vicePresidentEduHistory,
			WorkExperience:   vicePresidentWorkHistory,
			Gender:           e.VicePresident.Gender,
			BirthPlace:       e.VicePresident.BirthPlace,
			BirthDate:        e.VicePresident.BirthDate,
			Religion:         e.VicePresident.Religion,
			LastEducation:    e.VicePresident.LastEducation,
			Job:              e.VicePresident.Job,
			PhotoPath:        e.VicePresident.PhotoPath,
		}
	}

	if e.PairDetail != nil {
		msg.Detail = ConvertPairDetailToMessage(e.PairDetail)
	}

	return msg
}

func ConvertPairDetailToMessage(detail *PairDetail) *event.ElectionPairDetailMessage {
	if detail == nil {
		return nil
	}

	workPrograms := make([]event.WorkProgramMessage, len(detail.WorkProgram))
	for i, program := range detail.WorkProgram {
		workPrograms[i] = event.WorkProgramMessage{
			ProgramName:  program.ProgramName,
			ProgramPhoto: program.ProgramPhoto,
			ProgramDesc:  program.ProgramDesc,
		}
	}

	return &event.ElectionPairDetailMessage{
		ID:             detail.ID.String(),
		ElectionPairID: detail.ElectionPairID.String(),
		Vision:         detail.Vision,
		Mission:        detail.Mission,
		WorkProgram:    workPrograms,
		ProgramDocs:    detail.ProgramDocs,
	}
}

func CreateElectionActivationMessage(electionPairID string, isActive bool, transactionHash string) *event.ElectionActivationMessage {
	return &event.ElectionActivationMessage{
		BaseModelMessage: event.BaseModelMessage{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
		},
		ElectionPairID:  electionPairID,
		IsActive:        isActive,
		ActivatedAt:     time.Now(),
		TransactionHash: transactionHash,
	}
}

func (p *Party) ToMessageModel() *event.PartyMessage {
	return &event.PartyMessage{
		BaseModelMessage: event.BaseModelMessage{
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			IsDeleted: p.IsDeleted,
		},
		ID:       p.ID.String(),
		Name:     p.Name,
		LogoPath: p.LogoPath,
	}
}

func (s *SupportingParty) ToMessageModel() *event.SupportingPartyMessage {
	msg := &event.SupportingPartyMessage{
		BaseModelMessage: event.BaseModelMessage{
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			IsDeleted: s.IsDeleted,
		},
		ID:             s.ID.String(),
		ElectionPairID: s.ElectionPairID.String(),
		PartyID:        s.PartyID.String(),
	}

	if s.Party != nil {
		msg.PartyName = s.Party.Name
		msg.PartyLogoPath = s.Party.LogoPath
	}
	return msg
}

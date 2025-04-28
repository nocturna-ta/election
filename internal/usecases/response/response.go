package response

type CandidateInfoResponse struct {
	FullName           string `json:"full_name"`
	EducationHistory   string `json:"education_history"`
	WorkExperience     string `json:"work_experience"`
	LegalRecordHistory string `json:"legal_record_history"`
	PhotoPath          string `json:"photo_path"`
}

type ElectionPairResponse struct {
	ID            string                `json:"id"`
	ElectionNo    string                `json:"election_no"`
	VoteCount     int                   `json:"vote_count"`
	IsActive      bool                  `json:"is_active"`
	PairPhotoPath string                `json:"pair_photo_path"`
	President     CandidateInfoResponse `json:"president"`
	VicePresident CandidateInfoResponse `json:"vice_president"`
}

type ElectionPairDetailResponse struct {
	ID             string   `json:"id"`
	ElectionPairID string   `json:"election_pair_id"`
	Vision         string   `json:"vision"`
	Mission        string   `json:"mission"`
	WorkProgram    string   `json:"work_program"`
	ProgramDocs    []string `json:"program_docs,omitempty"`
}

type ElectionPairActivationResponse struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
}

type ElectionPairListResponse struct {
	Pairs []ElectionPairResponse `json:"pairs"`
	Total int                    `json:"total"`
}

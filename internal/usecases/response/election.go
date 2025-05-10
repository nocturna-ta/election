package response

type CandidateInfoResponse struct {
	FullName         string                     `json:"full_name"`
	EducationHistory []EducationHistoryResponse `json:"education_history"`
	WorkExperience   []WorkHistoryResponse      `json:"work_experience"`
	Gender           string                     `json:"gender"`
	BirthPlace       string                     `json:"birth_place"`
	BirthDate        string                     `json:"birth_date"`
	Religion         string                     `json:"religion"`
	LastEducation    string                     `json:"last_education"`
	Job              string                     `json:"job"`
	PhotoPath        string                     `json:"photo_path"`
}

type EducationHistoryResponse struct {
	InstituteName string `json:"institute_name"`
	Year          string `json:"year"`
}

type WorkHistoryResponse struct {
	InstituteName string `json:"institute_name"`
	Position      string `json:"position"`
	Year          string `json:"year"`
}

type WorkProgramResponse struct {
	ProgramName  string   `json:"program_name"`
	ProgramPhoto string   `json:"program_photo"`
	ProgramDesc  []string `json:"program_desc"`
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
	ID             string                `json:"id"`
	ElectionPairID string                `json:"election_pair_id"`
	Vision         string                `json:"vision"`
	Mission        string                `json:"mission"`
	WorkProgram    []WorkProgramResponse `json:"work_program"`
	ProgramDocs    string                `json:"program_docs"`
}

type ElectionPairActivationResponse struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
}

type ElectionPairListResponse struct {
	Pairs []ElectionPairResponse `json:"pairs"`
	Total int                    `json:"total"`
}

type ElectionPairFullResponse struct {
	ElectionPairResponse
	Detail            ElectionPairDetailResponse `json:"detail,omitempty"`
	SupportingParties []SupportingPartyResponse  `json:"supporting_parties,omitempty"`
}

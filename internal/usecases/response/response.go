package response

type CandidateResponse struct {
	ID            string   `json:"id"`
	NameCandidate []string `json:"name_candidate"`
	ElectionNo    string   `json:"election_no"`
	VoteCount     int      `json:"vote_count"`
	IsActive      bool     `json:"is_active"`
}

type CandidateActivation struct {
	IsActive bool `json:"is_active"`
}

type CandidateDetailResponse struct {
	ID           string `json:"id"`
	CandidateID  string `json:"candidate_id"`
	Biodata      string `json:"biodata"`
	Visi         string `json:"visi"`
	Misi         string `json:"misi"`
	ProgramKerja string `json:"program_kerja"`
}

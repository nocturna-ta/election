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

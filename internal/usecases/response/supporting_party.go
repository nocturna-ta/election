package response

type SupportingPartyResponse struct {
	ID             string        `json:"id"`
	ElectionPairID string        `json:"election_pair_id"`
	PartyID        string        `json:"party_id"`
	Party          PartyResponse `json:"party"`
}

type SupportingPartiesResponse struct {
	Parties []SupportingPartyResponse `json:"parties"`
	Total   int                       `json:"total"`
}

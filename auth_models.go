package tesla

type MFAVerify struct {
	Data MFAData `json:"data"`
}

type MFAData struct {
	Id          string `json:"id"`
	ChallengeId string `json:"challengeId"`
	FactorId    string `json:"factorId"`
	PassCode    string `json:"passCode"`
	Approved    bool   `json:"approved"`
	Flagged     bool   `json:"flagged"`
	Valid       bool   `json:"valid"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

package api

// the response object of allContestsURL
type contests struct {
	Adventures []adventure `json:"adventures"`
}

type adventure struct {
	Title string `json:"title"`
	ID    string `json:"id"`

	Competition    int    `json:"competition"`
	CompetitionStr string `json:"competition__str"`

	RegStartTime int64 `json:"reg_start_ms"`
	RegEndTime   int64 `json:"reg_end_ms"`

	Challenges []challenge `json:"challenges"`
}

package main

type apiResponse struct {
	ScoreboardSize int `json:"full_scoreboard_size"`

	Challenge  challenge   `json:"challenge"`
	UserScores []userScore `json:"user_scores"`
}

type test struct {
	Point int `json:"value"`

	Solved int `json:"num_solved"`

	Status       int    `json:"type"`
	StatusString string `json:"type__str"` // VISIBLE / HIDDEN
}

type tasks []task

func (t tasks) Len() int {
	return len(t)
}
func (t tasks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t tasks) Less(i, j int) bool {
	// sort by score
	// tie break by title

	a := 0
	for _, p := range t[i].Tests {
		a += p.Point
	}
	b := 0
	for _, p := range t[j].Tests {
		b += p.Point
	}

	if a == b {
		return t[i].Title < t[j].Title
	}
	return a < b
}

type task struct {
	TaskID string `json:"id"`
	Title  string `json:"title"`

	Attempts int `json:"num_attempted"`

	Tests []test `json:"tests"`
}

type challenge struct {
	ContestID      string `json:"id"`
	Title          string `json:"title"`
	AdditionalInfo string `json:"additional_info"`

	Tasks tasks `json:"tasks"`

	ResultStatus       int    `json:"result_status"`
	ResultStatusString string `json:"result_status__str"` // PARTIALLY_HIDDEN

	StartMS int64 `json:"start_ms"`
	EndMS   int64 `json:"end_ms"`

	AreResultsFinal     bool `json:"are_results_final"`
	IsPracticeAvailable bool `json:"is_practice_available"`
}

type taskInfo struct {
	TaskID string `json:"task_id"`
	Point  int    `json:"score"`

	Attempts int `json:"total_attempts"`

	AC            int `json:"tests_definitely_solved"`
	PretestPassed int `json:"tests_possibly_solved"`

	WA   int   `json:"penalty_attempts"`
	WAms int64 `json:"penalty_micros"`
}

type userScores []userScore

func (t userScores) Len() int {
	return len(t)
}
func (t userScores) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t userScores) Less(i, j int) bool {
	// sort by rank
	return t[i].Rank < t[j].Rank
}

type userScore struct {
	Handle  string `json:"displayname"`
	Country string `json:"country"`
	Rank    int    `json:"rank"`

	Score  int   `json:"score_1"`
	Score2 int64 `json:"score_2"`

	TasksInfo []taskInfo `json:"task_info"`
}

package api

import (
	"encoding/json"
	"log"
	"sort"
	"sync"
)

type apiResponse struct {
	ScoreboardSize int `json:"full_scoreboard_size"`

	Challenge  challenge   `json:"challenge"`
	UserScores []userScore `json:"user_scores"`
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

type scoreboardPaginationPayload struct {
	// {"min_rank":11,"num_consecutive_users":10}
	StartingRank       int `json:"min_rank"`
	ConsecutiveRecords int `json:"num_consecutive_users"`
}

func (data *ContestMetadata) getScoreboardPaginationPayload(startingRank, consecutiveRecords int) string {
	payload := scoreboardPaginationPayload{StartingRank: startingRank, ConsecutiveRecords: consecutiveRecords}
	res, err := json.Marshal(payload)
	handleErr(err)
	return encodeToBase64(res)
}

func (data *ContestMetadata) fetchAllContestantData(country string) {
	var wg sync.WaitGroup
	for i := 1; i <= data.TotalContestants; i += data.StepSize {
		go func(starting, step int) {
			wg.Add(1)
			log.Println("Starting from record", starting)

			param := make([]interface{}, 2)
			param[0] = starting
			param[1] = step
			response := data.fetchAPIResponseBody(scoreboardType, param)

			data.Lock()
			data.UserScores = append(data.UserScores, response.UserScores...)
			data.Unlock()

			log.Println("Done", starting)
			wg.Done()
		}(i, data.StepSize)
	}
	wg.Wait()

	// sort by rank
	sort.Sort(data.UserScores)
}

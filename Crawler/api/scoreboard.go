package api

import (
	"encoding/json"
	"sort"
	"sync"
)

type scoreboardResponse struct {
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
	// tie break by name
	if t[i].Rank == t[j].Rank {
		return t[i].Handle < t[j].Handle
	}
	return t[i].Rank < t[j].Rank
}

type userScore struct {
	Handle  string `json:"displayname"`
	Country string `json:"country"`
	Rank    int    `json:"rank"`

	Score  int   `json:"score_1"`
	Score2 int64 `json:"score_2"`

	TasksInfo []taskInfo `json:"task_info"`

	isEmpty bool
}

type scoreboardPaginationPayload struct {
	// {"min_rank":11,"num_consecutive_users":10}
	StartingRank       int `json:"min_rank"`
	ConsecutiveRecords int `json:"num_consecutive_users"`
}

func getScoreboardPaginationPayload(startingRank, consecutiveRecords int) string {
	payload := scoreboardPaginationPayload{StartingRank: startingRank, ConsecutiveRecords: consecutiveRecords}
	res, err := json.Marshal(payload)
	handleErr(err)
	return encodeToBase64(res)
}

func (data *ContestData) fetchScoreboard(startingRank int, contestantChannel chan userScores, pool chan bool) {
	// log.Println("Starting from rank", startingRank)

	param := make([]interface{}, 2)
	param[0] = startingRank
	param[1] = data.stepSize
	response := fetchAPIResponse(scoreboardType, data.contestID, param).(*scoreboardResponse)

	contestantChannel <- response.UserScores

	// log.Println("Done", startingRank)
	<-pool
}

func (data *ContestData) mergeContestants(contestantChannel chan userScores, wg *sync.WaitGroup) {
	for {
		data.contestants = append(data.contestants, <-contestantChannel...)
		wg.Done()
	}
}

func (data *ContestData) fetchAllContestantData(concurrentFetch int) {
	contestantChannel := make(chan userScores, concurrentFetch)
	pool := make(chan bool, concurrentFetch)

	var wg sync.WaitGroup
	go data.mergeContestants(contestantChannel, &wg)

	for i := 1; i <= data.totalContestants; i += data.stepSize {
		pool <- true
		wg.Add(1)
		go data.fetchScoreboard(i, contestantChannel, pool)
	}

	wg.Wait()

	// sort users by rank
	sort.Sort(data.contestants)
}

package api

import (
	"encoding/json"
	"log"
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

func getScoreboardPaginationPayload(startingRank, consecutiveRecords int) string {
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
			response := fetchAPIResponse(scoreboardType, data.ContestID, param).(*scoreboardResponse)

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

func (data *ContestData) fetchScoreboard(startingRank int, contestantChannel chan userScores, pool chan bool, wg *sync.WaitGroup) {
	log.Println("Starting from rank", startingRank)

	param := make([]interface{}, 2)
	param[0] = startingRank
	param[1] = data.stepSize
	response := fetchAPIResponse(scoreboardType, data.contestID, param).(*scoreboardResponse)

	contestantChannel <- response.UserScores

	log.Println("Done", startingRank)
	<-pool
	wg.Done()
}

func (data *ContestData) mergeContestants(contestantChannel chan userScores) {
	for {
		tmp, more := <-contestantChannel
		if more == false {
			return
		}
		data.contestants = append(data.contestants, tmp...)
	}
}

func (data *ContestData) fetchAllContestantData(concurrentFetch int) {
	contestantChannel := make(chan userScores, concurrentFetch)
	pool := make(chan bool, concurrentFetch)
	go data.mergeContestants(contestantChannel)

	var wg sync.WaitGroup
	for i := 1; i <= data.totalContestants; i += data.stepSize {
		pool <- true
		wg.Add(1)
		go data.fetchScoreboard(i, contestantChannel, pool, &wg)
	}

	wg.Wait()
	close(pool)
	close(contestantChannel)

	// sort users by rank
	sort.Sort(data.contestants)
}

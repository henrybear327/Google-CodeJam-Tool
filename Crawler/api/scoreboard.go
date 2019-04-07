package api

import (
	"encoding/json"
	"log"
	"sort"
	"sync"
)

const scoreboardPaginationURL = "https://codejam.googleapis.com/scoreboard/%s/poll?p=%s"

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
			response := data.fetchResponseBody(2, param)

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

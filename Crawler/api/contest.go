package api

import "sort"

func (data *ContestMetadata) FetchContestInfo() {
	data.Lock()
	defer data.Unlock()

	// init
	data.UserScores = nil

	// fetch contest data
	param := make([]interface{}, 2)
	param[0] = 1
	param[1] = 10
	result := data.fetchResponseBody(2, param)

	// set scoreboard
	data.TotalContestants = result.ScoreboardSize

	// set problem order
	data.ContestInfo = result.Challenge
	sort.Sort(data.ContestInfo.Tasks)

	// set step size
	if data.StepSize == 0 {
		data.StepSize = 100
	}
}

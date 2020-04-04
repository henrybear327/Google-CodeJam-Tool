package api

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

// ContestData is the object for holding data of a specific contest
type ContestData struct {
	contestID string

	lastChallengeFetchTime  int64
	lastScoreboardFetchTime int64

	totalContestants int
	stepSize         int

	challenge   challenge
	contestants userScores

	sync.Mutex
}

// the response object of specificContestURL
type contestResponse struct {
	Challenge challenge  `json:"challenge"`
	Languages []language `json:"languages"`
}

type language struct {
	ID         int    `json:"id"`
	IDStr      string `json:"id__str"`
	Name       string `json:"name"`
	HelperText string `json:"helper_text"`
}

type test struct {
	Point int `json:"value"`

	Solved int `json:"num_solved"`

	Status       int    `json:"type"`      // 1 / 2
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
	ID            string `json:"id"`
	Title         string `json:"title"`
	StatementHTML string `json:"statement"`
	AnalysisHTML  string `json:"analysis"`

	Attempts    int    `json:"num_attempted"`
	Tests       []test `json:"tests"`
	totalPoints int

	InputType       int    `json:"trial_input_type"`
	InputTypeString string `json:"trial_input_type__str"`
}

type challenge struct {
	ContestID      string `json:"id"`
	Title          string `json:"title"`
	AdditionalInfo string `json:"additional_info"`
	RecapHTML      string `json:"recap"`

	Tasks tasks `json:"tasks"`

	ResultStatus       int    `json:"result_status"`
	ResultStatusString string `json:"result_status__str"` // PARTIALLY_HIDDEN

	StartTime int64 `json:"start_ms"`
	EndTime   int64 `json:"end_ms"`

	AreResultsFinal     bool `json:"are_results_final"`
	IsPracticeAvailable bool `json:"is_practice_available"`

	MyUserType    int    `json:"my_user_type"`
	MyUserTypeStr string `json:"my_user_type__str"`

	ShowScoreboard bool `json:"show_scoreboard"`
}

// FetchContest gathers the information for a specific contest
// This is the main object for a contest
func (data *ContestData) FetchContest(ID string, concurrentFetch int, forceFetch bool) error {
	data.Lock()
	defer data.Unlock()

	data.contestID = ID

	// get problem sets
	if forceFetch || len(data.challenge.ContestID) == 0 {
		data.lastChallengeFetchTime = time.Now().UnixNano()

		log.Println("New contest object, fetching challenges")
		param := make([]interface{}, 1)
		param[0] = ID
		response := fetchAPIResponse(specificContestType, data.contestID, param).(*contestResponse)

		data.challenge = response.Challenge
		sort.Sort(data.challenge.Tasks) // reorder the tasks

		for _, task := range data.challenge.Tasks {
			task.totalPoints = 0
			for _, subtask := range task.Tests {
				task.totalPoints += subtask.Point
			}
		}

		fmt.Println("Is the result finalized?", data.challenge.AreResultsFinal)
		log.Println("Done!")
	}

	// get scoreboard
	{
		data.lastScoreboardFetchTime = time.Now().UnixNano()

		log.Println("Fetching entire scoreboard")
		// get scoreboard metadata
		data.stepSize = 100

		param := make([]interface{}, 2)
		param[0] = 1
		param[1] = 10
		scoreboardResponse := fetchAPIResponse(scoreboardType, data.contestID, param).(*scoreboardResponse)
		data.totalContestants = scoreboardResponse.ScoreboardSize

		// fetch the entire scoreboard
		data.fetchAllContestantData(concurrentFetch)

		log.Println("Done!")
	}

	return nil
}

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
)

type ContestMetadata struct {
	ContestID string

	TotalContestants int
	StepSize         int

	ContestInfo challenge
	UserScores  userScores

	sync.Mutex
}

func fetchResponse(url string) []byte {
	// log.Println("url", url)
	resp, err := http.Get(url)
	handleErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	// log.Println(string(body))

	result := decodeFromBase64(body)
	// log.Println(string(result))

	return result
}

func (data *ContestMetadata) fetchResponseBody(fetchType int, param []interface{}) *apiResponse {
	url := ""
	if fetchType == 1 { // handle search
		handle := param[0].(string)
		url = fmt.Sprintf(findHandleURL, data.ContestID, data.getHandleSearchPayload(handle))
	} else if fetchType == 2 { // dump scoreboard
		starting := param[0].(int)
		step := param[1].(int)
		url = fmt.Sprintf(scoreboardPaginationURL, data.ContestID, data.getScoreboardPaginationPayload(starting, step))
	} else if fetchType == 3 {
		url = param[0].(string)
	} else {
		log.Fatalln("Unknown option")
	}

	result := fetchResponse(url)

	var response apiResponse
	err := json.Unmarshal(result, &response)
	handleErr(err)
	// log.Println(response)

	return &response
}

func (data *ContestMetadata) printUserRecord(user *userScore) {
	fmt.Printf("+====== %-15s (%s) ======+\n", user.Handle, user.Country)
	fmt.Printf("Rank %v Score %v\n\n", user.Rank, user.Score)

	for _, curTask := range data.ContestInfo.Tasks {
		fmt.Printf("%-30s ", curTask.Title)
		totalPoint := 0
		for _, p := range curTask.Tests {
			totalPoint += p.Point
		}

		for _, task := range user.TasksInfo {
			if curTask.TaskID == task.TaskID {
				fmt.Printf("%2d / %2d\n", task.Point, totalPoint)
				goto found
			}
		}

		// not attempted, so no record is present
		fmt.Printf("%2d / %2d\n", 0, totalPoint)
	found:
	}
}

func (data *ContestMetadata) GetHandleResults(handles []string) {
	results := make(userScores, len(handles))
	ch := make(chan userScore)
	for _, handle := range handles {
		go data.fetchHandleResult(handle, ch)
	}

	for i := 0; i < len(handles); i++ {
		results[i] = <-ch
	}

	sort.Sort(results)
	for _, user := range results {
		data.printUserRecord(&user)
	}
}

func (data *ContestMetadata) GetAllContestantData(country string) {
	data.fetchAllContestantData(country)

	for _, contestant := range data.UserScores {
		if len(country) > 0 && contestant.Country == country {
			data.printUserRecord(&contestant)
		} else if len(country) == 0 {
			data.printUserRecord(&contestant)
		}
	}
}

func (data *ContestMetadata) GetJSONResponse(url string) {
	response := fetchResponse(url)
	fmt.Println(string(response))
}

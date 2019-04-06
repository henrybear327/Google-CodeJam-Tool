package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
)

type contestMetadata struct {
	contestID string

	totalContestants int
	stepSize         int

	contestInfo challenge
	userScores  userScores

	sync.Mutex
}

const findHandleURL = "https://codejam.googleapis.com/scoreboard/%s/find?p=%s"
const scoreboardPaginationURL = "https://codejam.googleapis.com/scoreboard/%s/poll?p=%s"

type handleSearchPayload struct {
	/*
		{"nickname":"henrybear327","scoreboard_page_size":10}
	*/
	Handle   string `json:"nickname"`
	PageSize int    `json:"scoreboard_page_size"`
}

func (data *contestMetadata) getHandleSearchPayload(handle string) string {
	payload := handleSearchPayload{Handle: handle, PageSize: 1}
	res, err := json.Marshal(payload)
	handleErr(err)
	return encodeToBase64(res)
}

type scoreboardPaginationPayload struct {
	// {"min_rank":11,"num_consecutive_users":10}
	StartingRank       int `json:"min_rank"`
	ConsecutiveRecords int `json:"num_consecutive_users"`
}

func (data *contestMetadata) getScoreboardPaginationPayload(startingRank, consecutiveRecords int) string {
	payload := scoreboardPaginationPayload{StartingRank: startingRank, ConsecutiveRecords: consecutiveRecords}
	res, err := json.Marshal(payload)
	handleErr(err)
	return encodeToBase64(res)
}

func (data *contestMetadata) fetchResponseBody(fetchType int, param []interface{}) *apiResponse {
	url := ""
	if fetchType == 1 { // handle search
		handle := param[0].(string)
		url = fmt.Sprintf(findHandleURL, data.contestID, data.getHandleSearchPayload(handle))
	} else if fetchType == 2 { // dump scoreboard
		starting := param[0].(int)
		step := param[1].(int)
		url = fmt.Sprintf(scoreboardPaginationURL, data.contestID, data.getScoreboardPaginationPayload(starting, step))
	} else {
		log.Fatalln("Unknown option")
	}

	// log.Println("url", url)
	resp, err := http.Get(url)
	handleErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	// log.Println(string(body))

	result := decodeFromBase64(body)
	// log.Println(string(result))

	var response apiResponse
	err = json.Unmarshal(result, &response)
	handleErr(err)
	// log.Println(response)

	return &response
}

func (data *contestMetadata) fetchContestInfo() {
	data.Lock()
	defer data.Unlock()

	// init
	data.userScores = nil

	// fetch contest data
	param := make([]interface{}, 2)
	param[0] = 1
	param[1] = 10
	result := data.fetchResponseBody(2, param)

	// set scoreboard
	data.totalContestants = result.ScoreboardSize

	// set problem order
	data.contestInfo = result.Challenge
	sort.Sort(data.contestInfo.Tasks)

	// set step size
	if data.stepSize == 0 {
		data.stepSize = 100
	}
}

func (data contestMetadata) printUserRecord(user *userScore) {
	fmt.Printf("+====== %15s (%s) ======+\n", user.Handle, user.Country)
	fmt.Printf("Rank %v Score %v\n\n", user.Rank, user.Score)

	for _, curTask := range data.contestInfo.Tasks {
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

func (data *contestMetadata) prettyPrint(response *apiResponse, isSingle bool) {
	data.Lock()
	defer data.Unlock()

	if isSingle {
		if len(response.UserScores) != 1 {
			log.Fatalln("Incorrect user count", len(response.UserScores))
		}
		user := response.UserScores[0]
		data.printUserRecord(&user)
	} else {

	}
}

func (data *contestMetadata) GetHandleResults(handles []string) {
	var wg sync.WaitGroup
	for _, handle := range handles {
		wg.Add(1)
		go data.fetchHandleResult(handle, &wg)
	}
	wg.Wait()
}

func (data *contestMetadata) fetchHandleResult(handle string, wg *sync.WaitGroup) {
	param := make([]interface{}, 1)
	param[0] = handle
	response := data.fetchResponseBody(1, param)

	data.prettyPrint(response, true)
	wg.Done()
}

func (data *contestMetadata) GetAllContestantData(country string) {
	data.fetchAllContestantData(country)

	for _, contestant := range data.userScores {
		if len(country) > 0 && contestant.Country == country {
			data.printUserRecord(&contestant)
		}
	}
}

func (data *contestMetadata) fetchAllContestantData(country string) {
	var wg sync.WaitGroup
	for i := 1; i <= data.totalContestants; i += data.stepSize {
		go func(starting, step int) {
			wg.Add(1)
			log.Println("Starting from record", starting)

			param := make([]interface{}, 2)
			param[0] = starting
			param[1] = step
			response := data.fetchResponseBody(2, param)

			data.Lock()
			data.userScores = append(data.userScores, response.UserScores...)
			data.Unlock()

			data.prettyPrint(response, false)

			log.Println("Done", starting)
			wg.Done()
		}(i, data.stepSize)
	}
	wg.Wait()

	// sort by rank
	sort.Sort(data.userScores)
}

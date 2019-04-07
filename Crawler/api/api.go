package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type apiType int

const (
	allContestsType apiType = iota
	scoreboardType
	specificContestType
	specificHandleType
)

const (
	allContestsAPI     string = "https://codejam.googleapis.com/poll?p=e30"
	scoreboardAPI      string = "https://codejam.googleapis.com/scoreboard/%s/poll?p=%s"
	specificContestAPI string = "https://codejam.googleapis.com/dashboard/%s/poll?p=e30"
	specificHandleAPI  string = "https://codejam.googleapis.com/scoreboard/%s/find?p=%s"
)

func fetchAPI(url string) []byte {
	// log.Println("url", url)
	resp, err := http.Get(url)
	handleErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	// log.Println(string(body))

	result, err := decodeFromBase64(body)
	handleErr(err)
	// log.Println(string(result))

	return result
}

func fetchAPIResponse(fetchType apiType, contestID string, param []interface{}) interface{} {
	url := ""

	switch fetchType {
	case specificHandleType: // handle search
		handle := param[0].(string)
		url = fmt.Sprintf(specificHandleAPI, contestID, getHandleSearchPayload(handle))
	case scoreboardType: // dump scoreboard
		starting := param[0].(int)
		step := param[1].(int)
		url = fmt.Sprintf(scoreboardAPI, contestID, getScoreboardPaginationPayload(starting, step))
	case specificContestType:
		contestID := param[0].(string)
		url = fmt.Sprintf(specificContestAPI, contestID)
	case allContestsType:
		url = fmt.Sprintf(allContestsAPI)
	default:
		log.Fatalln("Unknown option")
	}

	result := fetchAPI(url)

	switch fetchType {
	case specificHandleType: // handle search
		var response scoreboardResponse
		err := json.Unmarshal(result, &response)
		handleErr(err)
		// log.Println(response)

		return &response
	case scoreboardType: // dump scoreboard
		var response scoreboardResponse
		err := json.Unmarshal(result, &response)
		handleErr(err)
		// log.Println(response)

		return &response
	case specificContestType:
		var response contestResponse
		err := json.Unmarshal(result, &response)
		handleErr(err)
		// log.Println(response)

		return &response
	case allContestsType:
		var response contestsResponse
		err := json.Unmarshal(result, &response)
		handleErr(err)
		// log.Println(response)

		return &response
	default:
		log.Fatalln("Unknown option")
		return nil
	}
}

func (data *ContestData) printUserRecord(user *userScore) {
	fmt.Printf("+====== %-15s (%s) ======+\n", user.Handle, user.Country)
	fmt.Printf("Rank %v Score %v\n\n", user.Rank, user.Score)

	for _, curTask := range data.challenge.Tasks {
		fmt.Printf("%-30s ", curTask.Title)
		totalPoint := 0
		for _, p := range curTask.Tests {
			totalPoint += p.Point
		}

		for _, task := range user.TasksInfo {
			if curTask.ID == task.TaskID {
				fmt.Printf("%2d / %2d\n", task.Point, totalPoint)
				goto found
			}
		}

		// not attempted, so no record is present
		fmt.Printf("%2d / %2d\n", 0, totalPoint)
	found:
	}
}

func (data *ContestData) GetHandleResults(handles []string, forceFetch bool) {
	results := make(userScores, len(handles))
	idx := 0

	if forceFetch {
		ch := make(chan userScore)
		for _, handle := range handles {
			go data.fetchHandleResult(handle, ch)
		}

		for i := 0; i < len(handles); i++ {
			tmp := <-ch
			if tmp.isEmpty {
				continue
			}
			results[idx] = tmp
			idx++
		}
	} else {
		// TODO: optimize this part with map
		for _, handle := range handles {
			for _, cachedHandle := range data.contestants {
				if cachedHandle.Handle == handle {
					results[idx] = cachedHandle
					idx++

					goto hasMatch
				}
			}

			log.Println("Handle", handle, "not found")
		hasMatch:
		}

	}

	results = results[:idx]
	sort.Sort(results)
	for _, user := range results {
		data.printUserRecord(&user)
	}
}

func (data *ContestData) GetAllContestantData(country string) {
	for _, contestant := range data.contestants {
		if len(country) > 0 {
			if contestant.Country == country {
				data.printUserRecord(&contestant)
			}
		} else if len(country) == 0 {
			data.printUserRecord(&contestant)
		}
	}
}

// GetJSONResponse dumps the response from the specified url
// The url must be one of the api requests
func GetJSONResponse(url string) {
	response := fetchAPI(url)
	fmt.Println(string(response))
}

func GetContestListing() {
	response := fetchAPIResponse(allContestsType, "", nil).(*contestsResponse)

	for _, contestGroup := range response.Adventures {
		fmt.Printf("%-30s (%s)\n", contestGroup.Title, contestGroup.CompetitionStr)
		for _, contest := range contestGroup.Challenges {
			fmt.Printf("%-30s %-10s %s\n", contest.Title, contest.ContestID, contest.AdditionalInfo)
		}
		fmt.Println()
	}
}

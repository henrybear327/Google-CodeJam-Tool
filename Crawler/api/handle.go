package api

import (
	"encoding/json"
	"log"
)

const findHandleURL = "https://codejam.googleapis.com/scoreboard/%s/find?p=%s"

type handleSearchPayload struct {
	/*
		{"nickname":"henrybear327","scoreboard_page_size":10}
	*/
	Handle   string `json:"nickname"`
	PageSize int    `json:"scoreboard_page_size"`
}

func (data *ContestMetadata) getHandleSearchPayload(handle string) string {
	payload := handleSearchPayload{Handle: handle, PageSize: 1}
	res, err := json.Marshal(payload)
	handleErr(err)
	return encodeToBase64(res)
}

func (data *ContestMetadata) fetchHandleResult(handle string, ch chan userScore) {
	param := make([]interface{}, 1)
	param[0] = handle
	response := data.fetchResponseBody(1, param)

	if len(response.UserScores) != 1 {
		log.Fatalln("Incorrect user count", len(response.UserScores))
	}

	ch <- response.UserScores[0]
}

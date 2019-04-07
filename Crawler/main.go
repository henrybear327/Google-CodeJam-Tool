package main

import (
	"GCJ-side-project/Crawler/api"
	"flag"
	"log"
)

var operation = flag.Int("operation", 1, "1 = fetch handles, 2 = fetch all contestant, 3 = fetch and get decoded json response")
var url = flag.String("url", "", "For operation 3")
var country = flag.String("country", "", "The country name (for operation 2)")

func main() {
	flag.Parse()

	// init
	// load config file
	config := parseConfigFile()

	// init the contest
	contest := api.ContestData{}
	err := contest.FetchContest(config.ContestID, config.ConcurrentFetch, false)
	if err != nil {
		log.Fatalln(err)
	}

	switch *operation {
	case 1:
		contest.GetHandleResults(config.Handles, true)
	case 2:
		contest.GetAllContestantData(*country)
	case 3:
		api.GetJSONResponse(*url)
	default:
		log.Fatalln("No such operation")
	}
}

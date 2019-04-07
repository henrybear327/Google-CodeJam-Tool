package main

import (
	"Google-CodeJam-Tool/Crawler/api"
	"flag"
	"log"
)

var operation = flag.Int("operation", 1, "1 = fetch handles, 2 = fetch all contestant, 3 = fetch and get decoded json response, 4 = get contest listing")
var url = flag.String("url", "", "For operation 3")
var country = flag.String("country", "", "The country name (for operation 2)")

func main() {
	flag.Parse()

	switch *operation {
	case 1, 2:
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
		default:
			log.Fatalln("Golang is broken")
		}
	case 3:
		api.GetJSONResponse(*url)
	case 4:
		api.GetContestListing()
	default:
		log.Fatalln("No such operation")
	}
}

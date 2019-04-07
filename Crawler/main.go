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

	apiObject := api.ContestMetadata{ContestID: config.ContestID, StepSize: 100}
	apiObject.FetchContestInfo()

	switch *operation {
	case 1:
		apiObject.GetHandleResults(config.Handles)
	case 2:
		apiObject.GetAllContestantData(*country)
	case 3:
		apiObject.GetJSONResponse(*url)
	default:
		log.Fatalln("No such operation")
	}
}

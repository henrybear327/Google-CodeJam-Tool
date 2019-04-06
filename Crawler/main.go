package main

import (
	"flag"
	"log"
)

var operation = flag.Int("operation", 1, "1 = fetch handles, 2 = fetch all contestant")

func main() {
	flag.Parse()

	// special characters causing base64 decoding error
	// test := "Pożeracz_pączków_z_lukrem"
	// lol := encodeToBase64([]byte(test))
	// fmt.Println(lol)
	// fmt.Println(decodeFromBase64([]byte(lol)))

	// init
	config := parseConfigFile()
	api := contestMetadata{contestID: config.ContestID, stepSize: 100}
	api.fetchContestInfo()

	switch *operation {
	case 1:
		api.GetHandleResults(config.Handles)
	case 2:
		// api.GetAllContestantData("Taiwan")
		api.GetAllContestantData("")
	default:
		log.Fatalln("No such operation")
	}
}

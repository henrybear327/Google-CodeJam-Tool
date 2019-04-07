package api

import (
	"encoding/base64"
	"log"
)

func handleErr(err error) {
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func encodeToBase64(input []byte) string {
	str := base64.RawURLEncoding.EncodeToString(input)
	return str
}

func decodeFromBase64(input []byte) ([]byte, error) {
	res := make([]byte, len(input))
	n, err := base64.RawURLEncoding.Decode(res, input)
	if err != nil {
		return nil, err
	}

	return res[:n], nil
}

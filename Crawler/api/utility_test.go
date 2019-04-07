package api

import (
	"testing"
)

func TestBase64Operations(t *testing.T) {
	// special characters causing base64 decoding error
	test := "Pożeracz_pączków_z_lukrem"
	encoded := encodeToBase64([]byte(test))
	_, err := decodeFromBase64([]byte(encoded))
	if err != nil {
		t.Fatal("Should decode base64 string successfully", err)
	}
}

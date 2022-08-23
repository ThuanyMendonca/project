package utils

import (
	"bytes"
	"encoding/json"
)

func MockInvalidPayload() *bytes.Buffer {
	payload := "Invalid_Payload"
	reqBody, _ := json.Marshal(&payload)

	return bytes.NewBuffer(reqBody)
}

func MockValidPayload(payload interface{}) *bytes.Buffer {
	reqBody, _ := json.Marshal(&payload)

	return bytes.NewBuffer(reqBody)
}

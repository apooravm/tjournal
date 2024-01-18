package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LogReqPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Log      string `json:"log"`
}

type JournError struct {
	Err    error
	Simple string
}

func (ce *JournError) Error() string {
	return fmt.Sprintf("%v", ce.Simple)
}

func main() {
	if err := sendPostReq("https://multi-serve.onrender.com/api/journal/", LogReqPayload{Username: "sb@nsAd*/", Password: "xxxxx", Log: "Day has been great"}); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Logged! 🚀")
	}
}

func sendPostReq(url string, body interface{}) error {
	// []byte
	payload, err := json.Marshal(body)
	if err != nil {
		return &JournError{
			Err:    err,
			Simple: "Error marshalling object",
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return &JournError{
			Err:    err,
			Simple: "Error creating request",
		}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &JournError{
			Err:    err,
			Simple: "Error sending request. ",
		}
	}

	fmt.Println(res.Status)

	res.Body.Close()
	return nil
}

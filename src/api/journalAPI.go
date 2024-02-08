package api

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

type JournalDB struct {
	Url      string
	Username string
	Password string
}

func (journal *JournalDB) ReadJournalLogs() (*[]ReadJournalLogRes, error) {
	payload, err := json.Marshal(UserAuth{
		Username: journal.Username,
		Password: journal.Password,
	})
	if err != nil {
		fmt.Println("Error marshalling")
		return nil, err
	}

	req, err := http.NewRequest("GET", journal.Url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return nil, err
	}

	defer res.Body.Close()

	var journalLogs []ReadJournalLogRes

	if err := json.NewDecoder(res.Body).Decode(&journalLogs); err != nil {
		fmt.Println("Error Unmarshallng data")
		return nil, err
	}

	return &journalLogs, nil
}

func (journal *JournalDB) CreateJournalLog(log string, title string, tags *[]string) (*JournalMessage, error) {
	payload, err := json.Marshal(CreateJournalLogReq{
		Username: journal.Username,
		Password: journal.Password,
		Log:      log,
		Tags:     *tags,
		Title:    title,
	})
	if err != nil {
		fmt.Println("Error creating data payload")
		return nil, err
	}

	req, err := http.NewRequest("POST", journal.Url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Println("Log created!", res.Status)
		return &JournalMessage{Message: res.Status, Code: res.StatusCode, Simple: "good"}, nil

	} else {
		fmt.Println("Something went wrong.", res.Status)
		return &JournalMessage{Message: res.Status, Code: res.StatusCode, Simple: "bad"}, nil
	}
}

// Create a copy of the original log obj and edit that itself. This becomes the new log
func (journal *JournalDB) UpdateJournalLog(prevLog *ReadJournalLogRes) (*JournalMessage, error) {
	payload, err := json.Marshal(UpdateLogReq{
		Username: journal.Username,
		Password: journal.Password,
		Log:      prevLog.Log,
		Tags:     prevLog.Tags,
		Title:    prevLog.Title,
		Log_Id:   prevLog.Log_Id,
	})
	if err != nil {
		fmt.Println("Error marshalling payload")
		return nil, err
	}

	req, err := http.NewRequest("PUT", journal.Url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Println("Updated Successfully")
		return &JournalMessage{Message: res.Status, Code: res.StatusCode, Simple: "good"}, nil
	} else {
		fmt.Println("Something went wrong")
		return &JournalMessage{Message: res.Status, Code: res.StatusCode, Simple: "bad"}, nil
	}
}

func (journal *JournalDB) DeleteJournalLog(log_id int) (*JournalMessage, error) {
	payload, err := json.Marshal(DeleteJournalLogReq{
		Username: journal.Username,
		Password: journal.Password,
		Log_Id:   log_id,
	})
	if err != nil {
		fmt.Println("Error marshalling payload")
		return nil, err
	}

	req, err := http.NewRequest("DELETE", journal.Url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Println("Deleted Successfully")
		return &JournalMessage{Message: res.Status, Code: res.StatusCode, Simple: "good"}, nil
	} else {
		fmt.Println("Something went wrong", res.Status)
		return &JournalMessage{Message: res.Status, Code: res.StatusCode, Simple: "bad"}, nil
	}
}

// func WriteJournalLog() {
// 	var log string
// 	var title string
// 	var tags_str string

// 	fmt.Println("log: ")
// 	scanner := bufio.NewScanner(os.Stdin)
// 	if scanner.Scan() {
// 		log = scanner.Text()
// 	}

// 	fmt.Println("title: ")
// 	scanner = bufio.NewScanner(os.Stdin)
// 	if scanner.Scan() {
// 		title = scanner.Text()
// 	}

// 	fmt.Println("tags (separated with spaces): ")
// 	scanner = bufio.NewScanner(os.Stdin)
// 	if scanner.Scan() {
// 		tags_str = scanner.Text()
// 	}

// 	tags := strings.Split(tags_str, " ")

//		// fmt.Println("LOG:", log)
//		// fmt.Println("TITLE:", title)
//		// fmt.Println("TAGS:", tags)
//		CreateJournalLog(log, title, &tags)
//	}

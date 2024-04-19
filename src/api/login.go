package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type UserAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRes struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

// Error response object from the server
type ServerErrorRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Simple  string `json:"simple"`
}

func (e ServerErrorRes) Error() string {
	return fmt.Sprintf("Error %d: %s; %s", e.Code, e.Simple, e.Message)
}

// Logs in the user and returns the received token
func LoginUser(urlEndpoint string, email string, password string) (*AuthRes, error) {
	payload, err := json.Marshal(UserAuth{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, ServerErrorRes{
			Code:    400,
			Message: err.Error(),
			Simple:  "Error marshaling the data.",
		}
	}
	req, err := http.NewRequest("POST", urlEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, ServerErrorRes{
			Code:    400,
			Message: err.Error(),
			Simple:  "Error creating request.",
		}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, ServerErrorRes{
			Code:    400,
			Message: err.Error(),
			Simple:  "Error sending request.",
		}
	}

	defer res.Body.Close()

	switch {
	case res.StatusCode >= 200 && res.StatusCode < 300:
		// Success
		var tokenObj AuthRes
		if err := json.NewDecoder(res.Body).Decode(&tokenObj); err != nil {
			return nil, ServerErrorRes{
				Code:    400,
				Message: err.Error(),
				Simple:  "Error unmarshaling response.",
			}
		}

		tokenStr := strings.Split(tokenObj.Token, " ")[1]
		tokenObj.Token = tokenStr
		return &tokenObj, nil

	default:
		var errorRes ServerErrorRes
		if err := json.NewDecoder(res.Body).Decode(&errorRes); err != nil {
			return nil, ServerErrorRes{
				Code:    400,
				Message: err.Error(),
				Simple:  "Error unmarshaling response.",
			}
		}

		return nil, errorRes
	}
}

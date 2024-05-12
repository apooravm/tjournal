package api

import (
	"fmt"
	"net"
	"time"
	"net/http"
)

// Check whether the user has an active internet connection
func userIsConnected() bool {
	timeout := time.Duration(2 * time.Second)

	conn, err := net.DialTimeout("tcp", "www.google.com:80", timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}

// Ping and check whether the server is online
func  CheckServerStatus(pingUrl string) (bool, error) {
	req, err := http.NewRequest("GET", pingUrl, nil)
	if err != nil {
		return false, fmt.Errorf("error creating a request: %s", err.Error())
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		// Either server is offline or the user has no internet
		return false, nil
	}

	defer res.Body.Close()
	return true, nil
}


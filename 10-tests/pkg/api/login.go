package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// start with capital letter so it is exported
func doLoginRequest(client ClientIface, requestURL, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("marshall error: %s", err)
	}

	res, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("HTTP POST error: %s", err)
	}

	defer res.Body.Close()

	// now we're ready to read the body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}

	// parse data if response code is 200
	if res.StatusCode != 200 {
		return "", fmt.Errorf("invalid output (HTTP Code: %d): %s", res.StatusCode, string(resBody))
	}

	if !json.Valid(resBody) {
		return "", RequestError{
			HTTPCode: res.StatusCode,
			Body:     string(resBody),
			Err:      "Invalid JSON body error",
		}
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", RequestError{
			HTTPCode: res.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("Page unmarsall error: %s", err),
		}
	}

	if loginResponse.Token == "" {
		return "", RequestError{
			HTTPCode: res.StatusCode,
			Body:     string(resBody),
			Err:      "Empty token replied",
		}
	}

	return loginResponse.Token, nil
}

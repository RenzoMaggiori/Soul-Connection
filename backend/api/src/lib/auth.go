package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type LoginCredentials struct {
	XGroupAuthentication string
	AuthEmail            string
	AuthPassword         string
}

func Auth(credentials LoginCredentials) (string, error) {
	resp, err := Fetch(&http.Client{}, FetchRequest{
		Method:  "POST",
		Url:     fmt.Sprintf("%s/api/employees/login", ApiBaseUri),
		Body:    map[string]string{"email": credentials.AuthEmail, "password": credentials.AuthPassword},
		Headers: map[string]string{"Content-Type": "application/json", "X-Group-Authorization": credentials.XGroupAuthentication},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var loginResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		detail, ok := loginResponse["detail"]
		if !ok {
			return "", errors.New("could not authenticate to api")
		}
		return "", errors.New(detail)
	}
	jwt, ok := loginResponse["access_token"]
	if !ok {
		return "", errors.New("could not parse authentication response from api")
	}
	return jwt, nil
}

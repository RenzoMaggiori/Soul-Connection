package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type FetchRequest struct {
	Method  string
	Url     string
	Body    map[string]string
	Headers map[string]string
}

func Fetch(client *http.Client, request FetchRequest) (*http.Response, error) {
	var body io.Reader
	if request.Body != nil {
		jsonBody, err := json.Marshal(request.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonBody)
	} else {
		body = nil
	}
	req, err := http.NewRequest(request.Method, request.Url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

package router

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func post(url string, bodyString string, header map[string]string) (int, string, error) {
	body := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return 0, "", err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	rs, err := client.Do(req)
	if err != nil {
		return 0, "", errors.New("Failed to create named key request: " + err.Error())
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return 0, "", errors.New("Failed to parse response body: " + err.Error())
	}

	return rs.StatusCode, string(bodyBytes), nil
}

func parseRequestBody(body io.ReadCloser) map[string]string {
	decoder := json.NewDecoder(body)
	var decodedBody map[string]string
	decoder.Decode(&decodedBody)

	return decodedBody
}

func copyHeaders(from http.Header) map[string]string {
	to := map[string]string{}
	for key, values := range from {
		to[key] = values[0]
	}

	return to
}

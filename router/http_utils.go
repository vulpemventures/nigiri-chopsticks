package router

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

func get(url string, headers map[string]string) (int, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return resp.StatusCode, string(respBody), nil
}

func post(url string, bodyString string, headers map[string]string) (int, string, error) {
	body := strings.NewReader(bodyString)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return resp.StatusCode, string(respBody), nil
}

func parseRequestBody(body io.ReadCloser) map[string]string {
	decoder := json.NewDecoder(body)
	var decodedBody map[string]string
	decoder.Decode(&decodedBody)

	return decodedBody
}

package rpc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

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

func parseResponseError(resp string) string {
	out := map[string]map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	return out["error"]["message"].(string)
}

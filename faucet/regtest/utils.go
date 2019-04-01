package regtestfaucet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func listunspent(url string) (int, string, error) {
	body := `{
		"jsonrpc": "1.0",
		"id": "2",
		"method": "listunspent",
		"params": []
	}`

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		return status, "", fmt.Errorf("an error occured while listing unspents with status %d: %s", status, resp)
	}

	out := map[string]interface{}{}
	err = json.Unmarshal([]byte(resp), &out)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("an error occured while unmarshaling unspents: %s", err)
	}
	utxo, _ := json.Marshal(out["result"].([]interface{})[0])

	return status, string(utxo), nil
}

func getrawchangeaddress(url string) (int, string, error) {
	body := `{
		"jsonrpc": "1.0",
		"id": "faucet",
		"method": "getrawchangeaddress",
		"params": []
	}`

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		return status, "", fmt.Errorf("an error occured while getting change address with status %d: %s", status, resp)
	}

	out := map[string]string{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"], nil
}

func createrawtransaction(url string, utxo string, address string, changeAddress string) (int, string, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "faucet",
		"method": "createrawtransaction",
		"params": [
			[%s],
			[
				{"%s": 1},
				{"%s": 48.99}
			]
		]
	}`, utxo, address, changeAddress)

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		return status, "", fmt.Errorf("an error occured while creating raw transction with status %d: %s", status, resp)
	}

	out := map[string]string{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"], nil
}

func signrawtransaction(url string, tx string) (int, string, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "faucet",
		"method": "signrawtransactionwithwallet",
		"params": ["%s"]
	}`, tx)

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		return status, "", fmt.Errorf("an error occured while signing transaction with status %d: %s", status, resp)
	}

	out := map[string]map[string]interface{}{}
	json.Unmarshal([]byte(resp), &out)

	return status, out["result"]["hex"].(string), nil
}

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

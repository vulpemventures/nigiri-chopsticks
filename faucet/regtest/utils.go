package regtestfaucet

import (
	"encoding/json"
	"errors"
	"fmt"
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

func listunspent(url string) (int, string, error) {
	body := `{"jsonrpc": "1.0", "id": "2", "method": "listunspent", "params": []}`

	status, resp, err := post(url, body, nil)
	if err != nil {
		return status, "", err
	}
	if status != http.StatusOK {
		return status, "", fmt.Errorf("an error occured while listing unspents with status %d: %s", status, resp)
	}

	type vout struct {
		TxID string `json:"txid"`
		VOut uint   `json:"vout"`
	}

	type response struct {
		Result []vout `json:"result"`
	}

	out := &response{}
	err = json.Unmarshal([]byte(resp), out)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("error while unmarshaling unspents: %s", err)
	}
	utxo, _ := json.Marshal(out.Result[0])

	return status, string(utxo), nil
}

func createrawtransaction(url string, utxo string, address string) (int, string, error) {
	status, changeAddress, err := getchangeaddress(url)

	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "2",
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
	body := fmt.Sprintf(`{"jsonrpc": "1.0", "id": "2", "method": "signrawtransactionwithwallet", "params": ["%s"]}`, tx)

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

func getchangeaddress(url string) (int, string, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "1.0",
		"id": "2",
		"method": "getrawchangeaddress",
		"params": []
		}`)
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

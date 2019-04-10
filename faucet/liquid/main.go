package liquidfaucet

import (
	"encoding/json"
	"net/http"

	"github.com/vulpemventures/nigiri-chopsticks/faucet"
	"github.com/vulpemventures/nigiri-chopsticks/helpers"
)

type liquidfaucet struct {
	URL       string
	rpcClient *helpers.RpcClient
}

// NewFaucet initialize a liquid faucet and returns it as interface
func NewFaucet(url string, client *helpers.RpcClient) faucet.Faucet {
	return &liquidfaucet{url, client}
}

func (f *liquidfaucet) NewTransaction(address string) (int, string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "sendtoaddress", []interface{}{address, 1})
	if err != nil {
		return status, "", err
	}

	return status, resp.(string), nil
}

func (f *liquidfaucet) Fund() (int, []string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "getbalance", nil)
	if err != nil {
		return status, nil, err
	}

	balance := resp.(map[string]interface{})["bitcoin"].(float64)
	if balance <= 0 {
		status, resp, err := handleRPCRequest(f.rpcClient, "generate", []interface{}{101})
		if err != nil {
			return status, nil, err
		}

		blockHashes := []string{}
		for _, b := range resp.([]interface{}) {
			blockHashes = append(blockHashes, b.(string))
		}
		return status, blockHashes, nil
	}

	return 200, nil, nil
}

func handleRPCRequest(client *helpers.RpcClient, method string, params []interface{}) (int, interface{}, error) {
	status, resp, err := client.Call(method, params)
	if err != nil {
		return status, "", err
	}
	var out interface{}
	err = json.Unmarshal(resp.Result, &out)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return status, out, nil
}

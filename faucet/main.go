package faucet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vulpemventures/nigiri-chopsticks/helpers"
)

type Faucet struct {
	URL       string
	rpcClient *helpers.RpcClient
}

// NewFaucet initialize a liquid faucet and returns it as interface
func NewFaucet(url string, client *helpers.RpcClient) *Faucet {
	return &Faucet{url, client}
}

func (f *Faucet) NewTransaction(address string) (int, string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "sendtoaddress", []interface{}{address, 1})
	if err != nil {
		return status, "", err
	}

	return status, resp.(string), nil
}

// liquid starts with initialfreecoins = 21,000,000 LBTC so we just need to
// "mature" the balance mining 101 blocks if not already mined
func (f *Faucet) Fund() (int, []string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "getblockcount", nil)
	if err != nil {
		return status, nil, err
	}

	if blockCount := resp.(float64); blockCount <= 0 {
		return f.Mine(101)
	}

	return 200, nil, nil
}

func (f *Faucet) Mine(blocks int) (int, []string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "getnewaddress", nil)
	if err != nil {
		return status, nil, err
	}
	address := resp.(string)

	status, resp, err = handleRPCRequest(f.rpcClient, "generatetoaddress", []interface{}{blocks, address})
	if err != nil {
		return status, nil, err
	}

	blockHashes := []string{}
	for _, b := range resp.([]interface{}) {
		blockHashes = append(blockHashes, b.(string))
	}

	return status, blockHashes, nil
}

func (f *Faucet) Mint(address string, quantity float64) (int, string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "issueasset", []interface{}{quantity, 0, false})
	if err != nil {
		return status, "", err
	}
	asset := resp.(map[string]interface{})["asset"].(string)

	status, tx, err := handleRPCRequest(f.rpcClient, "sendtoaddress", []interface{}{address, quantity, "", "", false, false, 1, "UNSET", asset})
	if err != nil {
		return status, "", err
	}

	resp = fmt.Sprintf(`{"asset": %s, "txId": %s}`, asset, tx.(string))
	return status, resp.(string), nil
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

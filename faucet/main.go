package faucet

import (
	"encoding/json"
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

// SendBitcoinTransaction calls the sendtoaddress method of the bitcoin node to the given address with the fractional amount
func (f *Faucet) SendBitcoinTransaction(address string, amount float64) (int, string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "sendtoaddress", []interface{}{address, amount})
	if err != nil {
		return status, "", err
	}

	return status, resp.(string), nil
}

// SendLiquidTransaction calls the sendtoaddress method of the elements node to the given address with the fractional amount of the given asset hash.
// If asset hash is empty will send Liquid Bitcoin
func (f *Faucet) SendLiquidTransaction(address string, amount float64, asset string) (int, string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "sendtoaddress", []interface{}{address, amount, "", "", false, false, 1, "UNSET", asset})
	if err != nil {
		return status, "", err
	}

	return status, resp.(string), nil
}

// Fund  "mature" the balance mining 1 block if not already mined
//liquid starts with initialfreecoins = 21,000,000 LBTC
func (f *Faucet) Fund() (int, []string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "getblockcount", nil)
	if err != nil {
		return status, nil, err
	}

	if blockCount := resp.(float64); blockCount <= 0 {
		return f.Mine(1)
	}

	return 200, nil, nil
}

// Mine will generated block versus an address of the wallet
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

// Mint issue a new Liquid asset
func (f *Faucet) Mint(address string, quantity float64) (int, map[string]interface{}, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "issueasset", []interface{}{quantity, 0, false})
	if err != nil {
		return status, nil, err
	}
	decodedResp := resp.(map[string]interface{})
	asset := decodedResp["asset"].(string)
	issuanceInput := map[string]interface{}{
		"txid": decodedResp["txid"].(string),
		"vin":  decodedResp["vin"].(float64),
	}

	status, tx, err := f.SendLiquidTransaction(address, quantity, asset)
	if err != nil {
		return status, nil, err
	}

	res := make(map[string]interface{})
	res["asset"] = asset
	res["txId"] = tx
	res["issuance_txin"] = issuanceInput

	return status, res, nil
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

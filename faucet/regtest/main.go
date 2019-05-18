package regtestfaucet

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/vulpemventures/nigiri-chopsticks/faucet"
	"github.com/vulpemventures/nigiri-chopsticks/helpers"
)

type regtestfaucet struct {
	URL       string
	rpcClient *helpers.RpcClient
}

// NewFaucet initialize a regtest faucet and returns it as interface
func NewFaucet(url string, client *helpers.RpcClient) faucet.Faucet {
	return &regtestfaucet{url, client}
}

func (f *regtestfaucet) NewTransaction(address string) (int, string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "listunspent", nil)
	if err != nil {
		return status, "", err
	}
	utxo := resp.([]interface{})[0]
	utxoAmount := utxo.(map[string]interface{})["amount"].(float64)

	status, resp, err = handleRPCRequest(f.rpcClient, "getrawchangeaddress", nil)
	if err != nil {
		return status, "", err
	}
	changeAddress := resp.(string)
	changeAmount := math.Round((utxoAmount-1.01)*100) / 100

	vin := []interface{}{utxo}
	vout := []map[string]float64{
		map[string]float64{address: 1.0},
		map[string]float64{changeAddress: changeAmount},
	}
	status, resp, err = handleRPCRequest(f.rpcClient, "createrawtransaction", []interface{}{vin, vout})
	if err != nil {
		return status, "", err
	}
	tx := resp.(string)

	status, resp, err = handleRPCRequest(f.rpcClient, "signrawtransactionwithwallet", []interface{}{tx})
	if err != nil {
		return status, "", err
	}

	signedTx := resp.(map[string]interface{})["hex"].(string)

	return status, signedTx, nil
}

func (f *regtestfaucet) Fund() (int, []string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "getbalance", nil)
	if err != nil {
		return status, nil, err
	}

	if balance := resp.(float64); balance <= 0 {
		return f.Mine(101)
	}

	return 200, nil, nil
}

func (f *regtestfaucet) Mine(blocks int) (int, []string, error) {
	status, resp, err := handleRPCRequest(f.rpcClient, "getnewaddress", nil)
	if err != nil {
		return status, nil, err
	}
	address := resp.(string)

	status, resp, err = handleRPCRequest(f.rpcClient, "generatetoaddress", []interface{}{101, address})
	if err != nil {
		return status, nil, err
	}

	blockHashes := []string{}
	for _, b := range resp.([]interface{}) {
		blockHashes = append(blockHashes, b.(string))
	}
	return status, blockHashes, nil
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

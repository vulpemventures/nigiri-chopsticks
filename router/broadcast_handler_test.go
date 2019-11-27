package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/altafan/nigiri-chopsticks/helpers"
)

func TestBroadcastLiquidTransaction(t *testing.T) {
	r := NewTestRouter(withLiquid)
	client := r.RPCClient

	signedTx, err := getSignedTransaction(client)
	if err != nil {
		t.Fatal(err)
	}

	blockCountResp := blockCountRequest(r)
	if blockCountResp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, blockCountResp.Code)
	}
	prevBlockCount, _ := strconv.Atoi(blockCountResp.Body.String())

	resp := broadcastRequest(r, signedTx)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status: %d, got: %d with error: %s\n", http.StatusOK, resp.Code, resp.Body.String())
	}

	time.Sleep(5 * time.Second)

	blockCountResp = blockCountRequest(r)
	if blockCountResp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, blockCountResp.Code)
	}
	blockCount, _ := strconv.Atoi(blockCountResp.Body.String())
	blockMined := blockCount - prevBlockCount
	if blockMined != 10 {
		t.Fatalf("Expected 1 block mined, got %d\n", blockMined)
	}
}

func broadcastRequest(r *Router, txHex string) *httptest.ResponseRecorder {
	payload := []byte(txHex)
	req, _ := http.NewRequest("POST", "/tx", bytes.NewBuffer(payload))

	return doRequest(r, req)
}

func getSignedTransaction(r *helpers.RpcClient) (string, error) {
	confidentialAddress, err := getNewAddress(r)
	if err != nil {
		return "", err
	}
	unconfidentialAddress, err := getUnconfidentialAddress(r, confidentialAddress)
	if err != nil {
		return "", err
	}

	rawTx, err := createRawTransaction(r, unconfidentialAddress, 1.0)
	if err != nil {
		return "", err
	}

	fundedTx, err := fundRawTransaction(r, rawTx)
	if err != nil {
		return "", err
	}

	blindedTx, err := blindRawTransaction(r, fundedTx)
	if err != nil {
		return "", err
	}

	signedTx, err := signRawTransaction(r, blindedTx)
	if err != nil {
		return "", err
	}

	return signedTx, nil
}

func getNewAddress(r *helpers.RpcClient) (string, error) {
	_, resp, err := handleRPCRequest(r, "getnewaddress", nil)
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}

func getUnconfidentialAddress(r *helpers.RpcClient, confidential string) (string, error) {
	_, resp, err := handleRPCRequest(r, "validateaddress", []interface{}{confidential})
	if err != nil {
		return "", err
	}
	return resp.(map[string]interface{})["unconfidential"].(string), nil
}

func createRawTransaction(r *helpers.RpcClient, receiver string, amount float64) (string, error) {
	inputs := []interface{}{}
	outputs := map[string]float64{receiver: amount}
	_, resp, err := handleRPCRequest(r, "createrawtransaction", []interface{}{inputs, outputs})
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}

func fundRawTransaction(r *helpers.RpcClient, txHex string) (string, error) {
	_, resp, err := handleRPCRequest(r, "fundrawtransaction", []interface{}{txHex})
	if err != nil {
		return "", err
	}
	return resp.(map[string]interface{})["hex"].(string), nil
}

func blindRawTransaction(r *helpers.RpcClient, txHex string) (string, error) {
	_, resp, err := handleRPCRequest(r, "blindrawtransaction", []interface{}{txHex})
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}

func signRawTransaction(r *helpers.RpcClient, txHex string) (string, error) {
	_, resp, err := handleRPCRequest(r, "signrawtransactionwithwallet", []interface{}{txHex})
	if err != nil {
		return "", err
	}
	return resp.(map[string]interface{})["hex"].(string), nil
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

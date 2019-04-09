package liquidfaucet

import (
	"github.com/vulpemventures/nigiri-chopsticks/faucet"
	"github.com/vulpemventures/nigiri-chopsticks/helpers/liquidrpc"
)

type liquidfaucet struct {
	URL string
}

// NewFaucet initialize a liquid faucet and returns it as interface
func NewFaucet(url string) faucet.Faucet {
	return &liquidfaucet{url}
}

func (f *liquidfaucet) NewTransaction(address string) (int, string, error) {
	status, txHash, err := liquidrpc.Sendtoaddress(f.URL, address, 1)
	if err != nil {
		return status, "", err
	}

	return status, txHash, nil
}

func (f *liquidfaucet) Fund() (int, []string, error) {
	status, balance, err := liquidrpc.Getbalance(f.URL)
	if err != nil {
		return status, nil, err
	}

	if balance <= 0 {
		return liquidrpc.Mine(f.URL, 200)
	}

	return 200, nil, nil
}
